package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// MuleRunGatewayHandler handles API gateway requests for the MuleRun platform.
type MuleRunGatewayHandler struct {
	muleRunService        *service.MuleRunGatewayService
	vendorMediaService    *service.VendorMediaService
	gatewayService        *service.GatewayService
	usageService          *service.UsageService
	usageRecordWorkerPool *service.UsageRecordWorkerPool
	maxAccountSwitches    int
}

// NewMuleRunGatewayHandler creates a new MuleRun gateway handler.
func NewMuleRunGatewayHandler(
	muleRunService *service.MuleRunGatewayService,
	vendorMediaService *service.VendorMediaService,
	gatewayService *service.GatewayService,
	usageService *service.UsageService,
	usageRecordWorkerPool *service.UsageRecordWorkerPool,
) *MuleRunGatewayHandler {
	return &MuleRunGatewayHandler{
		muleRunService:        muleRunService,
		vendorMediaService:    vendorMediaService,
		gatewayService:        gatewayService,
		usageService:          usageService,
		usageRecordWorkerPool: usageRecordWorkerPool,
		maxAccountSwitches:    10,
	}
}

// Messages handles POST /v1/messages for MuleRun platform.
func (h *MuleRunGatewayHandler) Messages(c *gin.Context) {
	h.handleTextRequest(c, EndpointMessages)
}

// ChatCompletions handles POST /v1/chat/completions for MuleRun platform.
func (h *MuleRunGatewayHandler) ChatCompletions(c *gin.Context) {
	h.handleTextRequest(c, EndpointChatCompletions) // forward to MuleRun's OpenAI-compatible endpoint
}

// Vendor handles /vendors/* requests for MuleRun platform.
func (h *MuleRunGatewayHandler) Vendor(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key

	// Failover loop
	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			"vendor", // model name placeholder
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		fwdResult, err := h.muleRunService.ForwardVendor(c.Request.Context(), c, account)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("mulerun vendor failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			// Terminal error already written to client
			return
		}

		// Success — record usage asynchronously
		if fwdResult != nil && h.usageRecordWorkerPool != nil {
			h.recordVendorUsage(apiKey, account, fwdResult)
		}
		return
	}

	// All accounts exhausted
	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// VendorMediaCreateTask handles POST /v1/contents/generations/tasks
// Creates a multimedia generation task using Ark/Seedance standard format.
func (h *MuleRunGatewayHandler) VendorMediaCreateTask(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	arkReq, err := service.ParseArkCreateRequest(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": err.Error(),
			},
		})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key
	requestModel := arkReq.Model

	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			requestModel,
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		arkResp, err := h.vendorMediaService.CreateTask(c.Request.Context(), c, account, arkReq)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("vendor media create failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
					"model", requestModel,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "upstream_error",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusAccepted, arkResp)
		return
	}

	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// VendorMediaQueryTask handles GET /v1/contents/generations/tasks/:id
// Queries a multimedia generation task using Ark/Seedance standard format.
func (h *MuleRunGatewayHandler) VendorMediaQueryTask(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	compositeID := c.Param("id")
	if compositeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": "task ID is required",
			},
		})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key

	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			"vendor-media-query",
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		arkResp, err := h.vendorMediaService.QueryTask(c.Request.Context(), c, account, compositeID)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("vendor media query failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "upstream_error",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, arkResp)
		return
	}

	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// OpenAICreateImage handles POST /v1/openai/images/generations
// Creates an image using OpenAI Image API format, with automatic polling.
func (h *MuleRunGatewayHandler) OpenAICreateImage(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	oiReq, err := service.ParseOpenAIImageRequest(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": err.Error(),
			},
		})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key
	requestModel := oiReq.Model

	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			requestModel,
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		openAIResp, err := h.vendorMediaService.CreateImage(c.Request.Context(), c, account, oiReq)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("openai image create failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
					"model", requestModel,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "upstream_error",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, openAIResp)
		return
	}

	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// ───────────────────────── internal helpers ───────────────────────

// GoogleCreateVideo handles POST /v1/google/models/{model}:predictLongRunning
// Creates a video using Google Veo format, with automatic polling.
func (h *MuleRunGatewayHandler) GoogleCreateVideo(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	gvReq, err := service.ParseGoogleVideoRequest(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": err.Error(),
			},
		})
		return
	}

	// Extract model from URL if not in body: /v1/google/models/{model}:predictLongRunning
	if gvReq.Model == "" {
		gvReq.Model = extractModelFromAction(c)
	}
	if gvReq.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"type": "invalid_request_error", "message": "model is required"},
		})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key
	requestModel := gvReq.Model

	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			requestModel,
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		googleResp, err := h.vendorMediaService.CreateGoogleVideo(c.Request.Context(), c, account, gvReq)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("google video create failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
					"model", requestModel,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "upstream_error",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, googleResp)
		return
	}

	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// GoogleCreateImage handles POST /v1/google/models/{model}:generateContent
// Creates an image using Google Gemini format, with automatic polling and base64 encoding.
func (h *MuleRunGatewayHandler) GoogleCreateImage(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	giReq, err := service.ParseGoogleImageRequest(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"type":    "invalid_request_error",
				"message": err.Error(),
			},
		})
		return
	}

	// Extract model from URL if not in body: /v1/google/models/{model}:generateContent
	if giReq.Model == "" {
		giReq.Model = extractModelFromAction(c)
	}
	if giReq.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"type": "invalid_request_error", "message": "model is required"},
		})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key
	requestModel := giReq.Model

	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		result, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			requestModel,
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts",
				},
			})
			return
		}

		account := result.Account
		if result.ReleaseFunc != nil {
			defer result.ReleaseFunc()
		}
		if result.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), result.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		googleResp, err := h.vendorMediaService.CreateGoogleImage(c.Request.Context(), c, account, giReq)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("google image create failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
					"model", requestModel,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"type":    "upstream_error",
					"message": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, googleResp)
		return
	}

	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// extractModelFromAction extracts the model name from a Google-style URL path.
// e.g., "/veo-3.1-generate-preview:predictLongRunning" → "veo-3.1-generate-preview"
func extractModelFromAction(c *gin.Context) string {
	modelAction := c.Param("modelAction")
	if modelAction == "" {
		return ""
	}
	// Remove leading slash
	modelAction = strings.TrimPrefix(modelAction, "/")
	// Remove :action suffix
	if idx := strings.LastIndex(modelAction, ":"); idx > 0 {
		modelAction = modelAction[:idx]
	}
	return modelAction
}

// handleTextRequest is the shared logic for Messages and ChatCompletions.
func (h *MuleRunGatewayHandler) handleTextRequest(c *gin.Context, upstreamEndpoint string) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse model and stream flag
	var parsed struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request body"})
		return
	}

	groupID := apiKey.GroupID
	sessionHash := apiKey.Key
	requestModel := parsed.Model
	if requestModel == "" {
		requestModel = "default"
	}

	// Failover loop
	failedIDs := make(map[int64]struct{})
	var lastErr error

	for attempt := 0; attempt < h.maxAccountSwitches; attempt++ {
		selResult, err := h.gatewayService.SelectAccountWithLoadAwareness(
			c.Request.Context(),
			groupID,
			sessionHash,
			requestModel,
			failedIDs,
			"",
			apiKey.UserID,
		)
		if err != nil {
			slog.Warn("mulerun account selection failed",
				"error", err,
				"attempt", attempt,
			)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "pool_error",
					"message": "No available MuleRun accounts: " + err.Error(),
				},
			})
			return
		}

		account := selResult.Account
		if selResult.ReleaseFunc != nil {
			defer selResult.ReleaseFunc()
		}
		if selResult.WaitPlan != nil {
			acquired, release := h.waitForSlot(c.Request.Context(), selResult.WaitPlan)
			if !acquired {
				failedIDs[account.ID] = struct{}{}
				continue
			}
			if release != nil {
				defer release()
			}
		}

		fwdReq := &service.MuleRunForwardRequest{
			Body:     body,
			Model:    requestModel,
			Stream:   parsed.Stream,
			Endpoint: upstreamEndpoint,
		}

		fwdResult, err := h.muleRunService.Forward(c.Request.Context(), c, account, fwdReq)
		if err != nil {
			if failoverErr, ok := err.(*service.UpstreamFailoverError); ok {
				slog.Info("mulerun text failover",
					"account_id", account.ID,
					"status", failoverErr.StatusCode,
					"attempt", attempt,
					"model", requestModel,
				)
				failedIDs[account.ID] = struct{}{}
				lastErr = err

				// If the response has a body, it might be a terminal error
				if failoverErr.StatusCode == 400 || failoverErr.StatusCode == 401 {
					if len(failoverErr.ResponseBody) > 0 {
						c.Data(failoverErr.StatusCode, "application/json", failoverErr.ResponseBody)
					} else {
						c.JSON(failoverErr.StatusCode, gin.H{
							"error": gin.H{
								"type":    "upstream_error",
								"message": "MuleRun upstream error",
							},
						})
					}
					return
				}
				continue
			}
			// Terminal error — already written to response
			return
		}

		// Success — record usage asynchronously
		if fwdResult != nil && h.usageRecordWorkerPool != nil {
			h.recordTextUsage(apiKey, account, fwdResult)
		}
		return
	}

	// All accounts exhausted
	if lastErr != nil {
		if fe, ok := lastErr.(*service.UpstreamFailoverError); ok && len(fe.ResponseBody) > 0 {
			c.Data(fe.StatusCode, "application/json", fe.ResponseBody)
			return
		}
	}
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": gin.H{
			"type":    "pool_error",
			"message": "All MuleRun accounts exhausted after retries",
		},
	})
}

// waitForSlot checks if a concurrency slot is available.
// NOTE: This is currently a non-blocking check — it does NOT actually wait for a slot.
// A proper implementation requires a semaphore/channel in AccountWaitPlan.
// It only verifies that the context is not cancelled and the timeout has not elapsed.
func (h *MuleRunGatewayHandler) waitForSlot(ctx context.Context, plan *service.AccountWaitPlan) (bool, func()) {
	if plan == nil {
		return true, nil
	}
	deadline := time.After(plan.Timeout)
	select {
	case <-ctx.Done():
		return false, nil
	case <-deadline:
		return false, nil
	default:
		return true, nil
	}
}

// recordTextUsage records usage for text endpoint requests.
func (h *MuleRunGatewayHandler) recordTextUsage(
	apiKey *service.APIKey,
	account *service.Account,
	result *service.ForwardResult,
) {
	slog.Debug("mulerun text usage",
		"user_id", apiKey.UserID,
		"account_id", account.ID,
		"model", result.Model,
		"stream", result.Stream,
		"duration_ms", result.Duration.Milliseconds(),
		"input_tokens", result.Usage.InputTokens,
		"output_tokens", result.Usage.OutputTokens,
	)
}

// recordVendorUsage records usage for vendor endpoint requests.
func (h *MuleRunGatewayHandler) recordVendorUsage(
	apiKey *service.APIKey,
	account *service.Account,
	result *service.ForwardResult,
) {
	slog.Debug("mulerun vendor usage",
		"user_id", apiKey.UserID,
		"account_id", account.ID,
		"duration_ms", result.Duration.Milliseconds(),
	)
}
