// =====================
// 模型列表（硬编码，与 new-api 一致）
// =====================

// OpenAI
const openaiModels = [
  // GPT-5.2 系列
  'gpt-5.2', 'gpt-5.2-2025-12-11', 'gpt-5.2-chat-latest',
  'gpt-5.2-pro', 'gpt-5.2-pro-2025-12-11',
  // GPT-5.5 系列
  'gpt-5.5',
  // GPT-5.4 系列
  'gpt-5.4', 'gpt-5.4-mini', 'gpt-5.4-2026-03-05',
  // GPT-5.3 / Codex 系列
  'gpt-5.3-codex', 'gpt-5.3-codex-spark', 'codex-auto-review',
  'gpt-4o-audio-preview', 'gpt-4o-realtime-preview',
  // GPT Image 系列
  'gpt-image-1', 'gpt-image-1.5', 'gpt-image-2'
]

// Anthropic Claude
export const claudeModels = [
  'claude-3-5-sonnet-20241022', 'claude-3-5-sonnet-20240620',
  'claude-3-5-haiku-20241022',
  'claude-3-7-sonnet-20250219',
  'claude-sonnet-4-20250514', 'claude-opus-4-20250514',
  'claude-opus-4-1-20250805',
  'claude-sonnet-4-5-20250929', 'claude-haiku-4-5-20251001',
  'claude-opus-4-5-20251101',
  'claude-opus-4-6',
  'claude-opus-4-7',
  'claude-opus-4-8',
  'claude-sonnet-4-6',
  'claude-fable-5'
]

// Google Gemini
const geminiModels = [
  // Keep in sync with backend curated Gemini lists.
  // This list is intentionally conservative (models commonly available across OAuth/API key).
  'gemini-3.1-flash-image',
  'gemini-2.5-flash-image',
  'gemini-2.0-flash',
  'gemini-2.5-flash',
  'gemini-2.5-pro',
  'gemini-3.5-flash',
  'gemini-3-flash-preview',
  'gemini-3-pro-preview'
]

// Antigravity 官方支持的模型（精确匹配）
// 基于官方 API 返回的模型列表，只支持 Claude 4.5+ 和 Gemini 2.5+
const antigravityModels = [
  // Claude 4.5+ 系列
  'claude-fable-5',
  'claude-opus-4-6',
  'claude-opus-4-6-thinking',
  'claude-opus-4-7',
  'claude-opus-4-8',
  'claude-opus-4-5-thinking',
  'claude-sonnet-4-6',
  'claude-sonnet-4-5',
  'claude-sonnet-4-5-thinking',
  // Gemini 2.5 系列
  'gemini-3.1-flash-image',
  'gemini-2.5-flash-image',
  'gemini-2.5-flash',
  'gemini-2.5-flash-lite',
  'gemini-2.5-flash-thinking',
  'gemini-2.5-pro',
  // Gemini 3 系列
  'gemini-3-flash',
  'gemini-3-pro-high',
  'gemini-3-pro-low',
  // Gemini 3.1 系列
  'gemini-3.1-pro-high',
  'gemini-3.1-pro-low',
  'gemini-3-pro-image',
  // 其他
  'gpt-oss-120b-medium',
  'tab_flash_lite_preview'
]

// MuleRun 支持的模型（文本 + Vendor 多媒体端点）
const mulerunModels = [
  // ── 文本模型（Chat Completions /v1/chat/completions）──
  // Gemini 系列
  'gemini-3-pro-preview', 'gemini-2.5-pro', 'gemini-2.5-flash',
  'gemini-2.5-flash-lite', 'gemini-2.5-flash-image-preview', 'gemini-2.5-flash-image',
  // Qwen 系列
  'qwen3-max', 'qwen-plus', 'qwen-flash', 'qwen3-vl-plus', 'qwen3-vl-flash',
  // 智谱
  'glm-4.6',
  // ── 文本模型（Anthropic /v1/messages — 需要高级套餐）──
  'claude-sonnet-4-6', 'claude-opus-4-6', 'claude-haiku-4-5',
  // ── 文本模型（Responses API /v1/responses）──
  'gpt-5.4', 'gpt-5.2', 'gpt-5', 'gpt-5-mini', 'gpt-5-nano', 'gpt-5-codex',
  // ── Vendor: Alibaba 系列（/vendors/alibaba/v1/...）──
  'alibaba/happy-horse-1-0-t2v/generation', 'alibaba/happy-horse-1-0-i2v/generation',
  'alibaba/wan2.6-t2v/generation', 'alibaba/wan2.6-t2i/generation',
  'alibaba/wan2.6-image/generation', 'alibaba/wan2.6-i2v/generation',
  'alibaba/wan2.5-t2v-preview/generation', 'alibaba/wan2.5-t2i-preview/generation',
  'alibaba/wan2.5-i2v-preview/generation', 'alibaba/wan2.5-i2i-preview/generation',
  'alibaba/wan2.2-t2v-plus/generation', 'alibaba/wan2.2-i2v-plus/generation', 'alibaba/wan2.2-i2v-flash/generation',
  'alibaba/wan2.1-vace-plus/generation', 'alibaba/wan2.1-kf2v-plus/generation',
  // ── Vendor: ByteDance Seedance 系列 ──
  'bytedance/seedance-2.0/text-to-video',
  'bytedance/seedance-2.0/image-to-video',
  'bytedance/seedance-2.0/reference-to-video',
  'bytedance/seedance-2.0-fast/text-to-video',
  'bytedance/seedance-2.0-fast/image-to-video',
  'bytedance/seedance-2.0-fast/reference-to-video',
  // ── Vendor: Google 系列 ──
  'google/veo3/generation',
  'google/nano-banana-pro/generation', 'google/nano-banana-pro/edit',
  'google/nano-banana-2/generation', 'google/nano-banana-2/edit',
  'google/nano-banana/generation', 'google/nano-banana/edit',
  // ── Vendor: KlingAI 系列 ──
  'klingai/kling-v3-t2v/generation', 'klingai/kling-v3-i2v/generation',
  'klingai/kling-v3-omni-t2v/generation', 'klingai/kling-v3-omni-i2v/generation',
  'klingai/kling-v3-omni-ref2v/generation', 'klingai/kling-v3-omni-v2v/generation',
  'klingai/kling-v3-omni-v2v-edit/generation',
  // ── Vendor: Midjourney ──
  'midjourney/diffusion/generation', 'midjourney/video/generation',
  // ── Vendor: MiniMax ──
  'minimax/music-2.5/generation', 'minimax/music-2.0/generation',
  'minimax/speech-2.8-hd/generation', 'minimax/speech-2.8-turbo/generation',
  // ── Vendor: OpenAI ──
  'openai/sora/generation', 'openai/sora-2/generation',
  'openai/gpt-image-2/generation', 'openai/gpt-image-2/edit'
]

// 智谱 GLM
const zhipuModels = [
  'glm-4', 'glm-4v', 'glm-4-plus', 'glm-4-0520',
  'glm-4-air', 'glm-4-airx', 'glm-4-long', 'glm-4-flash',
  'glm-4v-plus', 'glm-4.5', 'glm-4.6',
  'glm-3-turbo', 'glm-4-alltools',
  'chatglm_turbo', 'chatglm_pro', 'chatglm_std', 'chatglm_lite',
  'cogview-3', 'cogvideo'
]

// 阿里 通义千问
const qwenModels = [
  'qwen-turbo', 'qwen-plus', 'qwen-max', 'qwen-max-longcontext', 'qwen-long',
  'qwen2-72b-instruct', 'qwen2-57b-a14b-instruct', 'qwen2-7b-instruct',
  'qwen2.5-72b-instruct', 'qwen2.5-32b-instruct', 'qwen2.5-14b-instruct',
  'qwen2.5-7b-instruct', 'qwen2.5-3b-instruct', 'qwen2.5-1.5b-instruct',
  'qwen2.5-coder-32b-instruct', 'qwen2.5-coder-14b-instruct', 'qwen2.5-coder-7b-instruct',
  'qwen3-235b-a22b',
  'qwq-32b', 'qwq-32b-preview'
]

// DeepSeek
const deepseekModels = [
  'deepseek-chat', 'deepseek-coder', 'deepseek-reasoner',
  'deepseek-v3', 'deepseek-v3-0324',
  'deepseek-r1', 'deepseek-r1-0528',
  'deepseek-r1-distill-qwen-32b', 'deepseek-r1-distill-qwen-14b', 'deepseek-r1-distill-qwen-7b',
  'deepseek-r1-distill-llama-70b', 'deepseek-r1-distill-llama-8b'
]

// Mistral
const mistralModels = [
  'mistral-small-latest', 'mistral-medium-latest', 'mistral-large-latest',
  'open-mistral-7b', 'open-mixtral-8x7b', 'open-mixtral-8x22b',
  'codestral-latest', 'codestral-mamba',
  'pixtral-12b-2409', 'pixtral-large-latest'
]

// Meta Llama
const metaModels = [
  'llama-3.3-70b-instruct',
  'llama-3.2-90b-vision-instruct', 'llama-3.2-11b-vision-instruct',
  'llama-3.2-3b-instruct', 'llama-3.2-1b-instruct',
  'llama-3.1-405b-instruct', 'llama-3.1-70b-instruct', 'llama-3.1-8b-instruct',
  'llama-3-70b-instruct', 'llama-3-8b-instruct',
  'codellama-70b-instruct', 'codellama-34b-instruct', 'codellama-13b-instruct'
]

// xAI Grok
const xaiModels = [
  'grok-4', 'grok-4-0709',
  'grok-3-beta', 'grok-3-mini-beta', 'grok-3-fast-beta',
  'grok-2', 'grok-2-vision', 'grok-2-image',
  'grok-beta', 'grok-vision-beta'
]

// Cohere
const cohereModels = [
  'command-a-03-2025',
  'command-r', 'command-r-plus',
  'command-r-08-2024', 'command-r-plus-08-2024',
  'c4ai-aya-23-35b', 'c4ai-aya-23-8b',
  'command', 'command-light'
]

// Yi (01.AI)
const yiModels = [
  'yi-large', 'yi-large-turbo', 'yi-large-rag',
  'yi-medium', 'yi-medium-200k',
  'yi-spark', 'yi-vision',
  'yi-1.5-34b-chat', 'yi-1.5-9b-chat', 'yi-1.5-6b-chat'
]

// Moonshot/Kimi
const moonshotModels = [
  'moonshot-v1-8k', 'moonshot-v1-32k', 'moonshot-v1-128k',
  'kimi-latest'
]

// 字节跳动 豆包
const doubaoModels = [
  'doubao-pro-256k', 'doubao-pro-128k', 'doubao-pro-32k', 'doubao-pro-4k',
  'doubao-lite-128k', 'doubao-lite-32k', 'doubao-lite-4k',
  'doubao-vision-pro-32k', 'doubao-vision-lite-32k',
  'doubao-1.5-pro-256k', 'doubao-1.5-pro-32k', 'doubao-1.5-lite-32k',
  'doubao-1.5-pro-vision-32k', 'doubao-1.5-thinking-pro'
]

// MiniMax
const minimaxModels = [
  'abab6.5-chat', 'abab6.5s-chat', 'abab6.5s-chat-pro',
  'abab6-chat',
  'abab5.5-chat', 'abab5.5s-chat'
]

// 百度 文心
const baiduModels = [
  'ernie-4.0-8k-latest', 'ernie-4.0-8k', 'ernie-4.0-turbo-8k',
  'ernie-3.5-8k', 'ernie-3.5-128k',
  'ernie-speed-8k', 'ernie-speed-128k', 'ernie-speed-pro-128k',
  'ernie-lite-8k', 'ernie-lite-pro-128k',
  'ernie-tiny-8k'
]

// 讯飞 星火
const sparkModels = [
  'spark-desk', 'spark-desk-v1.1', 'spark-desk-v2.1',
  'spark-desk-v3.1', 'spark-desk-v3.5', 'spark-desk-v4.0',
  'spark-lite', 'spark-pro', 'spark-max', 'spark-ultra'
]

// 腾讯 混元
const hunyuanModels = [
  'hunyuan-lite', 'hunyuan-standard', 'hunyuan-standard-256k',
  'hunyuan-pro', 'hunyuan-turbo', 'hunyuan-large',
  'hunyuan-vision', 'hunyuan-code'
]

// Perplexity
const perplexityModels = [
  'sonar', 'sonar-pro', 'sonar-reasoning',
  'llama-3-sonar-small-32k-online', 'llama-3-sonar-large-32k-online',
  'llama-3-sonar-small-32k-chat', 'llama-3-sonar-large-32k-chat'
]

// 所有模型（去重）
const allModelsList: string[] = [
  ...openaiModels,
  ...claudeModels,
  ...geminiModels,
  ...antigravityModels,
  ...mulerunModels,
  ...zhipuModels,
  ...qwenModels,
  ...deepseekModels,
  ...mistralModels,
  ...metaModels,
  ...xaiModels,
  ...cohereModels,
  ...yiModels,
  ...moonshotModels,
  ...doubaoModels,
  ...minimaxModels,
  ...baiduModels,
  ...sparkModels,
  ...hunyuanModels,
  ...perplexityModels
]

// 转换为下拉选项格式
export const allModels = allModelsList.map(m => ({ value: m, label: m }))

// =====================
// 预设映射
// =====================

const anthropicPresetMappings = [
  { label: 'Fable 5', from: 'claude-fable-5', to: 'claude-fable-5', color: 'bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-900/30 dark:text-rose-400' },
  { label: 'Sonnet 4', from: 'claude-sonnet-4-20250514', to: 'claude-sonnet-4-20250514', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: 'Sonnet 4.5', from: 'claude-sonnet-4-5-20250929', to: 'claude-sonnet-4-5-20250929', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  { label: 'Sonnet 4.6', from: 'claude-sonnet-4-6', to: 'claude-sonnet-4-6', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  { label: 'Opus 4.5', from: 'claude-opus-4-5-20251101', to: 'claude-opus-4-5-20251101', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Opus 4.6', from: 'claude-opus-4-6', to: 'claude-opus-4-6', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Opus 4.7', from: 'claude-opus-4-7', to: 'claude-opus-4-7', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Opus 4.8', from: 'claude-opus-4-8', to: 'claude-opus-4-8', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Haiku 3.5', from: 'claude-3-5-haiku-20241022', to: 'claude-3-5-haiku-20241022', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  { label: 'Haiku 4.5', from: 'claude-haiku-4-5-20251001', to: 'claude-haiku-4-5-20251001', color: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400' },
  { label: 'Opus->Sonnet', from: 'claude-opus-4-6', to: 'claude-sonnet-4-5-20250929', color: 'bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400' }
]

const openaiPresetMappings = [
  { label: 'GPT-4o', from: 'gpt-4o', to: 'gpt-4o', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  { label: 'GPT-4o Mini', from: 'gpt-4o-mini', to: 'gpt-4o-mini', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: 'GPT-4.1', from: 'gpt-4.1', to: 'gpt-4.1', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  { label: 'o1', from: 'o1', to: 'o1', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'o3', from: 'o3', to: 'o3', color: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400' },
  { label: 'GPT-5.3 Codex Spark', from: 'gpt-5.3-codex-spark', to: 'gpt-5.3-codex-spark', color: 'bg-teal-100 text-teal-700 hover:bg-teal-200 dark:bg-teal-900/30 dark:text-teal-400' },
  { label: 'GPT-5.2', from: 'gpt-5.2', to: 'gpt-5.2', color: 'bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900/30 dark:text-red-400' },
  { label: 'GPT-5.5', from: 'gpt-5.5', to: 'gpt-5.5', color: 'bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400' },
  { label: 'GPT-5.4', from: 'gpt-5.4', to: 'gpt-5.4', color: 'bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-900/30 dark:text-rose-400' },
  { label: 'Haiku→5.4', from: 'claude-haiku-4-5-20251001', to: 'gpt-5.4', color: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400' },
  { label: 'Opus→5.4', from: 'claude-opus-4-6', to: 'gpt-5.4', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Sonnet→5.4', from: 'claude-sonnet-4-6', to: 'gpt-5.4', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' }
]

const geminiPresetMappings = [
  { label: 'Flash 2.0', from: 'gemini-2.0-flash', to: 'gemini-2.0-flash', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: '2.5 Flash', from: 'gemini-2.5-flash', to: 'gemini-2.5-flash', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  { label: '2.5 Image', from: 'gemini-2.5-flash-image', to: 'gemini-2.5-flash-image', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: '2.5 Pro', from: 'gemini-2.5-pro', to: 'gemini-2.5-pro', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: '3.5 Flash', from: 'gemini-3.5-flash', to: 'gemini-3.5-flash', color: 'bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-900/30 dark:text-violet-400' },
  { label: '3.1 Image', from: 'gemini-3.1-flash-image', to: 'gemini-3.1-flash-image', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' }
]

// Antigravity 预设映射（支持通配符）
const antigravityPresetMappings = [
  // Claude 通配符映射
  { label: 'Claude→Sonnet', from: 'claude-*', to: 'claude-sonnet-4-5', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: 'Fable 5', from: 'claude-fable-5', to: 'claude-fable-5', color: 'bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-900/30 dark:text-rose-400' },
  { label: 'Sonnet→Sonnet', from: 'claude-sonnet-*', to: 'claude-sonnet-4-5', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  { label: 'Opus→Opus', from: 'claude-opus-*', to: 'claude-opus-4-6-thinking', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Haiku→Sonnet', from: 'claude-haiku-*', to: 'claude-sonnet-4-5', color: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400' },
  { label: 'Sonnet4→4.6', from: 'claude-sonnet-4-20250514', to: 'claude-sonnet-4-6', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: 'Sonnet4.5→4.6', from: 'claude-sonnet-4-5-20250929', to: 'claude-sonnet-4-6', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Sonnet3.5→4.6', from: 'claude-3-5-sonnet-20241022', to: 'claude-sonnet-4-6', color: 'bg-teal-100 text-teal-700 hover:bg-teal-200 dark:bg-teal-900/30 dark:text-teal-400' },
  { label: 'Opus4.5→4.6', from: 'claude-opus-4-5-20251101', to: 'claude-opus-4-6-thinking', color: 'bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-900/30 dark:text-violet-400' },
  // Gemini 3→3.1 映射
  { label: '3-Pro-Preview→3.1-Pro-High', from: 'gemini-3-pro-preview', to: 'gemini-3.1-pro-high', color: 'bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400' },
  { label: '3-Pro-High→3.1-Pro-High', from: 'gemini-3-pro-high', to: 'gemini-3.1-pro-high', color: 'bg-orange-100 text-orange-700 hover:bg-orange-200 dark:bg-orange-900/30 dark:text-orange-400' },
  { label: '3-Pro-Low→3.1-Pro-Low', from: 'gemini-3-pro-low', to: 'gemini-3.1-pro-low', color: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400' },
  { label: '3.1-Pro-High透传', from: 'gemini-3.1-pro-high', to: 'gemini-3.1-pro-high', color: 'bg-orange-100 text-orange-700 hover:bg-orange-200 dark:bg-orange-900/30 dark:text-orange-400' },
  { label: '3.1-Pro-Low透传', from: 'gemini-3.1-pro-low', to: 'gemini-3.1-pro-low', color: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400' },
  // Gemini 通配符映射
  { label: 'Gemini 3→Flash', from: 'gemini-3*', to: 'gemini-3-flash', color: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400' },
  { label: 'Gemini 2.5→Flash', from: 'gemini-2.5*', to: 'gemini-2.5-flash', color: 'bg-orange-100 text-orange-700 hover:bg-orange-200 dark:bg-orange-900/30 dark:text-orange-400' },
  { label: '2.5-Flash-Image透传', from: 'gemini-2.5-flash-image', to: 'gemini-2.5-flash-image', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: '3.1-Flash-Image透传', from: 'gemini-3.1-flash-image', to: 'gemini-3.1-flash-image', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: '3-Pro-Image→3.1', from: 'gemini-3-pro-image', to: 'gemini-3.1-flash-image', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: '3-Flash透传', from: 'gemini-3-flash', to: 'gemini-3-flash', color: 'bg-lime-100 text-lime-700 hover:bg-lime-200 dark:bg-lime-900/30 dark:text-lime-400' },
  { label: '2.5-Flash-Lite透传', from: 'gemini-2.5-flash-lite', to: 'gemini-2.5-flash-lite', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  // 精确映射
  { label: 'Sonnet 4.6', from: 'claude-sonnet-4-6', to: 'claude-sonnet-4-6', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Sonnet 4.5', from: 'claude-sonnet-4-5', to: 'claude-sonnet-4-5', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Opus 4.6', from: 'claude-opus-4-6', to: 'claude-opus-4-6-thinking', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Opus 4.6-thinking', from: 'claude-opus-4-6-thinking', to: 'claude-opus-4-6-thinking', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Opus 4.7', from: 'claude-opus-4-7', to: 'claude-opus-4-7', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Opus 4.8', from: 'claude-opus-4-8', to: 'claude-opus-4-8', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' }
]

// Bedrock 预设映射（与后端 DefaultBedrockModelMapping 保持一致）
const bedrockPresetMappings = [
  { label: 'Fable 5', from: 'claude-fable-5', to: 'anthropic.claude-fable-5', color: 'bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-900/30 dark:text-rose-400' },
  { label: 'Opus 4.6', from: 'claude-opus-4-6', to: 'us.anthropic.claude-opus-4-6-v1', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Opus 4.7', from: 'claude-opus-4-7', to: 'us.anthropic.claude-opus-4-7-v1', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Opus 4.8', from: 'claude-opus-4-8', to: 'us.anthropic.claude-opus-4-8-v1', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Sonnet 4.6', from: 'claude-sonnet-4-6', to: 'us.anthropic.claude-sonnet-4-6', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Opus 4.5', from: 'claude-opus-4-5-thinking', to: 'us.anthropic.claude-opus-4-5-20251101-v1:0', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Sonnet 4.5', from: 'claude-sonnet-4-5', to: 'us.anthropic.claude-sonnet-4-5-20250929-v1:0', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Haiku 4.5', from: 'claude-haiku-4-5', to: 'us.anthropic.claude-haiku-4-5-20251001-v1:0', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
]

// MuleRun 预设映射（文本模型 + Vendor 多媒体端点）
const mulerunPresetMappings = [
  // ── Chat Completions 文本模型 ──
  { label: 'Gemini 2.5 Pro', from: 'gemini-2.5-pro', to: 'gemini-2.5-pro', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: 'Gemini 2.5 Flash', from: 'gemini-2.5-flash', to: 'gemini-2.5-flash', color: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-400' },
  { label: 'Qwen3 Max', from: 'qwen3-max', to: 'qwen3-max', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'Qwen Flash', from: 'qwen-flash', to: 'qwen-flash', color: 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-400' },
  { label: 'GLM 4.6', from: 'glm-4.6', to: 'glm-4.6', color: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:text-emerald-400' },
  // ── Claude 系列（高级套餐）──
  { label: 'Claude Sonnet 4.6', from: 'claude-sonnet-4-6', to: 'claude-sonnet-4-6', color: 'bg-cyan-100 text-cyan-700 hover:bg-cyan-200 dark:bg-cyan-900/30 dark:text-cyan-400' },
  { label: 'Claude Opus 4.6', from: 'claude-opus-4-6', to: 'claude-opus-4-6', color: 'bg-pink-100 text-pink-700 hover:bg-pink-200 dark:bg-pink-900/30 dark:text-pink-400' },
  { label: 'Claude Haiku 4.5', from: 'claude-haiku-4-5', to: 'claude-haiku-4-5', color: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 dark:bg-indigo-900/30 dark:text-indigo-400' },
  // ── Responses API（GPT 系列）──
  { label: 'GPT-5.4', from: 'gpt-5.4', to: 'gpt-5.4', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  { label: 'GPT-5 Mini', from: 'gpt-5-mini', to: 'gpt-5-mini', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  // ── 视频生成 ──
  { label: 'Wan2.6 T2V', from: 'alibaba/wan2.6-t2v/generation', to: 'alibaba/wan2.6-t2v/generation', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: 'Wan2.6 T2I', from: 'alibaba/wan2.6-t2i/generation', to: 'alibaba/wan2.6-t2i/generation', color: 'bg-teal-100 text-teal-700 hover:bg-teal-200 dark:bg-teal-900/30 dark:text-teal-400' },
  { label: 'Happy Horse T2V', from: 'alibaba/happy-horse-1-0-t2v/generation', to: 'alibaba/happy-horse-1-0-t2v/generation', color: 'bg-sky-100 text-sky-700 hover:bg-sky-200 dark:bg-sky-900/30 dark:text-sky-400' },
  { label: 'Veo 3', from: 'google/veo3/generation', to: 'google/veo3/generation', color: 'bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400' },
  { label: 'Seedance 2.0 T2V', from: 'bytedance/seedance-2.0/text-to-video', to: 'bytedance/seedance-2.0/text-to-video', color: 'bg-rose-100 text-rose-700 hover:bg-rose-200 dark:bg-rose-900/30 dark:text-rose-400' },
  { label: 'Kling V3 T2V', from: 'klingai/kling-v3-t2v/generation', to: 'klingai/kling-v3-t2v/generation', color: 'bg-orange-100 text-orange-700 hover:bg-orange-200 dark:bg-orange-900/30 dark:text-orange-400' },
  { label: 'Kling V3 Omni T2V', from: 'klingai/kling-v3-omni-t2v/generation', to: 'klingai/kling-v3-omni-t2v/generation', color: 'bg-orange-100 text-orange-700 hover:bg-orange-200 dark:bg-orange-900/30 dark:text-orange-400' },
  { label: 'Sora', from: 'openai/sora/generation', to: 'openai/sora/generation', color: 'bg-slate-100 text-slate-700 hover:bg-slate-200 dark:bg-slate-900/30 dark:text-slate-400' },
  { label: 'Midjourney Video', from: 'midjourney/video/generation', to: 'midjourney/video/generation', color: 'bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-900/30 dark:text-violet-400' },
  // ── 图像生成 ──
  { label: 'Nano Banana Pro', from: 'google/nano-banana-pro/generation', to: 'google/nano-banana-pro/generation', color: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400' },
  { label: 'Nano Banana 2', from: 'google/nano-banana-2/generation', to: 'google/nano-banana-2/generation', color: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-400' },
  { label: 'Midjourney Diffusion', from: 'midjourney/diffusion/generation', to: 'midjourney/diffusion/generation', color: 'bg-violet-100 text-violet-700 hover:bg-violet-200 dark:bg-violet-900/30 dark:text-violet-400' },
  { label: 'GPT Image 2', from: 'openai/gpt-image-2/generation', to: 'openai/gpt-image-2/generation', color: 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-400' },
  // ── 音频生成 ──
  { label: 'MiniMax Music 2.5', from: 'minimax/music-2.5/generation', to: 'minimax/music-2.5/generation', color: 'bg-fuchsia-100 text-fuchsia-700 hover:bg-fuchsia-200 dark:bg-fuchsia-900/30 dark:text-fuchsia-400' },
  { label: 'MiniMax Speech HD', from: 'minimax/speech-2.8-hd/generation', to: 'minimax/speech-2.8-hd/generation', color: 'bg-fuchsia-100 text-fuchsia-700 hover:bg-fuchsia-200 dark:bg-fuchsia-900/30 dark:text-fuchsia-400' },
]

// Antigravity 默认映射（从后端 API 获取，与 constants.go 保持一致）
// 使用 fetchAntigravityDefaultMappings() 异步获取
import { getAntigravityDefaultModelMapping } from '@/api/admin/accounts'

let _antigravityDefaultMappingsCache: { from: string; to: string }[] | null = null

export async function fetchAntigravityDefaultMappings(): Promise<{ from: string; to: string }[]> {
  if (_antigravityDefaultMappingsCache !== null) {
    return _antigravityDefaultMappingsCache
  }
  try {
    const mapping = await getAntigravityDefaultModelMapping()
    _antigravityDefaultMappingsCache = Object.entries(mapping).map(([from, to]) => ({ from, to }))
  } catch (e) {
    console.warn('[fetchAntigravityDefaultMappings] API failed, using empty fallback', e)
    _antigravityDefaultMappingsCache = []
  }
  return _antigravityDefaultMappingsCache
}

// =====================
// 常用错误码
// =====================

export const commonErrorCodes = [
  { value: 401, label: 'Unauthorized' },
  { value: 403, label: 'Forbidden' },
  { value: 429, label: 'Rate Limit' },
  { value: 500, label: 'Server Error' },
  { value: 502, label: 'Bad Gateway' },
  { value: 503, label: 'Unavailable' },
  { value: 529, label: 'Overloaded' }
]

// =====================
// 辅助函数
// =====================

// 按平台获取模型
export function getModelsByPlatform(platform: string): string[] {
  switch (platform) {
    case 'openai': return openaiModels
    case 'anthropic':
    case 'claude': return claudeModels
    case 'gemini': return geminiModels
    case 'antigravity': return antigravityModels
    case 'mulerun': return mulerunModels
    case 'zhipu': return zhipuModels
    case 'qwen': return qwenModels
    case 'deepseek': return deepseekModels
    case 'mistral': return mistralModels
    case 'meta': return metaModels
    case 'xai': return xaiModels
    case 'cohere': return cohereModels
    case 'yi': return yiModels
    case 'moonshot': return moonshotModels
    case 'doubao': return doubaoModels
    case 'minimax': return minimaxModels
    case 'baidu': return baiduModels
    case 'spark': return sparkModels
    case 'hunyuan': return hunyuanModels
    case 'perplexity': return perplexityModels
    default: return claudeModels
  }
}

// 按平台获取预设映射
export function getPresetMappingsByPlatform(platform: string) {
  if (platform === 'openai') return openaiPresetMappings
  if (platform === 'gemini') return geminiPresetMappings
  if (platform === 'antigravity') return antigravityPresetMappings
  if (platform === 'mulerun') return mulerunPresetMappings
  if (platform === 'bedrock') return bedrockPresetMappings
  return anthropicPresetMappings
}

// =====================
// 构建模型映射对象（用于 API）
// =====================

// isValidWildcardPattern 校验通配符格式：* 只能放在末尾
// 导出供表单组件使用实时校验
export function isValidWildcardPattern(pattern: string): boolean {
  const starIndex = pattern.indexOf('*')
  if (starIndex === -1) return true // 无通配符，有效
  // * 必须在末尾，且只能有一个
  return starIndex === pattern.length - 1 && pattern.lastIndexOf('*') === starIndex
}

export type ModelRestrictionMode = 'whitelist' | 'mapping' | 'combined'

export interface ModelMappingEntry {
  from: string
  to: string
}

export function splitModelMappingObject(
  modelMapping?: Record<string, unknown> | null
): { allowedModels: string[]; modelMappings: ModelMappingEntry[] } {
  const allowedModels: string[] = []
  const modelMappings: ModelMappingEntry[] = []

  if (!modelMapping || typeof modelMapping !== 'object') {
    return { allowedModels, modelMappings }
  }

  for (const [rawFrom, rawTo] of Object.entries(modelMapping)) {
    if (typeof rawTo !== 'string') continue
    const from = rawFrom.trim()
    const to = rawTo.trim()
    if (!from || !to) continue

    if (from === to) {
      allowedModels.push(from)
    } else {
      modelMappings.push({ from, to })
    }
  }

  return { allowedModels, modelMappings }
}

export function buildModelMappingObject(
  mode: ModelRestrictionMode,
  allowedModels: string[],
  modelMappings: ModelMappingEntry[]
): Record<string, string> | null {
  const mapping: Record<string, string> = {}

  if (mode === 'whitelist' || mode === 'combined') {
    for (const model of allowedModels) {
      const normalizedModel = model.trim()
      if (!normalizedModel) continue
      // whitelist 模式的本意是"精确模型列表"，如果用户输入了通配符（如 claude-*），
      // 写入 model_mapping 会导致 GetMappedModel() 把真实模型映射成 "claude-*"，从而转发失败。
      // 因此这里跳过包含通配符的条目。
      if (!normalizedModel.includes('*')) {
        mapping[normalizedModel] = normalizedModel
      }
    }
  }

  if (mode === 'mapping' || mode === 'combined') {
    for (const m of modelMappings) {
      const from = m.from.trim()
      const to = m.to.trim()
      if (!from || !to) continue
      // 校验通配符格式：* 只能放在末尾
      if (!isValidWildcardPattern(from)) {
        console.warn(`[buildModelMappingObject] 无效的通配符格式，跳过: ${from}`)
        continue
      }
      // to 不允许包含通配符
      if (to.includes('*')) {
        console.warn(`[buildModelMappingObject] 目标模型不能包含通配符，跳过: ${from} -> ${to}`)
        continue
      }
      mapping[from] = to
    }
  }

  return Object.keys(mapping).length > 0 ? mapping : null
}
