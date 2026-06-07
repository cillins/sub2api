package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────────────────────────────
// MuleRunGatewayService — proxy for MuleRun platform accounts
// ─────────────────────────────────────────────────────────────────

const (
	muleRunDefaultBaseURL  = "https://api.mulerun.com"
	muleRunVendorTimeout   = 120 * time.Second
	muleRunReadTimeout     = 30 * time.Second
	muleRunTextTimeout     = 300 * time.Second
	muleRunMaxAccountTries = 10
)

// MuleRunGatewayService handles forwarding to MuleRun upstream.
type MuleRunGatewayService struct {
	httpUpstream HTTPUpstream
}

// NewMuleRunGatewayService creates a new MuleRun gateway service.
func NewMuleRunGatewayService(httpUpstream HTTPUpstream) *MuleRunGatewayService {
	return &MuleRunGatewayService{
		httpUpstream: httpUpstream,
	}
}

// ──────────────────────────── helpers ────────────────────────────

// getMuleRunBaseURL returns the upstream base URL for a MuleRun account.
func getMuleRunBaseURL(account *Account) string {
	if base := account.GetCredential("base_url"); base != "" {
		return strings.TrimRight(base, "/")
	}
	return muleRunDefaultBaseURL
}

// getMuleRunToken extracts the bearer token from a MuleRun account's credentials.
func getMuleRunToken(account *Account) (string, error) {
	// Prefer api_key (muk-xxx format for AI API calls), fall back to access_token (OAuth JWT)
	if tok := account.GetCredential("api_key"); tok != "" {
		return tok, nil
	}
	if tok := account.GetCredential("access_token"); tok != "" {
		return tok, nil
	}
	return "", fmt.Errorf("mulerun: no api_key or access_token in credentials for account %d", account.ID)
}

// muleRunShouldFailover returns true for HTTP status codes that should trigger
// account failover (try next key).
func muleRunShouldFailover(statusCode int) bool {
	switch statusCode {
	case 402, 403, 429, 502, 503, 504:
		return true
	}
	return false
}

// ────────────────────────────────────────────────────────────────
// Forward — text endpoints (/v1/messages, /v1/chat/completions)
// ────────────────────────────────────────────────────────────────

// MuleRunForwardRequest contains the parameters for a MuleRun text forward.
type MuleRunForwardRequest struct {
	Body      []byte            // raw request body from the client
	Model     string            // model name from the parsed request
	Stream    bool              // whether the client requested streaming
	Endpoint  string            // upstream endpoint path (e.g. "/v1/messages")
}

// Forward sends a text request to MuleRun upstream and returns the result.
// It handles both streaming (SSE) and non-streaming responses.
func (s *MuleRunGatewayService) Forward(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	req *MuleRunForwardRequest,
) (*ForwardResult, error) {
	startTime := time.Now()

	token, err := getMuleRunToken(account)
	if err != nil {
		return nil, &UpstreamFailoverError{StatusCode: 502, RetryableOnSameAccount: false}
	}

	baseURL := getMuleRunBaseURL(account)
	targetURL := baseURL + req.Endpoint

	// Build the upstream request
	upReq, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("mulerun: build request: %w", err)
	}

	upReq.Header.Set("Content-Type", "application/json")
	upReq.Header.Set("Authorization", "Bearer "+token)
	upReq.Header.Set("Accept", "text/event-stream")

	// Determine proxy
	proxyURL := ""
	if account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	// Execute with timeout
	upCtx, upCancel := context.WithTimeout(context.WithoutCancel(ctx), muleRunTextTimeout)
	defer upCancel()
	upReq = upReq.WithContext(upCtx)

	resp, err := s.httpUpstream.Do(upReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		slog.Warn("mulerun upstream request failed",
			"account_id", account.ID,
			"error", err,
		)
		return nil, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: false,
		}
	}
	defer resp.Body.Close()

	// Check for failover-worthy errors
	if resp.StatusCode >= 400 && muleRunShouldFailover(resp.StatusCode) {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		slog.Info("mulerun upstream returned failover error",
			"account_id", account.ID,
			"status", resp.StatusCode,
			"body_preview", truncateBody(body, 200),
		)
		return nil, &UpstreamFailoverError{
			StatusCode:      resp.StatusCode,
			ResponseBody:    body,
			ResponseHeaders: resp.Header,
		}
	}

	// Non-2xx non-failover errors: return as terminal error
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		c.Status(resp.StatusCode)
		for k, vs := range resp.Header {
			for _, v := range vs {
				c.Header(k, v)
			}
		}
		c.Writer.Write(body)
		return nil, fmt.Errorf("mulerun upstream error: %d", resp.StatusCode)
	}

	// Success path: streaming or non-streaming
	isStream := req.Stream || isEventStreamResponse(resp.Header)

	if isStream {
		return s.handleStreamResponse(ctx, c, resp, account, startTime, req.Model)
	}
	return s.handleNonStreamResponse(ctx, c, resp, account, startTime, req.Model)
}

// ─────────────────────── streaming response ──────────────────────

func (s *MuleRunGatewayService) handleStreamResponse(
	ctx context.Context,
	c *gin.Context,
	resp *http.Response,
	account *Account,
	startTime time.Time,
	model string,
) (*ForwardResult, error) {
	// Set SSE response headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 512*1024), 10*1024*1024) // 10MB max line

	var usage ClaudeUsage
	var firstTokenMs *int
	var clientDisconnect bool
	var gotTerminal bool

	for scanner.Scan() {
		line := scanner.Bytes()
		lineStr := string(line)

		// Write line to client
		_, writeErr := fmt.Fprintf(c.Writer, "%s\n", lineStr)
		if writeErr != nil {
			clientDisconnect = true
			// Continue draining to collect usage
		}
		c.Writer.Flush()

		// Track first token time
		if firstTokenMs == nil && !clientDisconnect {
			elapsed := int(time.Since(startTime).Milliseconds())
			firstTokenMs = &elapsed
		}

		// Parse usage from SSE data events
		if strings.HasPrefix(lineStr, "data: ") {
			data := strings.TrimPrefix(lineStr, "data: ")
			if data == "[DONE]" {
				gotTerminal = true
				break
			}
			// Try to extract usage from the data
			var evt map[string]json.RawMessage
			if json.Unmarshal([]byte(data), &evt) == nil {
				extractMuleRunUsage(evt, &usage)
				// Check for terminal event
				if evtType, ok := evt["type"]; ok {
					var tp string
					if json.Unmarshal(evtType, &tp) == nil {
						if tp == "message_stop" || tp == "message_delta" {
							// Try to get usage from message_delta
							if usg, ok := evt["usage"]; ok {
								json.Unmarshal(usg, &usage)
							}
						}
					}
				}
			}
		}

		// message_stop = terminal
		if strings.Contains(lineStr, "message_stop") {
			gotTerminal = true
		}
	}

	if err := scanner.Err(); err != nil && !clientDisconnect {
		return nil, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: true,
		}
	}

	if !gotTerminal && !clientDisconnect {
		// No terminal event — treat as potentially incomplete
		slog.Warn("mulerun stream ended without terminal event",
			"account_id", account.ID,
			"model", model,
		)
	}

	duration := time.Since(startTime)
	return &ForwardResult{
		Model:            model,
		UpstreamModel:    model,
		Stream:           true,
		Duration:         duration,
		FirstTokenMs:     firstTokenMs,
		ClientDisconnect: clientDisconnect,
		Usage:            usage,
	}, nil
}

// ──────────────────── non-streaming response ─────────────────────

func (s *MuleRunGatewayService) handleNonStreamResponse(
	_ context.Context,
	c *gin.Context,
	resp *http.Response,
	account *Account,
	startTime time.Time,
	model string,
) (*ForwardResult, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: true,
		}
	}

	// Parse usage from response
	var usage ClaudeUsage
	var respBody map[string]json.RawMessage
	if json.Unmarshal(body, &respBody) == nil {
		extractMuleRunUsage(respBody, &usage)
	}

	// Forward response headers
	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Header(k, v)
		}
	}
	c.Status(resp.StatusCode)
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)

	duration := time.Since(startTime)
	return &ForwardResult{
		Model:         model,
		UpstreamModel: model,
		Stream:        false,
		Duration:      duration,
		Usage:         usage,
	}, nil
}

// ────────────────────────────────────────────────────────────────
// ForwardVendor — vendor endpoints (/vendors/*)
// ────────────────────────────────────────────────────────────────

// ForwardVendor forwards a vendor API request to MuleRun upstream.
// It handles both POST (task creation) and GET (task polling) requests.
func (s *MuleRunGatewayService) ForwardVendor(
	ctx context.Context,
	c *gin.Context,
	account *Account,
) (*ForwardResult, error) {
	startTime := time.Now()

	token, err := getMuleRunToken(account)
	if err != nil {
		return nil, &UpstreamFailoverError{StatusCode: 502}
	}

	baseURL := getMuleRunBaseURL(account)
	// Preserve the full vendor path: /vendors/{provider}/v1/{model}/{action}[/{taskId}]
	vendorPath := c.Request.URL.Path
	targetURL := baseURL + vendorPath

	// Copy query string if present
	if qs := c.Request.URL.RawQuery; qs != "" {
		targetURL += "?" + qs
	}

	method := c.Request.Method
	timeout := muleRunReadTimeout
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		timeout = muleRunVendorTimeout
	}

	// Read request body for methods that have one
	var bodyReader io.Reader
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, fmt.Errorf("mulerun vendor: read body: %w", err)
		}
		bodyReader = bytes.NewReader(body)
	}

	upCtx, upCancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
	defer upCancel()

	upReq, err := http.NewRequestWithContext(upCtx, method, targetURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("mulerun vendor: build request: %w", err)
	}

	// Copy relevant headers
	upReq.Header.Set("Authorization", "Bearer "+token)
	upReq.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	if ct := c.GetHeader("Content-Type"); ct == "" && bodyReader != nil {
		upReq.Header.Set("Content-Type", "application/json")
	}
	upReq.Header.Set("Accept", "application/json")
	// Pass through MuleRun-specific headers
	for _, h := range []string{"X-Agent-Skills", "User-Agent"} {
		if v := c.GetHeader(h); v != "" {
			upReq.Header.Set(h, v)
		}
	}

	proxyURL := ""
	if account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(upReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		slog.Warn("mulerun vendor request failed",
			"account_id", account.ID,
			"method", method,
			"path", vendorPath,
			"error", err,
		)
		return nil, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: false,
		}
	}
	defer resp.Body.Close()

	// Check for failover-worthy errors
	if resp.StatusCode >= 400 && muleRunShouldFailover(resp.StatusCode) {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, &UpstreamFailoverError{
			StatusCode:      resp.StatusCode,
			ResponseBody:    body,
			ResponseHeaders: resp.Header,
		}
	}

	// Non-2xx non-failover errors: forward to client
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 65536))
		c.Status(resp.StatusCode)
		for k, vs := range resp.Header {
			for _, v := range vs {
				c.Header(k, v)
			}
		}
		c.Writer.Write(body)
		return nil, fmt.Errorf("mulerun vendor upstream error: %d", resp.StatusCode)
	}

	// Success: forward response to client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: true,
		}
	}

	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Header(k, v)
		}
	}
	c.Status(resp.StatusCode)
	c.Writer.Write(body)

	duration := time.Since(startTime)
	return &ForwardResult{
		Model:    "vendor",
		Stream:   false,
		Duration: duration,
	}, nil
}

// ──────────────────── utility functions ──────────────────────────

// extractMuleRunUsage attempts to extract token usage from a MuleRun response.
func extractMuleRunUsage(data map[string]json.RawMessage, usage *ClaudeUsage) {
	if usgRaw, ok := data["usage"]; ok {
		var u ClaudeUsage
		if json.Unmarshal(usgRaw, &u) == nil {
			usage.InputTokens += u.InputTokens
			usage.OutputTokens += u.OutputTokens
			usage.CacheCreationInputTokens += u.CacheCreationInputTokens
			usage.CacheReadInputTokens += u.CacheReadInputTokens
		}
	}
}

// truncateBody truncates a byte slice for logging purposes.
func truncateBody(body []byte, maxLen int) string {
	if len(body) <= maxLen {
		return string(body)
	}
	return string(body[:maxLen]) + "..."
}
