package service

import "encoding/json"

// ─────────────────────────────────────────────────────────────────
// Seedance / Volcengine Ark standard format types
// These types define the public-facing multimedia API contract.
// ─────────────────────────────────────────────────────────────────

// ArkCreateTaskRequest is the standard request format for creating
// a multimedia generation task (video, image, audio).
// Based on Volcengine Ark "创建视频生成任务 API".
type ArkCreateTaskRequest struct {
	Model    string          `json:"model"`
	Content  []ArkContent   `json:"content"`
	Resolution string        `json:"resolution,omitempty"`
	Ratio      string        `json:"ratio,omitempty"`
	Duration   *int          `json:"duration,omitempty"`
	Frames     *int          `json:"frames,omitempty"`
	Seed       *int          `json:"seed,omitempty"`
	CameraFixed *bool        `json:"camera_fixed,omitempty"`
	Watermark  *bool         `json:"watermark,omitempty"`
	GenerateAudio *bool      `json:"generate_audio,omitempty"`
	Draft        *bool       `json:"draft,omitempty"`
	ReturnLastFrame *bool    `json:"return_last_frame,omitempty"`
	ServiceTier     string   `json:"service_tier,omitempty"`
	ExecutionExpiresAfter *int `json:"execution_expires_after,omitempty"`
	CallbackURL     string   `json:"callback_url,omitempty"`
	SafetyIdentifier string  `json:"safety_identifier,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
	Tools            json.RawMessage `json:"tools,omitempty"`
}

// ArkContent represents a single content input item.
type ArkContent struct {
	Type      string           `json:"type"`
	Text      string           `json:"text,omitempty"`
	ImageURL  *ArkMediaURL     `json:"image_url,omitempty"`
	VideoURL  *ArkMediaURL     `json:"video_url,omitempty"`
	AudioURL  *ArkMediaURL     `json:"audio_url,omitempty"`
	DraftTask *ArkDraftTask    `json:"draft_task,omitempty"`
	Role      string           `json:"role,omitempty"`
}

// ArkMediaURL wraps a URL for media inputs (image/video/audio).
type ArkMediaURL struct {
	URL string `json:"url"`
}

// ArkDraftTask references a draft task for upgrading to final video.
type ArkDraftTask struct {
	ID string `json:"id"`
}

// ArkCreateTaskResponse is returned after creating a task.
type ArkCreateTaskResponse struct {
	ID string `json:"id"`
}

// ArkQueryTaskResponse is the standard response for querying a task.
type ArkQueryTaskResponse struct {
	ID            string              `json:"id"`
	Model         string              `json:"model,omitempty"`
	Status        string              `json:"status"`
	Error         *ArkError           `json:"error"`
	CreatedAt     int64               `json:"created_at,omitempty"`
	UpdatedAt     int64               `json:"updated_at,omitempty"`
	Content       *ArkOutputContent   `json:"content,omitempty"`
	Seed          *int                `json:"seed,omitempty"`
	Resolution    string              `json:"resolution,omitempty"`
	Ratio         string              `json:"ratio,omitempty"`
	Duration      *int                `json:"duration,omitempty"`
	Frames        *int                `json:"frames,omitempty"`
	FramesPerSecond *int              `json:"framespersecond,omitempty"`
	GenerateAudio *bool               `json:"generate_audio,omitempty"`
	Usage         *ArkUsage           `json:"usage,omitempty"`
	ServiceTier   string              `json:"service_tier,omitempty"`
}

// ArkError contains error details for failed tasks.
type ArkError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ArkOutputContent contains the generated output.
type ArkOutputContent struct {
	VideoURL     string `json:"video_url,omitempty"`
	LastFrameURL string `json:"last_frame_url,omitempty"`
}

// ArkUsage contains token usage information.
type ArkUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ─── Ark status constants ───

const (
	ArkStatusQueued    = "queued"
	ArkStatusRunning   = "running"
	ArkStatusSucceeded = "succeeded"
	ArkStatusFailed    = "failed"
	ArkStatusCancelled = "cancelled"
	ArkStatusExpired   = "expired"
)

// ─── Ark content type constants ───

const (
	ArkContentTypeText      = "text"
	ArkContentTypeImageURL  = "image_url"
	ArkContentTypeVideoURL  = "video_url"
	ArkContentTypeAudioURL  = "audio_url"
	ArkContentTypeDraftTask = "draft_task"
)

// ─── Ark role constants ───

const (
	ArkRoleFirstFrame     = "first_frame"
	ArkRoleLastFrame      = "last_frame"
	ArkRoleReferenceImage = "reference_image"
	ArkRoleReferenceVideo = "reference_video"
	ArkRoleReferenceAudio = "reference_audio"
)

// ─────────────────────────────────────────────────────────────────
// OpenAI Image API format types
// POST /v1/openai/images/generations
// POST /v1/openai/images/edits
// ─────────────────────────────────────────────────────────────────

// OpenAIImageRequest is the OpenAI Image API request format.
// See: https://developers.openai.com/api/docs/api-reference/images/create
type OpenAIImageRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              *int   `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // "b64_json" or "url"
	OutputFormat   string `json:"output_format,omitempty"`   // "png", "jpeg", "webp"
	OutputCompression *int `json:"output_compression,omitempty"`
	Background     string `json:"background,omitempty"` // "opaque", "transparent", "auto"
	Moderation     string `json:"moderation,omitempty"` // "auto", "low"
	User           string `json:"user,omitempty"`
}

// OpenAIImageResponse is the OpenAI Image API response format.
type OpenAIImageResponse struct {
	Created int64              `json:"created"`
	Data    []OpenAIImageData  `json:"data"`
}

// OpenAIImageData represents one generated image.
type OpenAIImageData struct {
	B64JSON       string `json:"b64_json,omitempty"`
	URL           string `json:"url,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// ─────────────────────────────────────────────────────────────────
// Google Video API format types (Veo / Gemini format)
// POST /v1/google/models/{model}:predictLongRunning
// GET  /v1/google/{operation_name}
// ─────────────────────────────────────────────────────────────────

// GoogleVideoRequest is the Google Veo video generation request format.
type GoogleVideoRequest struct {
	Model      string                    `json:"model"`
	Instances  []GoogleVideoInstance     `json:"instances"`
	Parameters *GoogleVideoParameters    `json:"parameters,omitempty"`
}

// GoogleVideoInstance represents one video generation input.
type GoogleVideoInstance struct {
	Prompt          string                      `json:"prompt"`
	Image           *GoogleInlineDataWrapper `json:"image,omitempty"`
	LastFrame       *GoogleInlineDataWrapper `json:"lastFrame,omitempty"`
	ReferenceImages []GoogleReferenceImage      `json:"referenceImages,omitempty"`
}

// GoogleVideoParameters contains video generation parameters.
type GoogleVideoParameters struct {
	AspectRatio      string `json:"aspectRatio,omitempty"`
	Resolution       string `json:"resolution,omitempty"`
	DurationSeconds  string `json:"durationSeconds,omitempty"`
	PersonGeneration string `json:"personGeneration,omitempty"`
	Seed             *int   `json:"seed,omitempty"`
}

// GoogleReferenceImage is a reference image for video generation.
type GoogleReferenceImage struct {
	Image         *GoogleInlineDataWrapper `json:"image,omitempty"`
	ReferenceType string                   `json:"referenceType,omitempty"`
}

// GoogleInlineDataWrapper wraps an inline data blob.
type GoogleInlineDataWrapper struct {
	InlineData *GoogleBlob `json:"inlineData,omitempty"`
}

// GoogleBlob contains base64-encoded data with MIME type.
type GoogleBlob struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// GoogleVideoCreateResponse is returned after creating a video generation operation.
type GoogleVideoCreateResponse struct {
	Name string `json:"name"`
}

// GoogleVideoOperationResponse is the polling response for a video operation.
type GoogleVideoOperationResponse struct {
	Name     string                      `json:"name"`
	Done     bool                        `json:"done"`
	Response *GoogleVideoOperationResult `json:"response,omitempty"`
	Error    *GoogleOperationError       `json:"error,omitempty"`
}

// GoogleVideoOperationResult wraps the final video response.
type GoogleVideoOperationResult struct {
	GenerateVideoResponse *GoogleGenerateVideoResponse `json:"generateVideoResponse,omitempty"`
}

// GoogleGenerateVideoResponse contains the generated video samples.
type GoogleGenerateVideoResponse struct {
	GeneratedSamples []GoogleVideoSample `json:"generatedSamples,omitempty"`
}

// GoogleVideoSample represents one generated video.
type GoogleVideoSample struct {
	Video *GoogleVideoFile `json:"video,omitempty"`
}

// GoogleVideoFile contains the video URI.
type GoogleVideoFile struct {
	URI string `json:"uri,omitempty"`
}

// GoogleOperationError contains error info for failed operations.
type GoogleOperationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ─────────────────────────────────────────────────────────────────
// Google Image API format types (Nano Banana / Gemini format)
// POST /v1/google/models/{model}:generateContent
// ─────────────────────────────────────────────────────────────────

// GoogleImageRequest is the Gemini generateContent request format for images.
type GoogleImageRequest struct {
	Model    string                `json:"model"`
	Contents []GoogleContent       `json:"contents"`
	Config   *GoogleGenerationConfig `json:"generationConfig,omitempty"`
}

// GoogleContent represents one content item in a Gemini request.
type GoogleContent struct {
	Parts []GooglePart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

// GooglePart is a single part in content (text or inline data).
type GooglePart struct {
	Text      string     `json:"text,omitempty"`
	InlineData *GoogleBlob `json:"inline_data,omitempty"`
}

// GoogleGenerationConfig contains generation configuration.
type GoogleGenerationConfig struct {
	ResponseModalities []string `json:"responseModalities,omitempty"`
}

// GoogleImageResponse is the Gemini generateContent response format.
type GoogleImageResponse struct {
	Candidates []GoogleCandidate `json:"candidates"`
}

// GoogleCandidate represents one generation candidate.
type GoogleCandidate struct {
	Content GoogleContentOutput `json:"content"`
}

// GoogleContentOutput contains the generated parts.
type GoogleContentOutput struct {
	Parts []GooglePart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}

// ─── MuleRun vendor task response (what MuleRun returns) ───

// MuleRunVendorTaskResponse represents the raw response from MuleRun
// vendor endpoints.
type MuleRunVendorTaskResponse struct {
	TaskInfo *MuleRunTaskInfo `json:"task_info,omitempty"`
	Videos   []string         `json:"videos,omitempty"`
	// Some vendors may return these directly
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

// MuleRunTaskInfo contains task metadata from MuleRun.
type MuleRunTaskInfo struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// MuleRun vendor status constants
const (
	MuleRunStatusPending    = "pending"
	MuleRunStatusProcessing = "processing"
	MuleRunStatusCompleted  = "completed"
	MuleRunStatusFailed     = "failed"
)
