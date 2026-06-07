package routes

import (
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterGatewayRoutes 注册 API 网关路由（Claude/OpenAI/Gemini 兼容）
func RegisterGatewayRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	apiKeyAuth middleware.APIKeyAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	opsService *service.OpsService,
	settingService *service.SettingService,
	cfg *config.Config,
) {
	bodyLimit := middleware.RequestBodyLimit(cfg.Gateway.MaxBodySize)
	clientRequestID := middleware.ClientRequestID()
	opsErrorLogger := handler.OpsErrorLoggerMiddleware(opsService)
	endpointNorm := handler.InboundEndpointMiddleware()

	// 未分组 Key 拦截中间件（按协议格式区分错误响应）
	requireGroupAnthropic := middleware.RequireGroupAssignment(settingService, middleware.AnthropicErrorWriter)
	requireGroupGoogle := middleware.RequireGroupAssignment(settingService, middleware.GoogleErrorWriter)

	// API网关（Claude API兼容）
	gateway := r.Group("/v1")
	gateway.Use(bodyLimit)
	gateway.Use(clientRequestID)
	gateway.Use(opsErrorLogger)
	gateway.Use(endpointNorm)
	gateway.Use(gin.HandlerFunc(apiKeyAuth))
	gateway.Use(requireGroupAnthropic)
	{
		// /v1/messages: auto-route based on group platform
		gateway.POST("/messages", func(c *gin.Context) {
			switch getGroupPlatform(c) {
			case service.PlatformOpenAI:
				h.OpenAIGateway.Messages(c)
			case service.PlatformMuleRun:
				if h.MuleRunGateway != nil {
					h.MuleRunGateway.Messages(c)
				} else {
					h.Gateway.Messages(c)
				}
			default:
				h.Gateway.Messages(c)
			}
		})
		// /v1/messages/count_tokens: OpenAI groups get 404
		gateway.POST("/messages/count_tokens", func(c *gin.Context) {
			if getGroupPlatform(c) == service.PlatformOpenAI {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{
					"type": "error",
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Token counting is not supported for this platform",
					},
				})
				return
			}
			h.Gateway.CountTokens(c)
		})
		gateway.GET("/models", h.Gateway.Models)
		gateway.GET("/usage", h.Gateway.Usage)
		// OpenAI Responses API: auto-route based on group platform
		gateway.POST("/responses", func(c *gin.Context) {
			if getGroupPlatform(c) == service.PlatformOpenAI {
				h.OpenAIGateway.Responses(c)
				return
			}
			h.Gateway.Responses(c)
		})
		gateway.POST("/responses/*subpath", func(c *gin.Context) {
			if getGroupPlatform(c) == service.PlatformOpenAI {
				h.OpenAIGateway.Responses(c)
				return
			}
			h.Gateway.Responses(c)
		})
		gateway.GET("/responses", h.OpenAIGateway.ResponsesWebSocket)
		// OpenAI Chat Completions API: auto-route based on group platform
		gateway.POST("/chat/completions", func(c *gin.Context) {
			switch getGroupPlatform(c) {
			case service.PlatformOpenAI:
				h.OpenAIGateway.ChatCompletions(c)
			case service.PlatformMuleRun:
				if h.MuleRunGateway != nil {
					h.MuleRunGateway.ChatCompletions(c)
				} else {
					h.Gateway.ChatCompletions(c)
				}
			default:
				h.Gateway.ChatCompletions(c)
			}
		})
		gateway.POST("/embeddings", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformOpenAI {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Embeddings API is not supported for this platform",
					},
				})
				return
			}
			h.OpenAIGateway.Embeddings(c)
		})
		gateway.POST("/images/generations", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformOpenAI {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Images API is not supported for this platform",
					},
				})
				return
			}
			h.OpenAIGateway.Images(c)
		})
		gateway.POST("/images/edits", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformOpenAI {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Images API is not supported for this platform",
					},
				})
				return
			}
			h.OpenAIGateway.Images(c)
		})
	}

	// Gemini 原生 API 兼容层（Gemini SDK/CLI 直连）
	gemini := r.Group("/v1beta")
	gemini.Use(bodyLimit)
	gemini.Use(clientRequestID)
	gemini.Use(opsErrorLogger)
	gemini.Use(endpointNorm)
	gemini.Use(middleware.APIKeyAuthWithSubscriptionGoogle(apiKeyService, subscriptionService, cfg))
	gemini.Use(requireGroupGoogle)
	{
		gemini.GET("/models", h.Gateway.GeminiV1BetaListModels)
		gemini.GET("/models/:model", h.Gateway.GeminiV1BetaGetModel)
		// Gin treats ":" as a param marker, but Gemini uses "{model}:{action}" in the same segment.
		gemini.POST("/models/*modelAction", h.Gateway.GeminiV1BetaModels)
	}

	// OpenAI Responses API（不带v1前缀的别名）— auto-route based on group platform
	responsesHandler := func(c *gin.Context) {
		if getGroupPlatform(c) == service.PlatformOpenAI {
			h.OpenAIGateway.Responses(c)
			return
		}
		h.Gateway.Responses(c)
	}
	r.POST("/responses", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, responsesHandler)
	r.POST("/responses/*subpath", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, responsesHandler)
	r.GET("/responses", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, h.OpenAIGateway.ResponsesWebSocket)
	codexDirect := r.Group("/backend-api/codex")
	codexDirect.Use(bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic)
	{
		codexDirect.POST("/responses", responsesHandler)
		codexDirect.POST("/responses/*subpath", responsesHandler)
		codexDirect.GET("/responses", h.OpenAIGateway.ResponsesWebSocket)
	}
	// OpenAI Chat Completions API（不带v1前缀的别名）— auto-route based on group platform
	r.POST("/chat/completions", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, func(c *gin.Context) {
		switch getGroupPlatform(c) {
		case service.PlatformOpenAI:
			h.OpenAIGateway.ChatCompletions(c)
		case service.PlatformMuleRun:
			if h.MuleRunGateway != nil {
				h.MuleRunGateway.ChatCompletions(c)
			} else {
				h.Gateway.ChatCompletions(c)
			}
		default:
			h.Gateway.ChatCompletions(c)
		}
	})
	r.POST("/embeddings", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformOpenAI {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Embeddings API is not supported for this platform",
				},
			})
			return
		}
		h.OpenAIGateway.Embeddings(c)
	})
	r.POST("/images/generations", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformOpenAI {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Images API is not supported for this platform",
				},
			})
			return
		}
		h.OpenAIGateway.Images(c)
	})
	r.POST("/images/edits", bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformOpenAI {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Images API is not supported for this platform",
				},
			})
			return
		}
		h.OpenAIGateway.Images(c)
	})

	// Antigravity 模型列表
	r.GET("/antigravity/models", gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, h.Gateway.AntigravityModels)

	// Antigravity 专用路由（仅使用 antigravity 账户，不混合调度）
	antigravityV1 := r.Group("/antigravity/v1")
	antigravityV1.Use(bodyLimit)
	antigravityV1.Use(clientRequestID)
	antigravityV1.Use(opsErrorLogger)
	antigravityV1.Use(endpointNorm)
	antigravityV1.Use(middleware.ForcePlatform(service.PlatformAntigravity))
	antigravityV1.Use(gin.HandlerFunc(apiKeyAuth))
	antigravityV1.Use(requireGroupAnthropic)
	{
		antigravityV1.POST("/messages", h.Gateway.Messages)
		antigravityV1.POST("/messages/count_tokens", h.Gateway.CountTokens)
		antigravityV1.GET("/models", h.Gateway.AntigravityModels)
		antigravityV1.GET("/usage", h.Gateway.Usage)
	}

	antigravityV1Beta := r.Group("/antigravity/v1beta")
	antigravityV1Beta.Use(bodyLimit)
	antigravityV1Beta.Use(clientRequestID)
	antigravityV1Beta.Use(opsErrorLogger)
	antigravityV1Beta.Use(endpointNorm)
	antigravityV1Beta.Use(middleware.ForcePlatform(service.PlatformAntigravity))
	antigravityV1Beta.Use(middleware.APIKeyAuthWithSubscriptionGoogle(apiKeyService, subscriptionService, cfg))
	antigravityV1Beta.Use(requireGroupGoogle)
	{
		antigravityV1Beta.GET("/models", h.Gateway.GeminiV1BetaListModels)
		antigravityV1Beta.GET("/models/:model", h.Gateway.GeminiV1BetaGetModel)
		antigravityV1Beta.POST("/models/*modelAction", h.Gateway.GeminiV1BetaModels)
	}

	// ─── MuleRun 专用路由 ─────────────────────────────────────────
	// Vendor API (图片/视频/音频生成等异步任务端点)
	// 仅在配置了 MuleRunGateway 时注册
	if h.MuleRunGateway != nil {
		vendors := r.Group("/vendors")
		vendors.Use(bodyLimit)
		vendors.Use(clientRequestID)
		vendors.Use(opsErrorLogger)
		vendors.Use(endpointNorm)
		vendors.Use(gin.HandlerFunc(apiKeyAuth))
		vendors.Use(requireGroupAnthropic)
		{
			// POST /vendors/:provider/v1/:model/:action — 创建任务
			vendors.POST("/*path", func(c *gin.Context) {
				if getGroupPlatform(c) == service.PlatformMuleRun {
					h.MuleRunGateway.Vendor(c)
					return
				}
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Vendor API is only supported for MuleRun platform groups",
					},
				})
			})
			// GET /vendors/:provider/v1/:model/:action/:taskId — 轮询任务状态
			vendors.GET("/*path", func(c *gin.Context) {
				if getGroupPlatform(c) == service.PlatformMuleRun {
					h.MuleRunGateway.Vendor(c)
					return
				}
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Vendor API is only supported for MuleRun platform groups",
					},
				})
			})
		}

		// ─── Seedance 多媒体 API（仅 MuleRun 平台）────────────────────
		// 使用火山方舟 Seedance 格式作为请求/响应标准
		// POST /v1/seedance/contents/generations/tasks     — 创建多媒体生成任务
		// GET  /v1/seedance/contents/generations/tasks/:id — 查询任务状态
		seedanceMediaTasks := r.Group("/v1/seedance/contents/generations/tasks")
		seedanceMediaTasks.Use(bodyLimit)
		seedanceMediaTasks.Use(clientRequestID)
		seedanceMediaTasks.Use(opsErrorLogger)
		seedanceMediaTasks.Use(endpointNorm)
		seedanceMediaTasks.Use(gin.HandlerFunc(apiKeyAuth))
		seedanceMediaTasks.Use(requireGroupAnthropic)
		{
			seedanceMediaTasks.POST("", func(c *gin.Context) {
				if getGroupPlatform(c) == service.PlatformMuleRun {
					h.MuleRunGateway.VendorMediaCreateTask(c)
					return
				}
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Seedance Multimedia API is only supported for MuleRun platform groups",
					},
				})
			})
			seedanceMediaTasks.GET("/:id", func(c *gin.Context) {
				if getGroupPlatform(c) == service.PlatformMuleRun {
					h.MuleRunGateway.VendorMediaQueryTask(c)
					return
				}
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Seedance Multimedia API is only supported for MuleRun platform groups",
					},
				})
			})
		}

		// ─── OpenAI Image API（仅 MuleRun 平台）───────────────────────
		// 使用 OpenAI Image API 格式，内部自动轮询 MuleRun 异步任务
		// POST /v1/openai/images/generations — 生成图片
		openAIImages := r.Group("/v1/openai/images")
		openAIImages.Use(bodyLimit)
		openAIImages.Use(clientRequestID)
		openAIImages.Use(opsErrorLogger)
		openAIImages.Use(endpointNorm)
		openAIImages.Use(gin.HandlerFunc(apiKeyAuth))
		openAIImages.Use(requireGroupAnthropic)
		{
			openAIImages.POST("/generations", func(c *gin.Context) {
				if getGroupPlatform(c) == service.PlatformMuleRun {
					h.MuleRunGateway.OpenAICreateImage(c)
					return
				}
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "OpenAI Image API (MuleRun) is only supported for MuleRun platform groups",
					},
				})
			})
		}
		// ─── Google Multimedia API（仅 MuleRun 平台）──────────────────────
		// 使用 Google Gemini 格式，内部自动轮询 MuleRun 异步任务
		// POST /v1/google/models/{model}:predictLongRunning — 生成视频 (Veo)
		// POST /v1/google/models/{model}:generateContent   — 生成图片 (Nano Banana)
		googleMedia := r.Group("/v1/google/models")
		googleMedia.Use(bodyLimit)
		googleMedia.Use(clientRequestID)
		googleMedia.Use(opsErrorLogger)
		googleMedia.Use(endpointNorm)
		googleMedia.Use(gin.HandlerFunc(apiKeyAuth))
		googleMedia.Use(requireGroupAnthropic)
		{
			googleMedia.POST("/*modelAction", func(c *gin.Context) {
				if getGroupPlatform(c) != service.PlatformMuleRun {
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"type":    "not_found_error",
							"message": "Google Multimedia API is only supported for MuleRun platform groups",
						},
					})
					return
				}
				modelAction := c.Param("modelAction")
				if strings.HasSuffix(modelAction, ":predictLongRunning") {
					h.MuleRunGateway.GoogleCreateVideo(c)
				} else if strings.HasSuffix(modelAction, ":generateContent") {
					h.MuleRunGateway.GoogleCreateImage(c)
				} else {
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"type":    "not_found_error",
							"message": "Unsupported action. Use :predictLongRunning for video or :generateContent for images",
						},
					})
				}
			})
		}
	}

}

// getGroupPlatform extracts the group platform from the API Key stored in context.
func getGroupPlatform(c *gin.Context) string {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey.Group == nil {
		return ""
	}
	return apiKey.Group.Platform
}
