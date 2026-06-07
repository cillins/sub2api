package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// ─────────────────────────────────────────────────────────────────
// VendorMediaService — handles standardized multimedia API requests
// using Ark/Seedance format as the public contract.
// ─────────────────────────────────────────────────────────────────

const vendorMediaTimeout = 120 * time.Second

// VendorMediaService processes Ark-format multimedia requests
// and forwards them to MuleRun vendor endpoints.
type VendorMediaService struct {
	httpUpstream HTTPUpstream
}

// NewVendorMediaService creates a new vendor media service.
func NewVendorMediaService(httpUpstream HTTPUpstream) *VendorMediaService {
	return &VendorMediaService{
		httpUpstream: httpUpstream,
	}
}

// CreateTask handles POST /v1/contents/generations/tasks
// It converts Ark format → MuleRun vendor format, forwards, and converts back.
func (s *VendorMediaService) CreateTask(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	arkReq *ArkCreateTaskRequest,
) (*ArkCreateTaskResponse, error) {
	cfg := LookupVendorModel(arkReq.Model)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported model: %s", arkReq.Model)
	}

	// Convert Ark request to MuleRun vendor format
	vendorBody, createPath, err := ArkToMuleRunRequest(arkReq, cfg)
	if err != nil {
		return nil, fmt.Errorf("convert request: %w", err)
	}

	// Forward to MuleRun
	respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodPost, createPath, vendorBody)
	if err != nil {
		return nil, err
	}

	if statusCode >= 400 {
		// Forward error to client as-is
		return nil, &UpstreamFailoverError{
			StatusCode:   statusCode,
			ResponseBody: respBody,
		}
	}

	// Convert MuleRun response to Ark format
	// Determine query path from create path (they share the same base)
	arkResp, err := MuleRunToArkCreateResponse(respBody, cfg.MuleRunQueryPathPrefix, arkReq.Model)
	if err != nil {
		return nil, fmt.Errorf("convert response: %w", err)
	}

	return arkResp, nil
}

// QueryTask handles GET /v1/contents/generations/tasks/:id
// It decodes the composite task ID, forwards to MuleRun, and converts back.
func (s *VendorMediaService) QueryTask(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	compositeID string,
) (*ArkQueryTaskResponse, error) {
	vendorQueryPath, upstreamTaskID, modelName, err := DecodeCompositeTaskID(compositeID)
	if err != nil {
		return nil, fmt.Errorf("decode task ID: %w", err)
	}

	// Build the full query path: vendorQueryPath + "/" + taskID
	queryPath := vendorQueryPath + "/" + upstreamTaskID

	// Forward GET to MuleRun
	respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodGet, queryPath, nil)
	if err != nil {
		return nil, err
	}

	if statusCode >= 400 {
		return nil, &UpstreamFailoverError{
			StatusCode:   statusCode,
			ResponseBody: respBody,
		}
	}

	// Convert MuleRun response to Ark format
	arkResp, err := MuleRunToArkQueryResponse(respBody, modelName)
	if err != nil {
		return nil, fmt.Errorf("convert query response: %w", err)
	}

	return arkResp, nil
}

// forwardToMuleRun sends a request to MuleRun upstream and returns the response body.
func (s *VendorMediaService) forwardToMuleRun(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	method string,
	path string,
	body []byte,
) ([]byte, int, error) {
	token, err := getMuleRunToken(account)
	if err != nil {
		return nil, 502, err
	}

	baseURL := getMuleRunBaseURL(account)
	targetURL := baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	upCtx, upCancel := context.WithTimeout(context.WithoutCancel(ctx), vendorMediaTimeout)
	defer upCancel()

	upReq, err := http.NewRequestWithContext(upCtx, method, targetURL, bodyReader)
	if err != nil {
		return nil, 502, fmt.Errorf("build upstream request: %w", err)
	}

	upReq.Header.Set("Authorization", "Bearer "+token)
	upReq.Header.Set("Content-Type", "application/json")
	upReq.Header.Set("Accept", "application/json")

	// Pass through user agent
	if ua := c.GetHeader("User-Agent"); ua != "" {
		upReq.Header.Set("User-Agent", ua)
	}

	proxyURL := ""
	if account.Proxy != nil {
		proxyURL = account.Proxy.URL()
	}

	resp, err := s.httpUpstream.Do(upReq, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		slog.Warn("vendor media upstream request failed",
			"account_id", account.ID,
			"method", method,
			"path", path,
			"error", err,
		)
		return nil, 502, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: false,
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 502, &UpstreamFailoverError{
			StatusCode:             502,
			RetryableOnSameAccount: true,
		}
	}

	return respBody, resp.StatusCode, nil
}

// CreateImage handles POST /v1/openai/images/generations
// It converts OpenAI Image API format → MuleRun vendor format, forwards,
// polls until task completes, fetches images, and returns OpenAI format response.
func (s *VendorMediaService) CreateImage(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	oiReq *OpenAIImageRequest,
) (*OpenAIImageResponse, error) {
	cfg := LookupVendorModel(oiReq.Model)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported model: %s", oiReq.Model)
	}
	if !cfg.HasTextToImage {
		return nil, fmt.Errorf("model %s does not support image generation", oiReq.Model)
	}

	// Convert OpenAI request to MuleRun vendor format
	vendorBody, createPath, err := OpenAIToMuleRunRequest(oiReq, cfg)
	if err != nil {
		return nil, fmt.Errorf("convert request: %w", err)
	}

	// Forward POST to MuleRun to create task
	respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodPost, createPath, vendorBody)
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, &UpstreamFailoverError{StatusCode: statusCode, ResponseBody: respBody}
	}

	// Parse task creation response
	var mrResp MuleRunVendorTaskResponse
	if err := json.Unmarshal(respBody, &mrResp); err != nil {
		return nil, fmt.Errorf("parse mulerun create response: %w", err)
	}
	upstreamTaskID := ""
	if mrResp.TaskInfo != nil {
		upstreamTaskID = mrResp.TaskInfo.ID
	}
	if upstreamTaskID == "" {
		upstreamTaskID = mrResp.ID
	}
	if upstreamTaskID == "" {
		return nil, fmt.Errorf("no task ID in mulerun response")
	}

	// Poll the task until it completes or times out
	queryPath := cfg.MuleRunQueryPathPrefix + "/" + upstreamTaskID
	imageURLs, err := s.pollTaskUntilComplete(ctx, c, account, queryPath)
	if err != nil {
		return nil, err
	}

	// Build OpenAI response
	openAIResp := MuleRunToOpenAIImageResponse(imageURLs, oiReq.Prompt, oiReq.ResponseFormat)

	// If response_format is b64_json, fetch URLs and convert to base64
	if oiReq.ResponseFormat == "b64_json" {
		for i := range openAIResp.Data {
			if openAIResp.Data[i].URL != "" {
				b64, fetchErr := s.fetchImageAsBase64(ctx, openAIResp.Data[i].URL, account)
				if fetchErr != nil {
					slog.Warn("failed to fetch image for base64 encoding",
						"url", openAIResp.Data[i].URL,
						"error", fetchErr,
					)
					// Fall back to returning URL
					continue
				}
				openAIResp.Data[i].B64JSON = b64
				openAIResp.Data[i].URL = ""
			}
		}
	}

	return openAIResp, nil
}

// pollTaskUntilComplete polls a MuleRun task until it completes or times out.
// Returns the image/video URLs from the completed task.
func (s *VendorMediaService) pollTaskUntilComplete(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	queryPath string,
) ([]string, error) {
	const (
		maxPollDuration = 5 * time.Minute
		pollInterval    = 3 * time.Second
	)

	deadline := time.Now().Add(maxPollDuration)
	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("task polling timeout after %v", maxPollDuration)
		}

		// Check context cancellation before each poll
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodGet, queryPath, nil)
		if err != nil {
			return nil, err
		}
		if statusCode >= 400 {
			return nil, &UpstreamFailoverError{StatusCode: statusCode, ResponseBody: respBody}
		}

		var mrResp MuleRunVendorTaskResponse
		if err := json.Unmarshal(respBody, &mrResp); err != nil {
			return nil, fmt.Errorf("parse mulerun query response: %w", err)
		}

		status := ""
		if mrResp.TaskInfo != nil {
			status = mrResp.TaskInfo.Status
		} else {
			status = mrResp.Status
		}

		switch status {
		case MuleRunStatusCompleted:
			// Extract output URLs (videos or images field)
			if len(mrResp.Videos) > 0 {
				return mrResp.Videos, nil
			}
			// Some image models might use a different field — check raw response
			var raw map[string]any
			json.Unmarshal(respBody, &raw)
			if images, ok := raw["images"].([]any); ok {
				urls := make([]string, 0, len(images))
				for _, img := range images {
					if s, ok := img.(string); ok {
						urls = append(urls, s)
					}
				}
				if len(urls) > 0 {
					return urls, nil
				}
			}
			// Fallback: try "output" or "result_url"
			if outputURL, ok := raw["result_url"].(string); ok && outputURL != "" {
				return []string{outputURL}, nil
			}
			return nil, fmt.Errorf("task completed but no output URLs found")

		case MuleRunStatusFailed:
			errMsg := "task failed"
			var raw map[string]any
			json.Unmarshal(respBody, &raw)
			if msg, ok := raw["error"].(string); ok {
				errMsg = msg
			}
			return nil, fmt.Errorf("upstream task failed: %s", errMsg)

		case MuleRunStatusPending, MuleRunStatusProcessing:
			// Continue polling — respect context cancellation during wait
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(pollInterval):
			}

		default:
			// Unknown status — treat as terminal to avoid infinite polling
			slog.Warn("task returned unknown status, treating as failed",
				"status", status,
				"query_path", queryPath,
			)
			return nil, fmt.Errorf("task returned unknown status: %s", status)
		}
	}
}

// fetchImageAsBase64 downloads an image from a URL and returns it as base64.
// It uses the account's proxy if configured.
func (s *VendorMediaService) fetchImageAsBase64(ctx context.Context, imageURL string, account *Account) (string, error) {
	fetchCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(fetchCtx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("build fetch request: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	if account != nil && account.Proxy != nil {
		if proxyParsed, parseErr := url.Parse(account.Proxy.URL()); parseErr == nil && proxyParsed.Scheme != "" {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyParsed),
			}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch image failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read image body: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// ParseArkCreateRequest parses and validates an Ark create task request.
func ParseArkCreateRequest(body []byte) (*ArkCreateTaskRequest, error) {
	var req ArkCreateTaskRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if len(req.Content) == 0 {
		return nil, fmt.Errorf("content is required")
	}
	// Validate content items
	hasText := false
	for _, c := range req.Content {
		if c.Type == ArkContentTypeText && c.Text != "" {
			hasText = true
		}
	}
	// At least one text or image content is needed for most models
	if !hasText {
		// Check if there's at least an image
		hasImage := false
		for _, c := range req.Content {
			if c.Type == ArkContentTypeImageURL && c.ImageURL != nil {
				hasImage = true
			}
		}
		if !hasImage {
			return nil, fmt.Errorf("at least one text or image content is required")
		}
	}
	return &req, nil
}

// CreateGoogleVideo handles POST /v1/google/models/{model}:predictLongRunning
// It converts Google Veo format → MuleRun vendor format, forwards,
// polls until task completes, and returns Google operation response format.
func (s *VendorMediaService) CreateGoogleVideo(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	gvReq *GoogleVideoRequest,
) (*GoogleVideoOperationResponse, error) {
	cfg := LookupVendorModel(gvReq.Model)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported model: %s", gvReq.Model)
	}
	if !cfg.HasTextToVideo && !cfg.HasImageToVideo {
		return nil, fmt.Errorf("model %s does not support video generation", gvReq.Model)
	}

	// Convert Google request to MuleRun vendor format
	vendorBody, createPath, err := GoogleVideoToMuleRunRequest(gvReq, cfg)
	if err != nil {
		return nil, fmt.Errorf("convert request: %w", err)
	}

	// Forward POST to MuleRun to create task
	respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodPost, createPath, vendorBody)
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, &UpstreamFailoverError{StatusCode: statusCode, ResponseBody: respBody}
	}

	// Parse task creation response
	var mrResp MuleRunVendorTaskResponse
	if err := json.Unmarshal(respBody, &mrResp); err != nil {
		return nil, fmt.Errorf("parse mulerun create response: %w", err)
	}
	upstreamTaskID := ""
	if mrResp.TaskInfo != nil {
		upstreamTaskID = mrResp.TaskInfo.ID
	}
	if upstreamTaskID == "" {
		upstreamTaskID = mrResp.ID
	}
	if upstreamTaskID == "" {
		return nil, fmt.Errorf("no task ID in mulerun response")
	}

	// Poll the task until it completes or times out
	queryPath := cfg.MuleRunQueryPathPrefix + "/" + upstreamTaskID
	videoURLs, err := s.pollTaskUntilComplete(ctx, c, account, queryPath)
	if err != nil {
		return nil, err
	}

	// Build Google operation response
	operationName := "operations/" + upstreamTaskID
	return MuleRunToGoogleVideoOperation(videoURLs, operationName), nil
}

// CreateGoogleImage handles POST /v1/google/models/{model}:generateContent
// It converts Google Gemini format → MuleRun vendor format, forwards,
// polls until task completes, fetches images as base64, and returns Gemini response.
func (s *VendorMediaService) CreateGoogleImage(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	giReq *GoogleImageRequest,
) (*GoogleImageResponse, error) {
	cfg := LookupVendorModel(giReq.Model)
	if cfg == nil {
		return nil, fmt.Errorf("unsupported model: %s", giReq.Model)
	}
	if !cfg.HasTextToImage {
		return nil, fmt.Errorf("model %s does not support image generation", giReq.Model)
	}

	// Convert Google request to MuleRun vendor format
	vendorBody, createPath, err := GoogleImageToMuleRunRequest(giReq, cfg)
	if err != nil {
		return nil, fmt.Errorf("convert request: %w", err)
	}

	// Forward POST to MuleRun to create task
	respBody, statusCode, err := s.forwardToMuleRun(ctx, c, account, http.MethodPost, createPath, vendorBody)
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, &UpstreamFailoverError{StatusCode: statusCode, ResponseBody: respBody}
	}

	// Parse task creation response
	var mrResp MuleRunVendorTaskResponse
	if err := json.Unmarshal(respBody, &mrResp); err != nil {
		return nil, fmt.Errorf("parse mulerun create response: %w", err)
	}
	upstreamTaskID := ""
	if mrResp.TaskInfo != nil {
		upstreamTaskID = mrResp.TaskInfo.ID
	}
	if upstreamTaskID == "" {
		upstreamTaskID = mrResp.ID
	}
	if upstreamTaskID == "" {
		return nil, fmt.Errorf("no task ID in mulerun response")
	}

	// Poll the task until it completes or times out
	queryPath := cfg.MuleRunQueryPathPrefix + "/" + upstreamTaskID
	imageURLs, err := s.pollTaskUntilComplete(ctx, c, account, queryPath)
	if err != nil {
		return nil, err
	}

	// Fetch images as base64 (Google format always returns inline_data)
	base64Images := make([]string, 0, len(imageURLs))
	for _, url := range imageURLs {
		b64, fetchErr := s.fetchImageAsBase64(ctx, url, account)
		if fetchErr != nil {
			slog.Warn("failed to fetch image for base64 encoding",
				"url", url,
				"error", fetchErr,
			)
			continue
		}
		base64Images = append(base64Images, b64)
	}

	// Build Google response
	return MuleRunToGoogleImageResponse(imageURLs, base64Images), nil
}
