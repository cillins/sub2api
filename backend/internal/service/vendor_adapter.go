package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ─────────────────────────────────────────────────────────────────
// VendorAdapter — model registry and format conversion
// ─────────────────────────────────────────────────────────────────

// VendorType identifies the vendor API family.
type VendorType string

const (
	VendorTypeSeedance VendorType = "seedance"
	VendorTypeGoogle   VendorType = "google"
	VendorTypeOpenAI   VendorType = "openai"
	VendorTypeKlingAI  VendorType = "klingai"
	VendorTypeAlibaba  VendorType = "alibaba"
	VendorTypeMidjourney VendorType = "midjourney"
	VendorTypeMiniMax  VendorType = "minimax"
)

// VendorModelConfig maps a user-facing model name to a MuleRun vendor endpoint.
type VendorModelConfig struct {
	// ModelName is the user-facing model name (e.g. "seedance-2.0").
	ModelName string
	// VendorType identifies the vendor family.
	VendorType VendorType
	// MuleRunCreatePath is the MuleRun POST endpoint for task creation.
	MuleRunCreatePath string
	// MuleRunQueryPathPrefix is the MuleRun GET prefix for task polling.
	// The full query path is: MuleRunQueryPathPrefix + "/" + taskID
	MuleRunQueryPathPrefix string
	// HasImageToVideo indicates if the model supports image-to-video.
	HasImageToVideo bool
	// HasTextToVideo indicates if the model supports text-to-video.
	HasTextToVideo bool
	// HasTextToImage for image generation models.
	HasTextToImage bool
	// HasTextToMusic for music generation models.
	HasTextToMusic bool
	// HasTextToSpeech for TTS models.
	HasTextToSpeech bool
	// OutputType: "video", "image", "audio"
	OutputType string
}

// vendorModelRegistry maps model names to their MuleRun configurations.
var vendorModelRegistry = map[string]*VendorModelConfig{
	// ── ByteDance Seedance ──
	"seedance-2.0": {
		ModelName:              "seedance-2.0",
		VendorType:             VendorTypeSeedance,
		MuleRunCreatePath:      "/vendors/bytedance/v1/seedance-2.0/text-to-video/generation",
		MuleRunQueryPathPrefix: "/vendors/bytedance/v1/seedance-2.0/text-to-video/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"seedance-2.0-fast": {
		ModelName:              "seedance-2.0-fast",
		VendorType:             VendorTypeSeedance,
		MuleRunCreatePath:      "/vendors/bytedance/v1/seedance-2.0-fast/text-to-video/generation",
		MuleRunQueryPathPrefix: "/vendors/bytedance/v1/seedance-2.0-fast/text-to-video/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	// ── Google ──
	"veo3": {
		ModelName:              "veo3",
		VendorType:             VendorTypeGoogle,
		MuleRunCreatePath:      "/vendors/google/v1/veo/generation",
		MuleRunQueryPathPrefix: "/vendors/google/v1/veo/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"nano-banana-pro": {
		ModelName:              "nano-banana-pro",
		VendorType:             VendorTypeGoogle,
		MuleRunCreatePath:      "/vendors/google/v1/nano-banana-pro/generation",
		MuleRunQueryPathPrefix: "/vendors/google/v1/nano-banana-pro/generation",
		HasTextToImage:         true,
		OutputType:             "image",
	},
	"nano-banana-2": {
		ModelName:              "nano-banana-2",
		VendorType:             VendorTypeGoogle,
		MuleRunCreatePath:      "/vendors/google/v1/nano-banana-2/generation",
		MuleRunQueryPathPrefix: "/vendors/google/v1/nano-banana-2/generation",
		HasTextToImage:         true,
		OutputType:             "image",
	},
	"nano-banana": {
		ModelName:              "nano-banana",
		VendorType:             VendorTypeGoogle,
		MuleRunCreatePath:      "/vendors/google/v1/nano-banana/generation",
		MuleRunQueryPathPrefix: "/vendors/google/v1/nano-banana/generation",
		HasTextToImage:         true,
		OutputType:             "image",
	},
	// ── OpenAI ──
	"sora": {
		ModelName:              "sora",
		VendorType:             VendorTypeOpenAI,
		MuleRunCreatePath:      "/vendors/openai/v1/sora/generation",
		MuleRunQueryPathPrefix: "/vendors/openai/v1/sora/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"sora-2": {
		ModelName:              "sora-2",
		VendorType:             VendorTypeOpenAI,
		MuleRunCreatePath:      "/vendors/openai/v1/sora-2/generation",
		MuleRunQueryPathPrefix: "/vendors/openai/v1/sora-2/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"gpt-image-2": {
		ModelName:              "gpt-image-2",
		VendorType:             VendorTypeOpenAI,
		MuleRunCreatePath:      "/vendors/openai/v1/gpt-image-2/generation",
		MuleRunQueryPathPrefix: "/vendors/openai/v1/gpt-image-2/generation",
		HasTextToImage:         true,
		OutputType:             "image",
	},
	// ── KlingAI ──
	"kling-v3": {
		ModelName:              "kling-v3",
		VendorType:             VendorTypeKlingAI,
		MuleRunCreatePath:      "/vendors/klingai/v1/kling-v3/text-to-video/generation",
		MuleRunQueryPathPrefix: "/vendors/klingai/v1/kling-v3/text-to-video/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"kling-v3-omni": {
		ModelName:              "kling-v3-omni",
		VendorType:             VendorTypeKlingAI,
		MuleRunCreatePath:      "/vendors/klingai/v1/kling-v3-omni/text-to-video/generation",
		MuleRunQueryPathPrefix: "/vendors/klingai/v1/kling-v3-omni/text-to-video/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	// ── Alibaba ──
	"wan2.6-t2v": {
		ModelName:              "wan2.6-t2v",
		VendorType:             VendorTypeAlibaba,
		MuleRunCreatePath:      "/vendors/alibaba/v1/wan2.6-t2v/generation",
		MuleRunQueryPathPrefix: "/vendors/alibaba/v1/wan2.6-t2v/generation",
		HasTextToVideo:         true,
		OutputType:             "video",
	},
	"wan2.6-i2v": {
		ModelName:              "wan2.6-i2v",
		VendorType:             VendorTypeAlibaba,
		MuleRunCreatePath:      "/vendors/alibaba/v1/wan2.6-i2v/generation",
		MuleRunQueryPathPrefix: "/vendors/alibaba/v1/wan2.6-i2v/generation",
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	"happy-horse-1.0": {
		ModelName:              "happy-horse-1.0",
		VendorType:             VendorTypeAlibaba,
		MuleRunCreatePath:      "/vendors/alibaba/v1/happy-horse-1-0-t2v/generation",
		MuleRunQueryPathPrefix: "/vendors/alibaba/v1/happy-horse-1-0-t2v/generation",
		HasTextToVideo:         true,
		HasImageToVideo:        true,
		OutputType:             "video",
	},
	// ── Midjourney ──
	"midjourney": {
		ModelName:              "midjourney",
		VendorType:             VendorTypeMidjourney,
		MuleRunCreatePath:      "/vendors/midjourney/v1/tob/diffusion",
		MuleRunQueryPathPrefix: "/vendors/midjourney/v1/tob/diffusion",
		HasTextToImage:         true,
		OutputType:             "image",
	},
	// ── MiniMax ──
	"minimax-music": {
		ModelName:              "minimax-music",
		VendorType:             VendorTypeMiniMax,
		MuleRunCreatePath:      "/vendors/minimax/v1/music-2.5/text-to-music/generation",
		MuleRunQueryPathPrefix: "/vendors/minimax/v1/music-2.5/text-to-music/generation",
		HasTextToMusic:         true,
		OutputType:             "audio",
	},
}

// LookupVendorModel finds a vendor model config by name.
func LookupVendorModel(modelName string) *VendorModelConfig {
	if cfg, ok := vendorModelRegistry[modelName]; ok {
		return cfg
	}
	// Try case-insensitive match
	lower := strings.ToLower(modelName)
	for k, v := range vendorModelRegistry {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return nil
}

// ListVendorModels returns all registered model names.
func ListVendorModels() []string {
	names := make([]string, 0, len(vendorModelRegistry))
	for k := range vendorModelRegistry {
		names = append(names, k)
	}
	return names
}

// ─────────────────────────────────────────────────────────────────
// Composite Task ID encoding/decoding
// Encodes vendor path + upstream task ID into a single opaque ID
// so that GET requests can be routed without external state.
// Format: base64url(JSON{v: vendorQueryPath, t: upstreamTaskID, m: model})
// ─────────────────────────────────────────────────────────────────

type compositeTaskPayload struct {
	VendorQueryPath string `json:"v"`
	TaskID          string `json:"t"`
	Model           string `json:"m"`
}

// EncodeCompositeTaskID creates an opaque task ID from vendor info + upstream ID.
func EncodeCompositeTaskID(vendorQueryPath, upstreamTaskID, modelName string) string {
	payload := compositeTaskPayload{
		VendorQueryPath: vendorQueryPath,
		TaskID:          upstreamTaskID,
		Model:           modelName,
	}
	data, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(data)
}

// DecodeCompositeTaskID extracts vendor path + upstream task ID from an opaque ID.
func DecodeCompositeTaskID(compositeID string) (vendorQueryPath, upstreamTaskID, modelName string, err error) {
	data, err := base64.RawURLEncoding.DecodeString(compositeID)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid task ID format: %w", err)
	}
	var payload compositeTaskPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", "", "", fmt.Errorf("invalid task ID content: %w", err)
	}
	if payload.VendorQueryPath == "" || payload.TaskID == "" {
		return "", "", "", fmt.Errorf("incomplete task ID payload")
	}
	return payload.VendorQueryPath, payload.TaskID, payload.Model, nil
}

// ─────────────────────────────────────────────────────────────────
// OpenAI → MuleRun request conversion (Image API)
// ─────────────────────────────────────────────────────────────────

// OpenAIToMuleRunRequest converts an OpenAI Image API request to a MuleRun vendor request body.
func OpenAIToMuleRunRequest(req *OpenAIImageRequest, cfg *VendorModelConfig) (body []byte, createPath string, err error) {
	if req.Prompt == "" {
		return nil, "", fmt.Errorf("prompt is required")
	}

	bodyMap := map[string]any{
		"prompt": req.Prompt,
	}

	// Map size to resolution/aspect_ratio
	if req.Size != "" {
		bodyMap["resolution"] = req.Size
	}

	// Map quality
	if req.Quality != "" {
		bodyMap["quality"] = req.Quality
	}

	// Map output format
	if req.OutputFormat != "" {
		bodyMap["output_format"] = req.OutputFormat
	}

	// Map output compression
	if req.OutputCompression != nil {
		bodyMap["output_compression"] = *req.OutputCompression
	}

	// Map background
	if req.Background != "" {
		bodyMap["background"] = req.Background
	}

	// Map n
	if req.N != nil && *req.N > 1 {
		bodyMap["n"] = *req.N
	}

	createPath = cfg.MuleRunCreatePath
	data, _ := json.Marshal(bodyMap)
	return data, createPath, nil
}

// MuleRunToOpenAIImageResponse converts a completed MuleRun task response
// to the OpenAI Image API format.
// imageURLs are the output URLs from the MuleRun vendor response.
// If responseFormat is "b64_json", the caller should fetch the URL and encode it.
func MuleRunToOpenAIImageResponse(imageURLs []string, revisedPrompt string, responseFormat string) *OpenAIImageResponse {
	resp := &OpenAIImageResponse{
		Created: time.Now().Unix(),
		Data:    make([]OpenAIImageData, 0, len(imageURLs)),
	}

	for _, url := range imageURLs {
		item := OpenAIImageData{
			RevisedPrompt: revisedPrompt,
		}
		if responseFormat == "url" {
			item.URL = url
		} else {
			// Default: b64_json — caller will fill this after fetching
			item.URL = url // Store URL temporarily for caller to fetch
		}
		resp.Data = append(resp.Data, item)
	}

	if len(resp.Data) == 0 {
		// Return at least one item with empty data
		resp.Data = append(resp.Data, OpenAIImageData{})
	}

	return resp
}

// ParseOpenAIImageRequest parses and validates an OpenAI Image API request.
func ParseOpenAIImageRequest(body []byte) (*OpenAIImageRequest, error) {
	var req OpenAIImageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if req.Prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}
	// Default response_format to b64_json (OpenAI default)
	if req.ResponseFormat == "" {
		req.ResponseFormat = "b64_json"
	}
	return &req, nil
}

// ─────────────────────────────────────────────────────────────────
// Google Video → MuleRun request conversion
// ─────────────────────────────────────────────────────────────────

// ParseGoogleVideoRequest parses and validates a Google video generation request.
func ParseGoogleVideoRequest(body []byte) (*GoogleVideoRequest, error) {
	var req GoogleVideoRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	// Model may come from URL path instead of body — don't require it here
	if len(req.Instances) == 0 {
		return nil, fmt.Errorf("instances is required")
	}
	if req.Instances[0].Prompt == "" {
		return nil, fmt.Errorf("prompt is required in instances[0]")
	}
	return &req, nil
}

// GoogleVideoToMuleRunRequest converts a Google Veo video request to a MuleRun vendor request body.
func GoogleVideoToMuleRunRequest(req *GoogleVideoRequest, cfg *VendorModelConfig) (body []byte, createPath string, err error) {
	inst := req.Instances[0]
	bodyMap := map[string]any{
		"prompt": inst.Prompt,
	}

	// Map first frame image (base64 inline data → URL not needed, MuleRun expects URL or base64)
	// For now, pass through the inline data as "image" if present
	// Note: MuleRun vendor API expects URL strings, but Google format uses base64.
	// We'll pass the base64 data and let the vendor handle it.

	// Map parameters
	if req.Parameters != nil {
		if req.Parameters.AspectRatio != "" {
			bodyMap["aspect_ratio"] = req.Parameters.AspectRatio
		}
		if req.Parameters.Resolution != "" {
			bodyMap["resolution"] = req.Parameters.Resolution
		}
		if req.Parameters.DurationSeconds != "" {
			// Convert string duration to int for MuleRun
			var dur int
			if _, scanErr := fmt.Sscanf(req.Parameters.DurationSeconds, "%d", &dur); scanErr == nil {
				bodyMap["duration"] = dur
			}
		}
		if req.Parameters.PersonGeneration != "" {
			bodyMap["person_generation"] = req.Parameters.PersonGeneration
		}
		if req.Parameters.Seed != nil {
			bodyMap["seed"] = *req.Parameters.Seed
		}
	}

	createPath = cfg.MuleRunCreatePath
	data, _ := json.Marshal(bodyMap)
	return data, createPath, nil
}

// MuleRunToGoogleVideoOperation converts MuleRun video URLs to a Google video operation response.
func MuleRunToGoogleVideoOperation(videoURLs []string, operationName string) *GoogleVideoOperationResponse {
	resp := &GoogleVideoOperationResponse{
		Name: operationName,
		Done: true,
		Response: &GoogleVideoOperationResult{
			GenerateVideoResponse: &GoogleGenerateVideoResponse{
				GeneratedSamples: make([]GoogleVideoSample, 0, len(videoURLs)),
			},
		},
	}

	for _, url := range videoURLs {
		resp.Response.GenerateVideoResponse.GeneratedSamples = append(
			resp.Response.GenerateVideoResponse.GeneratedSamples,
			GoogleVideoSample{
				Video: &GoogleVideoFile{URI: url},
			},
		)
	}

	if len(resp.Response.GenerateVideoResponse.GeneratedSamples) == 0 {
		resp.Response.GenerateVideoResponse.GeneratedSamples = append(
			resp.Response.GenerateVideoResponse.GeneratedSamples,
			GoogleVideoSample{},
		)
	}

	return resp
}

// ─────────────────────────────────────────────────────────────────
// Google Image → MuleRun request conversion
// ─────────────────────────────────────────────────────────────────

// ParseGoogleImageRequest parses and validates a Google image generation request.
func ParseGoogleImageRequest(body []byte) (*GoogleImageRequest, error) {
	var req GoogleImageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	// Model may come from URL path instead of body — don't require it here
	if len(req.Contents) == 0 {
		return nil, fmt.Errorf("contents is required")
	}
	// Extract text prompt to validate
	hasText := false
	for _, c := range req.Contents {
		for _, p := range c.Parts {
			if p.Text != "" {
				hasText = true
				break
			}
		}
	}
	if !hasText {
		return nil, fmt.Errorf("at least one text part is required in contents")
	}
	return &req, nil
}

// GoogleImageToMuleRunRequest converts a Google image request to a MuleRun vendor request body.
func GoogleImageToMuleRunRequest(req *GoogleImageRequest, cfg *VendorModelConfig) (body []byte, createPath string, err error) {
	// Extract text prompt from contents
	var prompt string
	for _, c := range req.Contents {
		for _, p := range c.Parts {
			if p.Text != "" {
				prompt = p.Text
				break
			}
		}
	}

	bodyMap := map[string]any{
		"prompt": prompt,
	}

	createPath = cfg.MuleRunCreatePath
	data, _ := json.Marshal(bodyMap)
	return data, createPath, nil
}

// MuleRunToGoogleImageResponse converts MuleRun image URLs/base64 to Google Gemini image response format.
func MuleRunToGoogleImageResponse(imageURLs []string, base64Images []string) *GoogleImageResponse {
	resp := &GoogleImageResponse{
		Candidates: []GoogleCandidate{
			{
				Content: GoogleContentOutput{
					Parts: make([]GooglePart, 0),
					Role:  "model",
				},
			},
		},
	}

	// If we have base64 images, use inline_data
	if len(base64Images) > 0 {
		for _, b64 := range base64Images {
			resp.Candidates[0].Content.Parts = append(resp.Candidates[0].Content.Parts, GooglePart{
				InlineData: &GoogleBlob{
					MimeType: "image/png",
					Data:     b64,
				},
			})
		}
	} else if len(imageURLs) > 0 {
		// Fall back to URLs — return as text with URIs (Google doesn't have URL format,
		// but we include the URL in text for reference)
		for _, url := range imageURLs {
			resp.Candidates[0].Content.Parts = append(resp.Candidates[0].Content.Parts, GooglePart{
				Text: url,
			})
		}
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		resp.Candidates[0].Content.Parts = append(resp.Candidates[0].Content.Parts, GooglePart{})
	}

	return resp
}

// ─────────────────────────────────────────────────────────────────
// Ark → MuleRun request conversion
// ─────────────────────────────────────────────────────────────────

// ArkToMuleRunRequest converts an Ark-format request to a MuleRun vendor request body.
// It also determines the correct MuleRun create path (t2v vs i2v).
func ArkToMuleRunRequest(arkReq *ArkCreateTaskRequest, cfg *VendorModelConfig) (body []byte, createPath string, err error) {
	// Extract text prompt and media from content array
	var prompt string
	var firstFrameURL, lastFrameURL string
	var refImages []string
	var refVideos []string
	var refAudios []string

	for _, c := range arkReq.Content {
		switch c.Type {
		case ArkContentTypeText:
			prompt = c.Text
		case ArkContentTypeImageURL:
			if c.ImageURL != nil {
				switch c.Role {
				case ArkRoleFirstFrame, "":
					firstFrameURL = c.ImageURL.URL
				case ArkRoleLastFrame:
					lastFrameURL = c.ImageURL.URL
				case ArkRoleReferenceImage:
					refImages = append(refImages, c.ImageURL.URL)
				}
			}
		case ArkContentTypeVideoURL:
			if c.VideoURL != nil {
				refVideos = append(refVideos, c.VideoURL.URL)
			}
		case ArkContentTypeAudioURL:
			if c.AudioURL != nil {
				refAudios = append(refAudios, c.AudioURL.URL)
			}
		}
	}

	// Determine if this is i2v or t2v
	hasImage := firstFrameURL != "" || len(refImages) > 0
	createPath = cfg.MuleRunCreatePath

	// For models with separate t2v/i2v endpoints, switch path
	if hasImage && cfg.HasImageToVideo {
		createPath = muleRunI2VPath(cfg)
	}

	// Build vendor-specific request body based on VendorType
	switch cfg.VendorType {
	case VendorTypeSeedance:
		return buildSeedanceVendorBody(arkReq, prompt, firstFrameURL, lastFrameURL), createPath, nil
	case VendorTypeGoogle:
		return buildGoogleVendorBody(arkReq, prompt, firstFrameURL, lastFrameURL, refImages), createPath, nil
	case VendorTypeOpenAI:
		return buildOpenAIVendorBody(arkReq, prompt, firstFrameURL), createPath, nil
	case VendorTypeKlingAI:
		return buildKlingVendorBody(arkReq, prompt, firstFrameURL, lastFrameURL), createPath, nil
	case VendorTypeAlibaba:
		return buildAlibabaVendorBody(arkReq, prompt, firstFrameURL), createPath, nil
	case VendorTypeMidjourney:
		return buildMidjourneyVendorBody(arkReq, prompt), createPath, nil
	case VendorTypeMiniMax:
		return buildMiniMaxVendorBody(arkReq, prompt), createPath, nil
	default:
		return nil, "", fmt.Errorf("unsupported vendor type: %s", cfg.VendorType)
	}
}

// muleRunI2VPath converts a t2v create path to an i2v path.
func muleRunI2VPath(cfg *VendorModelConfig) string {
	p := cfg.MuleRunCreatePath
	switch cfg.VendorType {
	case VendorTypeSeedance:
		return strings.Replace(p, "/text-to-video/", "/image-to-video/", 1)
	case VendorTypeKlingAI:
		return strings.Replace(p, "/text-to-video/", "/image-to-video/", 1)
	case VendorTypeAlibaba:
		// Happy horse has separate i2v endpoint
		if strings.Contains(p, "happy-horse") {
			return strings.Replace(p, "happy-horse-1-0-t2v", "happy-horse-1-0-i2v", 1)
		}
		// wan2.6-i2v is already i2v
		return p
	default:
		// Google, OpenAI use the same endpoint for both
		return p
	}
}

// ─── Vendor-specific body builders ───

func buildSeedanceVendorBody(arkReq *ArkCreateTaskRequest, prompt, firstFrame, lastFrame string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	if firstFrame != "" {
		body["image"] = firstFrame
	}
	if lastFrame != "" {
		body["last_frame_image"] = lastFrame
	}
	applyCommonParams(arkReq, body)
	b, _ := json.Marshal(body)
	return b
}

func buildGoogleVendorBody(arkReq *ArkCreateTaskRequest, prompt, firstFrame, lastFrame string, refImages []string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	if firstFrame != "" {
		body["image"] = firstFrame
	}
	if lastFrame != "" {
		body["last_frame"] = lastFrame
	}
	if len(refImages) > 0 {
		body["reference_images"] = refImages
	}
	if arkReq.Resolution != "" {
		body["resolution"] = arkReq.Resolution
	}
	if arkReq.Ratio != "" {
		body["aspect_ratio"] = arkReq.Ratio
	}
	if arkReq.Duration != nil {
		body["duration"] = *arkReq.Duration
	}
	b, _ := json.Marshal(body)
	return b
}

func buildOpenAIVendorBody(arkReq *ArkCreateTaskRequest, prompt, firstFrame string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	if firstFrame != "" {
		body["image"] = firstFrame
	}
	if arkReq.Resolution != "" {
		body["resolution"] = arkReq.Resolution
	}
	if arkReq.Ratio != "" {
		body["aspect_ratio"] = arkReq.Ratio
	}
	if arkReq.Duration != nil {
		body["duration"] = *arkReq.Duration
	}
	b, _ := json.Marshal(body)
	return b
}

func buildKlingVendorBody(arkReq *ArkCreateTaskRequest, prompt, firstFrame, lastFrame string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	if firstFrame != "" {
		body["image"] = firstFrame
	}
	if lastFrame != "" {
		body["image_tail"] = lastFrame
	}
	applyCommonParams(arkReq, body)
	b, _ := json.Marshal(body)
	return b
}

func buildAlibabaVendorBody(arkReq *ArkCreateTaskRequest, prompt, firstFrame string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	if firstFrame != "" {
		body["image"] = firstFrame
	}
	applyCommonParams(arkReq, body)
	b, _ := json.Marshal(body)
	return b
}

func buildMidjourneyVendorBody(arkReq *ArkCreateTaskRequest, prompt string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	b, _ := json.Marshal(body)
	return b
}

func buildMiniMaxVendorBody(arkReq *ArkCreateTaskRequest, prompt string) []byte {
	body := map[string]any{
		"prompt": prompt,
	}
	b, _ := json.Marshal(body)
	return b
}

// applyCommonParams applies shared parameters (resolution, aspect_ratio, etc.)
func applyCommonParams(arkReq *ArkCreateTaskRequest, body map[string]any) {
	if arkReq.Resolution != "" {
		body["resolution"] = arkReq.Resolution
	}
	if arkReq.Ratio != "" {
		body["aspect_ratio"] = arkReq.Ratio
	}
	if arkReq.Duration != nil {
		body["duration"] = *arkReq.Duration
	}
	if arkReq.Seed != nil {
		body["seed"] = *arkReq.Seed
	}
	if arkReq.GenerateAudio != nil {
		body["generate_audio"] = *arkReq.GenerateAudio
	}
	if arkReq.CameraFixed != nil {
		body["camera_fixed"] = *arkReq.CameraFixed
	}
	if arkReq.Watermark != nil {
		body["watermark"] = *arkReq.Watermark
	}
}

// ─────────────────────────────────────────────────────────────────
// MuleRun → Ark response conversion
// ─────────────────────────────────────────────────────────────────

// MuleRunToArkCreateResponse converts a MuleRun task creation response
// to the Ark standard format.
func MuleRunToArkCreateResponse(muleRunBody []byte, vendorQueryPath, modelName string) (*ArkCreateTaskResponse, error) {
	var mrResp MuleRunVendorTaskResponse
	if err := json.Unmarshal(muleRunBody, &mrResp); err != nil {
		return nil, fmt.Errorf("parse mulerun response: %w", err)
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

	// Encode composite task ID
	compositeID := EncodeCompositeTaskID(vendorQueryPath, upstreamTaskID, modelName)

	return &ArkCreateTaskResponse{
		ID: compositeID,
	}, nil
}

// MuleRunToArkQueryResponse converts a MuleRun task query response
// to the Ark standard format.
func MuleRunToArkQueryResponse(muleRunBody []byte, modelName string) (*ArkQueryTaskResponse, error) {
	var mrResp MuleRunVendorTaskResponse
	if err := json.Unmarshal(muleRunBody, &mrResp); err != nil {
		return nil, fmt.Errorf("parse mulerun response: %w", err)
	}

	resp := &ArkQueryTaskResponse{
		Model: modelName,
	}

	// Extract task info
	if mrResp.TaskInfo != nil {
		resp.ID = mrResp.TaskInfo.ID
		resp.Status = mapMuleRunStatus(mrResp.TaskInfo.Status)
		resp.CreatedAt = parseTimestamp(mrResp.TaskInfo.CreatedAt)
		resp.UpdatedAt = parseTimestamp(mrResp.TaskInfo.UpdatedAt)
	} else {
		resp.ID = mrResp.ID
		resp.Status = mapMuleRunStatus(mrResp.Status)
	}

	// Extract video URL
	if len(mrResp.Videos) > 0 {
		resp.Content = &ArkOutputContent{
			VideoURL: mrResp.Videos[0],
		}
	}

	return resp, nil
}

// mapMuleRunStatus converts MuleRun status to Ark status.
func mapMuleRunStatus(muleRunStatus string) string {
	switch muleRunStatus {
	case MuleRunStatusPending:
		return ArkStatusQueued
	case MuleRunStatusProcessing:
		return ArkStatusRunning
	case MuleRunStatusCompleted:
		return ArkStatusSucceeded
	case MuleRunStatusFailed:
		return ArkStatusFailed
	default:
		return muleRunStatus
	}
}

// parseTimestamp attempts to parse a timestamp string to Unix seconds.
func parseTimestamp(ts string) int64 {
	if ts == "" {
		return 0
	}
	// Try parsing as RFC3339
	if t, err := parseTimeRFC3339(ts); err == nil {
		return t
	}
	// Try parsing as Unix timestamp string
	var unix int64
	if _, err := fmt.Sscanf(ts, "%d", &unix); err == nil {
		return unix
	}
	return 0
}

func parseTimeRFC3339(s string) (int64, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.Unix(), nil
		}
	}
	return 0, fmt.Errorf("unable to parse timestamp: %s", s)
}
