package ratio_setting

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
)

// from songquanpeng/one-api
const (
	USD2RMB = 7.3 // 暂定 1 USD = 7.3 RMB
	USD     = 500 // $0.002 = 1 -> $1 = 500
	RMB     = USD / USD2RMB
)

// modelRatio
// https://platform.openai.com/docs/models/model-endpoint-compatibility
// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Blfmc9dlf
// https://openai.com/pricing
// TODO: when a new api is enabled, check the pricing here
// 1 === $0.002 / 1K tokens
// 1 === ￥0.014 / 1k tokens

var defaultModelRatio = map[string]float64{
	//"midjourney":                50,
	"gpt-4-gizmo-*":  15,
	"gpt-4o-gizmo-*": 2.5,
	"gpt-4-all":      15,
	"gpt-4o-all":     15,
	"gpt-4":          15,
	//"gpt-4-0314":                   15, //deprecated
	"gpt-4-0613": 15,
	"gpt-4-32k":  30,
	//"gpt-4-32k-0314":               30, //deprecated
	"gpt-4-32k-0613":                          30,
	"gpt-4-1106-preview":                      5,    // $10 / 1M tokens
	"gpt-4-0125-preview":                      5,    // $10 / 1M tokens
	"gpt-4-turbo-preview":                     5,    // $10 / 1M tokens
	"gpt-4-vision-preview":                    5,    // $10 / 1M tokens
	"gpt-4-1106-vision-preview":               5,    // $10 / 1M tokens
	"chatgpt-4o-latest":                       2.5,  // $5 / 1M tokens
	"gpt-4o":                                  1.25, // $2.5 / 1M tokens
	"gpt-4o-audio-preview":                    1.25, // $2.5 / 1M tokens
	"gpt-4o-audio-preview-2024-10-01":         1.25, // $2.5 / 1M tokens
	"gpt-4o-2024-05-13":                       2.5,  // $5 / 1M tokens
	"gpt-4o-2024-08-06":                       1.25, // $2.5 / 1M tokens
	"gpt-4o-2024-11-20":                       1.25, // $2.5 / 1M tokens
	"gpt-4o-realtime-preview":                 2.5,
	"gpt-4o-realtime-preview-2024-10-01":      2.5,
	"gpt-4o-realtime-preview-2024-12-17":      2.5,
	"gpt-4o-mini-realtime-preview":            0.3,
	"gpt-4o-mini-realtime-preview-2024-12-17": 0.3,
	"gpt-4.1":                          1.0,  // $2 / 1M tokens
	"gpt-4.1-2025-04-14":               1.0,  // $2 / 1M tokens
	"gpt-4.1-mini":                     0.2,  // $0.4 / 1M tokens
	"gpt-4.1-mini-2025-04-14":          0.2,  // $0.4 / 1M tokens
	"gpt-4.1-nano":                     0.05, // $0.1 / 1M tokens
	"gpt-4.1-nano-2025-04-14":          0.05, // $0.1 / 1M tokens
	"gpt-image-1":                      2.5,  // $5 / 1M tokens
	"o1":                               7.5,  // $15 / 1M tokens
	"o1-2024-12-17":                    7.5,  // $15 / 1M tokens
	"o1-preview":                       7.5,  // $15 / 1M tokens
	"o1-preview-2024-09-12":            7.5,  // $15 / 1M tokens
	"o1-mini":                          0.55, // $1.1 / 1M tokens
	"o1-mini-2024-09-12":               0.55, // $1.1 / 1M tokens
	"o1-pro":                           75.0, // $150 / 1M tokens
	"o1-pro-2025-03-19":                75.0, // $150 / 1M tokens
	"o3-mini":                          0.55,
	"o3-mini-2025-01-31":               0.55,
	"o3-mini-high":                     0.55,
	"o3-mini-2025-01-31-high":          0.55,
	"o3-mini-low":                      0.55,
	"o3-mini-2025-01-31-low":           0.55,
	"o3-mini-medium":                   0.55,
	"o3-mini-2025-01-31-medium":        0.55,
	"o3":                               1.0,  // $2 / 1M tokens
	"o3-2025-04-16":                    1.0,  // $2 / 1M tokens
	"o3-pro":                           10.0, // $20 / 1M tokens
	"o3-pro-2025-06-10":                10.0, // $20 / 1M tokens
	"o3-deep-research":                 5.0,  // $10 / 1M tokens
	"o3-deep-research-2025-06-26":      5.0,  // $10 / 1M tokens
	"o4-mini":                          0.55, // $1.1 / 1M tokens
	"o4-mini-2025-04-16":               0.55, // $1.1 / 1M tokens
	"o4-mini-deep-research":            1.0,  // $2 / 1M tokens
	"o4-mini-deep-research-2025-06-26": 1.0,  // $2 / 1M tokens
	"gpt-4o-mini":                      0.075,
	"gpt-4o-mini-2024-07-18":           0.075,
	"gpt-4-turbo":                      5, // $0.01 / 1K tokens
	"gpt-4-turbo-2024-04-09":           5, // $0.01 / 1K tokens
	"gpt-4.5-preview":                  37.5,
	"gpt-4.5-preview-2025-02-27":       37.5,
	"gpt-5":                            0.625,
	"gpt-5-2025-08-07":                 0.625,
	"gpt-5-chat-latest":                0.625,
	"gpt-5-mini":                       0.125,
	"gpt-5-mini-2025-08-07":            0.125,
	"gpt-5-nano":                       0.025,
	"gpt-5-nano-2025-08-07":            0.025,
	//"gpt-3.5-turbo-0301":           0.75, //deprecated
	"gpt-3.5-turbo":          0.25,
	"gpt-3.5-turbo-0613":     0.75,
	"gpt-3.5-turbo-16k":      1.5, // $0.003 / 1K tokens
	"gpt-3.5-turbo-16k-0613": 1.5,
	"gpt-3.5-turbo-instruct": 0.75, // $0.0015 / 1K tokens
	"gpt-3.5-turbo-1106":     0.5,  // $0.001 / 1K tokens
	"gpt-3.5-turbo-0125":     0.25,
	"babbage-002":            0.2, // $0.0004 / 1K tokens
	"davinci-002":            1,   // $0.002 / 1K tokens
	"text-ada-001":           0.2,
	"text-babbage-001":       0.25,
	"text-curie-001":         1,
	//"text-davinci-002":               10,
	//"text-davinci-003":               10,
	"text-davinci-edit-001":                     10,
	"code-davinci-edit-001":                     10,
	"whisper-1":                                 15,  // $0.006 / minute -> $0.006 / 150 words -> $0.006 / 200 tokens -> $0.03 / 1k tokens
	"tts-1":                                     7.5, // 1k characters -> $0.015
	"tts-1-1106":                                7.5, // 1k characters -> $0.015
	"tts-1-hd":                                  15,  // 1k characters -> $0.03
	"tts-1-hd-1106":                             15,  // 1k characters -> $0.03
	"davinci":                                   10,
	"curie":                                     10,
	"babbage":                                   10,
	"ada":                                       10,
	"text-embedding-3-small":                    0.01,
	"text-embedding-3-large":                    0.065,
	"text-embedding-ada-002":                    0.05,
	"text-search-ada-doc-001":                   10,
	"text-moderation-stable":                    0.1,
	"text-moderation-latest":                    0.1,
	"claude-instant-1":                          0.4,   // $0.8 / 1M tokens
	"claude-2.0":                                4,     // $8 / 1M tokens
	"claude-2.1":                                4,     // $8 / 1M tokens
	"claude-3-haiku-20240307":                   0.125, // $0.25 / 1M tokens
	"claude-3-5-haiku-20241022":                 0.5,   // $1 / 1M tokens
	"claude-haiku-4-5-20251001":                 0.5,   // $1 / 1M tokens
	"claude-3-sonnet-20240229":                  1.5,   // $3 / 1M tokens
	"claude-3-5-sonnet-20240620":                1.5,
	"claude-3-5-sonnet-20241022":                1.5,
	"claude-3-7-sonnet-20250219":                1.5,
	"claude-3-7-sonnet-20250219-thinking":       1.5,
	"claude-sonnet-4-20250514":                  1.5,
	"claude-sonnet-4-5-20250929":                1.5,
	"claude-opus-4-5-20251101":                  2.5,
	"claude-3-opus-20240229":                    7.5, // $15 / 1M tokens
	"claude-opus-4-20250514":                    7.5,
	"claude-opus-4-1-20250805":                  7.5,
	"ERNIE-4.0-8K":                              0.120 * RMB,
	"ERNIE-3.5-8K":                              0.012 * RMB,
	"ERNIE-3.5-8K-0205":                         0.024 * RMB,
	"ERNIE-3.5-8K-1222":                         0.012 * RMB,
	"ERNIE-Bot-8K":                              0.024 * RMB,
	"ERNIE-3.5-4K-0205":                         0.012 * RMB,
	"ERNIE-Speed-8K":                            0.004 * RMB,
	"ERNIE-Speed-128K":                          0.004 * RMB,
	"ERNIE-Lite-8K-0922":                        0.008 * RMB,
	"ERNIE-Lite-8K-0308":                        0.003 * RMB,
	"ERNIE-Tiny-8K":                             0.001 * RMB,
	"BLOOMZ-7B":                                 0.004 * RMB,
	"Embedding-V1":                              0.002 * RMB,
	"bge-large-zh":                              0.002 * RMB,
	"bge-large-en":                              0.002 * RMB,
	"tao-8k":                                    0.002 * RMB,
	"PaLM-2":                                    1,
	"gemini-1.5-pro-latest":                     1.25, // $3.5 / 1M tokens
	"gemini-1.5-flash-latest":                   0.075,
	"gemini-2.0-flash":                          0.05,
	"gemini-2.5-pro-exp-03-25":                  0.625,
	"gemini-2.5-pro-preview-03-25":              0.625,
	"gemini-2.5-pro":                            0.625,
	"gemini-2.5-flash-preview-04-17":            0.075,
	"gemini-2.5-flash-preview-04-17-thinking":   0.075,
	"gemini-2.5-flash-preview-04-17-nothinking": 0.075,
	"gemini-2.5-flash-preview-05-20":            0.075,
	"gemini-2.5-flash-preview-05-20-thinking":   0.075,
	"gemini-2.5-flash-preview-05-20-nothinking": 0.075,
	"gemini-2.5-flash-thinking-*":               0.075, // 用于为后续所有2.5 flash thinking budget 模型设置默认倍率
	"gemini-2.5-pro-thinking-*":                 0.625, // 用于为后续所有2.5 pro thinking budget 模型设置默认倍率
	"gemini-2.5-flash-lite-preview-thinking-*":  0.05,
	"gemini-2.5-flash-lite-preview-06-17":       0.05,
	"gemini-2.5-flash":                          0.15,
	"gemini-robotics-er-1.5-preview":            0.15,
	"gemini-embedding-001":                      0.075,
	"text-embedding-004":                        0.001,
	"chatglm_turbo":                             0.3572,     // ￥0.005 / 1k tokens
	"chatglm_pro":                               0.7143,     // ￥0.01 / 1k tokens
	"chatglm_std":                               0.3572,     // ￥0.005 / 1k tokens
	"chatglm_lite":                              0.1429,     // ￥0.002 / 1k tokens
	"glm-4":                                     7.143,      // ￥0.1 / 1k tokens
	"glm-4v":                                    0.05 * RMB, // ￥0.05 / 1k tokens
	"glm-4-alltools":                            0.1 * RMB,  // ￥0.1 / 1k tokens
	"glm-3-turbo":                               0.3572,
	"glm-4-plus":                                0.05 * RMB,
	"glm-4-0520":                                0.1 * RMB,
	"glm-4-air":                                 0.001 * RMB,
	"glm-4-airx":                                0.01 * RMB,
	"glm-4-long":                                0.001 * RMB,
	"glm-4-flash":                               0,
	"glm-4v-plus":                               0.01 * RMB,
	// ========== 阿里云通义千问 (价格单位: ¥/1M tokens, 转为 ¥/1K tokens * RMB) ==========
	// --- 旗舰模型 ---
	"qwen3-max":                    0.0025 * RMB,  // ¥2.5/1M tokens (≤32K)
	"qwen3-max-2026-01-23":         0.0025 * RMB,  // ¥2.5/1M tokens
	"qwen3-max-preview":            0.006 * RMB,   // ¥6/1M tokens
	"qwen-max":                     0.0024 * RMB,  // ¥2.4/1M tokens
	"qwen-max-latest":              0.0024 * RMB,  // ¥2.4/1M tokens
	// --- 中端模型 ---
	"qwen-plus":                    0.0008 * RMB,  // ¥0.8/1M tokens (≤128K)
	"qwen-plus-latest":             0.0008 * RMB,  // ¥0.8/1M tokens
	// --- 轻量模型 ---
	"qwen-flash":                   0.00015 * RMB, // ¥0.15/1M tokens
	"qwen-turbo":                   0.0003 * RMB,  // ¥0.3/1M tokens
	"qwen-turbo-latest":            0.0003 * RMB,  // ¥0.3/1M tokens
	// --- 长文本 ---
	"qwen-long":                    0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen-long-latest":             0.0005 * RMB,  // ¥0.5/1M tokens
	// --- 深度研究 ---
	"qwen-deep-research":           0.054 * RMB,   // ¥54/1M tokens
	// --- Qwen3 开源系列 ---
	"qwen3-235b-a22b":              0.002 * RMB,   // ¥2/1M tokens
	"qwen3-32b":                    0.002 * RMB,   // ¥2/1M tokens
	"qwen3-30b-a3b":                0.00075 * RMB, // ¥0.75/1M tokens
	"qwen3-14b":                    0.001 * RMB,   // ¥1/1M tokens
	"qwen3-8b":                     0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen3-4b":                     0.0003 * RMB,  // ¥0.3/1M tokens
	"qwen3-1.7b":                   0.0003 * RMB,  // ¥0.3/1M tokens
	"qwen3-0.6b":                   0.0003 * RMB,  // ¥0.3/1M tokens
	// --- Qwen3 thinking/instruct 变体 ---
	"qwen3-235b-a22b-thinking-2507":  0.002 * RMB,  // ¥2/1M tokens
	"qwen3-235b-a22b-instruct-2507":  0.002 * RMB,  // ¥2/1M tokens
	"qwen3-30b-a3b-thinking-2507":    0.00075 * RMB, // ¥0.75/1M tokens
	"qwen3-30b-a3b-instruct-2507":    0.00075 * RMB, // ¥0.75/1M tokens
	"qwen3-next-80b-a3b-thinking":    0.001 * RMB,  // ¥1/1M tokens
	"qwen3-next-80b-a3b-instruct":    0.001 * RMB,  // ¥1/1M tokens
	// --- Vision 主力模型 ---
	"qwen3-vl-plus":                0.001 * RMB,   // ¥1/1M tokens
	"qwen3-vl-flash":               0.00015 * RMB, // ¥0.15/1M tokens
	"qwen-vl-max":                  0.0016 * RMB,  // ¥1.6/1M tokens
	"qwen-vl-max-latest":           0.0016 * RMB,  // ¥1.6/1M tokens
	"qwen-vl-plus":                 0.0008 * RMB,  // ¥0.8/1M tokens
	"qwen-vl-plus-latest":          0.0008 * RMB,  // ¥0.8/1M tokens
	"qwen-vl-ocr":                  0.005 * RMB,   // ¥5/1M tokens
	"qwen-vl-ocr-latest":           0.0003 * RMB,  // ¥0.3/1M tokens
	// --- Vision 开源系列 ---
	"qwen3-vl-235b-a22b-thinking":  0.002 * RMB,   // ¥2/1M tokens
	"qwen3-vl-235b-a22b-instruct":  0.002 * RMB,   // ¥2/1M tokens
	"qwen3-vl-32b-thinking":        0.002 * RMB,   // ¥2/1M tokens
	"qwen3-vl-32b-instruct":        0.002 * RMB,   // ¥2/1M tokens
	"qwen3-vl-30b-a3b-thinking":    0.00075 * RMB, // ¥0.75/1M tokens
	"qwen3-vl-30b-a3b-instruct":    0.00075 * RMB, // ¥0.75/1M tokens
	"qwen3-vl-8b-thinking":         0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen3-vl-8b-instruct":         0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen2.5-vl-72b-instruct":      0.016 * RMB,   // ¥16/1M tokens
	"qwen2.5-vl-32b-instruct":      0.008 * RMB,   // ¥8/1M tokens
	"qwen2.5-vl-7b-instruct":       0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-vl-3b-instruct":       0.0012 * RMB,  // ¥1.2/1M tokens
	// --- Coder ---
	"qwen3-coder-plus":             0.004 * RMB,   // ¥4/1M tokens
	"qwen3-coder-flash":            0.001 * RMB,   // ¥1/1M tokens
	"qwen3-coder-480b-a35b-instruct": 0.006 * RMB, // ¥6/1M tokens
	"qwen3-coder-30b-a3b-instruct":   0.0015 * RMB, // ¥1.5/1M tokens
	"qwen-coder-plus":              0.0035 * RMB,  // ¥3.5/1M tokens
	"qwen-coder-plus-latest":       0.0035 * RMB,  // ¥3.5/1M tokens
	"qwen-coder-turbo":             0.002 * RMB,   // ¥2/1M tokens
	"qwen-coder-turbo-latest":      0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-coder-32b-instruct":   0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-coder-14b-instruct":   0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-coder-7b-instruct":    0.001 * RMB,   // ¥1/1M tokens
	// --- Math ---
	"qwen-math-plus":               0.004 * RMB,   // ¥4/1M tokens
	"qwen-math-turbo":              0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-math-72b-instruct":    0.004 * RMB,   // ¥4/1M tokens
	"qwen2.5-math-7b-instruct":     0.001 * RMB,   // ¥1/1M tokens
	// --- Qwen2.5 开源系列 ---
	"qwen2.5-72b-instruct":         0.004 * RMB,   // ¥4/1M tokens
	"qwen2.5-32b-instruct":         0.002 * RMB,   // ¥2/1M tokens
	"qwen2.5-14b-instruct":         0.001 * RMB,   // ¥1/1M tokens
	"qwen2.5-14b-instruct-1m":      0.001 * RMB,   // ¥1/1M tokens
	"qwen2.5-7b-instruct":          0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen2.5-7b-instruct-1m":       0.0005 * RMB,  // ¥0.5/1M tokens
	"qwen2.5-3b-instruct":          0.0003 * RMB,  // ¥0.3/1M tokens
	"text-embedding-v1":                         0.05,   // ￥0.0007 / 1k tokens
	"SparkDesk-v1.1":                            1.2858, // ￥0.018 / 1k tokens
	"SparkDesk-v2.1":                            1.2858, // ￥0.018 / 1k tokens
	"SparkDesk-v3.1":                            1.2858, // ￥0.018 / 1k tokens
	"SparkDesk-v3.5":                            1.2858, // ￥0.018 / 1k tokens
	"SparkDesk-v4.0":                            1.2858,
	"360GPT_S2_V9":                              0.8572, // ¥0.012 / 1k tokens
	"360gpt-turbo":                              0.0858, // ¥0.0012 / 1k tokens
	"360gpt-turbo-responsibility-8k":            0.8572, // ¥0.012 / 1k tokens
	"360gpt-pro":                                0.8572, // ¥0.012 / 1k tokens
	"360gpt2-pro":                               0.8572, // ¥0.012 / 1k tokens
	"embedding-bert-512-v1":                     0.0715, // ¥0.001 / 1k tokens
	"embedding_s1_v1":                           0.0715, // ¥0.001 / 1k tokens
	"semantic_similarity_s1_v1":                 0.0715, // ¥0.001 / 1k tokens
	"hunyuan":                                   7.143,  // ¥0.1 / 1k tokens  // https://cloud.tencent.com/document/product/1729/97731#e0e6be58-60c8-469f-bdeb-6c264ce3b4d0
	// https://platform.lingyiwanwu.com/docs#-计费单元
	// 已经按照 7.2 来换算美元价格
	"yi-34b-chat-0205":       0.18,
	"yi-34b-chat-200k":       0.864,
	"yi-vl-plus":             0.432,
	"yi-large":               20.0 / 1000 * RMB,
	"yi-medium":              2.5 / 1000 * RMB,
	"yi-vision":              6.0 / 1000 * RMB,
	"yi-medium-200k":         12.0 / 1000 * RMB,
	"yi-spark":               1.0 / 1000 * RMB,
	"yi-large-rag":           25.0 / 1000 * RMB,
	"yi-large-turbo":         12.0 / 1000 * RMB,
	"yi-large-preview":       20.0 / 1000 * RMB,
	"yi-large-rag-preview":   25.0 / 1000 * RMB,
	"command":                0.5,
	"command-nightly":        0.5,
	"command-light":          0.5,
	"command-light-nightly":  0.5,
	"command-r":              0.25,
	"command-r-plus":         1.5,
	"command-r-08-2024":      0.075,
	"command-r-plus-08-2024": 1.25,
	"deepseek-chat":          0.27 / 2,
	"deepseek-coder":         0.27 / 2,
	"deepseek-reasoner":      0.55 / 2, // 0.55 / 1k tokens
	// Perplexity online 模型对搜索额外收费，有需要应自行调整，此处不计入搜索费用
	"llama-3-sonar-small-32k-chat":   0.2 / 1000 * USD,
	"llama-3-sonar-small-32k-online": 0.2 / 1000 * USD,
	"llama-3-sonar-large-32k-chat":   1 / 1000 * USD,
	"llama-3-sonar-large-32k-online": 1 / 1000 * USD,
	// grok
	"grok-3-beta":           1.5,
	"grok-3-mini-beta":      0.15,
	"grok-2":                1,
	"grok-2-vision":         1,
	"grok-beta":             2.5,
	"grok-vision-beta":      2.5,
	"grok-3-fast-beta":      2.5,
	"grok-3-mini-fast-beta": 0.3,
	// submodel
	"NousResearch/Hermes-4-405B-FP8":          0.8,
	"Qwen/Qwen3-235B-A22B-Thinking-2507":      0.6,
	"Qwen/Qwen3-Coder-480B-A35B-Instruct-FP8": 0.8,
	"Qwen/Qwen3-235B-A22B-Instruct-2507":      0.3,
	"zai-org/GLM-4.5-FP8":                     0.8,
	"openai/gpt-oss-120b":                     0.5,
	"deepseek-ai/DeepSeek-R1-0528":            0.8,
	"deepseek-ai/DeepSeek-R1":                 0.8,
	"deepseek-ai/DeepSeek-V3-0324":            0.8,
	"deepseek-ai/DeepSeek-V3.1":               0.8,
}

var defaultModelPrice = map[string]float64{
	"suno_music":                     0.1,
	"suno_lyrics":                    0.01,
	"dall-e-3":                       0.04,
	"imagen-3.0-generate-002":        0.03,
	"black-forest-labs/flux-1.1-pro": 0.04,
	"gpt-4-gizmo-*":                  0.1,
	"mj_video":                       0.8,
	// 阿里云万相视频模型定价 (按秒计费)
	// 实际价格 = 基础价格 * 秒数 * 分辨率倍率
	// 分辨率倍率在 ProcessAliOtherRatios 中配置
	// ========== 图生视频 (i2v) ==========
	"wan2.6-i2v-flash":   0.6 / 7.2,  // ¥0.6/秒 (720P有声) -> 转换为美元
	"wan2.6-i2v":         0.6 / 7.2,  // ¥0.6/秒 (720P有声)
	"wan2.5-i2v-preview": 0.3 / 7.2,  // ¥0.3/秒 (480P有声)
	"wan2.2-i2v-flash":   0.1 / 7.2,  // ¥0.1/秒 (480P无声)
	"wan2.2-i2v-plus":    0.14 / 7.2, // ¥0.14/秒 (480P无声)
	"wanx2.1-i2v-plus":   0.6 / 7.2,  // ¥0.6/秒 (720P无声)
	"wanx2.1-i2v-turbo":  0.1 / 7.2,  // ¥0.1/秒 (480P无声)
	// ========== 文生视频 (t2v) ==========
	"wan2.6-t2v":         0.6 / 7.2,  // ¥0.6/秒 (720P有声)
	"wan2.5-t2v-preview": 0.3 / 7.2,  // ¥0.3/秒 (480P有声)
	"wan2.2-t2v-plus":    0.14 / 7.2, // ¥0.14/秒 (480P无声)
	"wanx2.1-t2v-plus":   0.6 / 7.2,  // ¥0.6/秒 (720P无声)
	"wanx2.1-t2v-turbo":  0.1 / 7.2,  // ¥0.1/秒 (480P无声)
	// ========== 文生图 (t2i) ==========
	"wanx2.1-t2i-turbo": 0.14 / 7.2, // ¥0.14/张 (512*1024以下基础价，按分辨率倍率计费)
	"wanx2.1-t2i-plus":  0.14 / 7.2, // ¥0.14/张 (512*1024以下基础价)
	"wanx-v1":           0.08 / 7.2, // ¥0.08/张 (1024*1024)
	"flux-schnell":      0.06 / 7.2, // ¥0.06/张
	"flux-dev":          0.12 / 7.2, // ¥0.12/张
	"flux-merged":       0.12 / 7.2, // ¥0.12/张
	// ========== 图片编辑 (Image Edit) ==========
	"wanx2.1-imageedit-v1":       0.24 / 7.2, // ¥0.24/张
	"wanx2.0-imageedit-plus":     0.18 / 7.2, // ¥0.18/张
	"wanx-style-repaint-v1":      0.16 / 7.2, // ¥0.16/张
	"wanx-background-generation": 0.12 / 7.2, // ¥0.12/张
	"wanx-sketch-to-image-v1":    0.12 / 7.2, // ¥0.12/张
	"flux-inpaint":               0.12 / 7.2, // ¥0.12/张
	"flux-redux":                 0.06 / 7.2, // ¥0.06/张
	"mj_imagine":                     0.1,
	"mj_edits":                       0.1,
	"mj_variation":                   0.1,
	"mj_reroll":                      0.1,
	"mj_blend":                       0.1,
	"mj_modal":                       0.1,
	"mj_zoom":                        0.1,
	"mj_shorten":                     0.1,
	"mj_high_variation":              0.1,
	"mj_low_variation":               0.1,
	"mj_pan":                         0.1,
	"mj_inpaint":                     0,
	"mj_custom_zoom":                 0,
	"mj_describe":                    0.05,
	"mj_upscale":                     0.05,
	"swap_face":                      0.05,
	"mj_upload":                      0.05,
	"sora-2":                         0.3,
	"sora-2-pro":                     0.5,
	"gpt-4o-mini-tts":                0.3,
}

var defaultAudioRatio = map[string]float64{
	"gpt-4o-audio-preview":         16,
	"gpt-4o-mini-audio-preview":    66.67,
	"gpt-4o-realtime-preview":      8,
	"gpt-4o-mini-realtime-preview": 16.67,
	"gpt-4o-mini-tts":              25,
}

var defaultAudioCompletionRatio = map[string]float64{
	"gpt-4o-realtime":      2,
	"gpt-4o-mini-realtime": 2,
	"gpt-4o-mini-tts":      1,
	"tts-1":                0,
	"tts-1-hd":             0,
	"tts-1-1106":           0,
	"tts-1-hd-1106":        0,
}

var (
	modelPriceMap      map[string]float64 = nil
	modelPriceMapMutex                    = sync.RWMutex{}
)
var (
	modelRatioMap      map[string]float64 = nil
	modelRatioMapMutex                    = sync.RWMutex{}
)

var (
	CompletionRatio      map[string]float64 = nil
	CompletionRatioMutex                    = sync.RWMutex{}
)

var defaultCompletionRatio = map[string]float64{
	"gpt-4-gizmo-*":  2,
	"gpt-4o-gizmo-*": 3,
	"gpt-4-all":      2,
	"gpt-image-1":    8,
	// ========== 阿里云通义千问 CompletionRatio (输出价/输入价) ==========
	// 旗舰模型
	"qwen3-max":                      4,    // ¥10/¥2.5
	"qwen3-max-2026-01-23":           4,
	"qwen3-max-preview":              4,    // ¥24/¥6
	"qwen-max":                       4,    // ¥9.6/¥2.4
	"qwen-max-latest":                4,
	// 中端模型
	"qwen-plus":                      2.5,  // ¥2/¥0.8 (非思考输出)
	"qwen-plus-latest":               2.5,
	// 轻量模型
	"qwen-flash":                     10,   // ¥1.5/¥0.15
	"qwen-turbo":                     2,    // ¥0.6/¥0.3 (非思考输出)
	"qwen-turbo-latest":              2,
	// 长文本 & 深度研究
	"qwen-long":                      4,    // ¥2/¥0.5
	"qwen-long-latest":               4,
	"qwen-deep-research":             3.019, // ¥163/¥54
	// Qwen3 开源系列 (非思考输出)
	"qwen3-235b-a22b":                4,    // ¥8/¥2
	"qwen3-32b":                      4,    // ¥8/¥2
	"qwen3-30b-a3b":                  4,    // ¥3/¥0.75
	"qwen3-14b":                      4,    // ¥4/¥1
	"qwen3-8b":                       4,    // ¥2/¥0.5
	"qwen3-4b":                       4,    // ¥1.2/¥0.3
	"qwen3-1.7b":                     4,
	"qwen3-0.6b":                     4,
	// Qwen3 thinking/instruct 变体
	"qwen3-235b-a22b-thinking-2507":  10,   // ¥20/¥2 (思考输出)
	"qwen3-235b-a22b-instruct-2507":  4,    // ¥8/¥2
	"qwen3-30b-a3b-thinking-2507":    10,   // ¥7.5/¥0.75
	"qwen3-30b-a3b-instruct-2507":    4,    // ¥3/¥0.75
	"qwen3-next-80b-a3b-thinking":    10,   // ¥10/¥1
	"qwen3-next-80b-a3b-instruct":    4,    // ¥4/¥1
	// Vision 主力模型
	"qwen3-vl-plus":                  10,   // ¥10/¥1
	"qwen3-vl-flash":                 10,   // ¥1.5/¥0.15
	"qwen-vl-max":                    2.5,  // ¥4/¥1.6
	"qwen-vl-max-latest":             2.5,
	"qwen-vl-plus":                   2.5,  // ¥2/¥0.8
	"qwen-vl-plus-latest":            2.5,
	"qwen-vl-ocr":                    1,    // ¥5/¥5
	"qwen-vl-ocr-latest":             1.667, // ¥0.5/¥0.3
	// Vision 开源系列
	"qwen3-vl-235b-a22b-thinking":    10,   // ¥20/¥2
	"qwen3-vl-235b-a22b-instruct":    4,    // ¥8/¥2
	"qwen3-vl-32b-thinking":          10,   // ¥20/¥2
	"qwen3-vl-32b-instruct":          4,    // ¥8/¥2
	"qwen3-vl-30b-a3b-thinking":      10,   // ¥7.5/¥0.75
	"qwen3-vl-30b-a3b-instruct":      4,    // ¥3/¥0.75
	"qwen3-vl-8b-thinking":           10,   // ¥5/¥0.5
	"qwen3-vl-8b-instruct":           4,    // ¥2/¥0.5
	"qwen2.5-vl-72b-instruct":        3,    // ¥48/¥16
	"qwen2.5-vl-32b-instruct":        3,    // ¥24/¥8
	"qwen2.5-vl-7b-instruct":         2.5,  // ¥5/¥2
	"qwen2.5-vl-3b-instruct":         3,    // ¥3.6/¥1.2
	// Coder
	"qwen3-coder-plus":               4,    // ¥16/¥4
	"qwen3-coder-flash":              4,    // ¥4/¥1
	"qwen3-coder-480b-a35b-instruct": 4,    // ¥24/¥6
	"qwen3-coder-30b-a3b-instruct":   4,    // ¥6/¥1.5
	"qwen-coder-plus":                2,    // ¥7/¥3.5
	"qwen-coder-plus-latest":         2,
	"qwen-coder-turbo":               3,    // ¥6/¥2
	"qwen-coder-turbo-latest":        3,
	"qwen2.5-coder-32b-instruct":     3,    // ¥6/¥2
	"qwen2.5-coder-14b-instruct":     3,
	"qwen2.5-coder-7b-instruct":      2,    // ¥2/¥1
	// Math
	"qwen-math-plus":                 3,    // ¥12/¥4
	"qwen-math-turbo":                3,    // ¥6/¥2
	"qwen2.5-math-72b-instruct":      3,
	"qwen2.5-math-7b-instruct":       2,    // ¥2/¥1
	// Qwen2.5 开源系列
	"qwen2.5-72b-instruct":           3,    // ¥12/¥4
	"qwen2.5-32b-instruct":           3,    // ¥6/¥2
	"qwen2.5-14b-instruct":           3,    // ¥3/¥1
	"qwen2.5-14b-instruct-1m":        3,
	"qwen2.5-7b-instruct":            2,    // ¥1/¥0.5
	"qwen2.5-7b-instruct-1m":         2,
	"qwen2.5-3b-instruct":            3,    // ¥0.9/¥0.3
}

// ========== 阶梯计价配置 ==========
// 阿里云部分模型按单次请求输入 Token 数量分阶梯定价
// 价格单位: 元/百万tokens, 转为 ratio (¥/1K tokens * RMB)

type PriceTier struct {
	MaxTokens       int     // Token 阈值，0 表示无上限
	InputRatio      float64 // 该阶梯的输入倍率 (替代 modelRatio)
	CompletionRatio float64 // 该阶梯的输出/输入倍率 (替代 completionRatio)
}

var defaultTieredPricing = map[string][]PriceTier{
	// qwen3-max: ≤32K / ≤128K / ≤252K
	"qwen3-max": {
		{MaxTokens: 32000, InputRatio: 0.0025 * RMB, CompletionRatio: 4},   // ¥2.5/10
		{MaxTokens: 128000, InputRatio: 0.004 * RMB, CompletionRatio: 4},   // ¥4/16
		{MaxTokens: 0, InputRatio: 0.007 * RMB, CompletionRatio: 4},        // ¥7/28
	},
	"qwen3-max-2026-01-23": {
		{MaxTokens: 32000, InputRatio: 0.0025 * RMB, CompletionRatio: 4},
		{MaxTokens: 128000, InputRatio: 0.004 * RMB, CompletionRatio: 4},
		{MaxTokens: 0, InputRatio: 0.007 * RMB, CompletionRatio: 4},
	},
	// qwen3-max-preview: ≤32K / ≤128K / ≤252K
	"qwen3-max-preview": {
		{MaxTokens: 32000, InputRatio: 0.006 * RMB, CompletionRatio: 4},    // ¥6/24
		{MaxTokens: 128000, InputRatio: 0.010 * RMB, CompletionRatio: 4},   // ¥10/40
		{MaxTokens: 0, InputRatio: 0.015 * RMB, CompletionRatio: 4},        // ¥15/60
	},
	// qwen-plus: ≤128K / ≤256K / ≤1M (非思考输出)
	"qwen-plus": {
		{MaxTokens: 128000, InputRatio: 0.0008 * RMB, CompletionRatio: 2.5},  // ¥0.8/2
		{MaxTokens: 256000, InputRatio: 0.0024 * RMB, CompletionRatio: 8.333}, // ¥2.4/20
		{MaxTokens: 0, InputRatio: 0.0048 * RMB, CompletionRatio: 10},        // ¥4.8/48
	},
	"qwen-plus-latest": {
		{MaxTokens: 128000, InputRatio: 0.0008 * RMB, CompletionRatio: 2.5},
		{MaxTokens: 256000, InputRatio: 0.0024 * RMB, CompletionRatio: 8.333},
		{MaxTokens: 0, InputRatio: 0.0048 * RMB, CompletionRatio: 10},
	},
	// qwen-flash: ≤128K / ≤256K / ≤1M
	"qwen-flash": {
		{MaxTokens: 128000, InputRatio: 0.00015 * RMB, CompletionRatio: 10},  // ¥0.15/1.5
		{MaxTokens: 256000, InputRatio: 0.0006 * RMB, CompletionRatio: 10},   // ¥0.6/6
		{MaxTokens: 0, InputRatio: 0.0012 * RMB, CompletionRatio: 10},        // ¥1.2/12
	},
	// qwen3-vl-plus: ≤32K / ≤128K / ≤256K
	"qwen3-vl-plus": {
		{MaxTokens: 32000, InputRatio: 0.001 * RMB, CompletionRatio: 10},     // ¥1/10
		{MaxTokens: 128000, InputRatio: 0.0015 * RMB, CompletionRatio: 10},   // ¥1.5/15
		{MaxTokens: 0, InputRatio: 0.003 * RMB, CompletionRatio: 10},         // ¥3/30
	},
	// qwen3-vl-flash: ≤32K / ≤128K / ≤256K
	"qwen3-vl-flash": {
		{MaxTokens: 32000, InputRatio: 0.00015 * RMB, CompletionRatio: 10},   // ¥0.15/1.5
		{MaxTokens: 128000, InputRatio: 0.0003 * RMB, CompletionRatio: 10},   // ¥0.3/3
		{MaxTokens: 0, InputRatio: 0.0006 * RMB, CompletionRatio: 10},        // ¥0.6/6
	},
	// qwen3-coder-plus: ≤32K / ≤128K / ≤256K / ≤1M
	"qwen3-coder-plus": {
		{MaxTokens: 32000, InputRatio: 0.004 * RMB, CompletionRatio: 4},      // ¥4/16
		{MaxTokens: 128000, InputRatio: 0.006 * RMB, CompletionRatio: 4},     // ¥6/24
		{MaxTokens: 256000, InputRatio: 0.010 * RMB, CompletionRatio: 4},     // ¥10/40
		{MaxTokens: 0, InputRatio: 0.020 * RMB, CompletionRatio: 10},         // ¥20/200
	},
	// qwen3-coder-flash: ≤32K / ≤128K / ≤256K / ≤1M
	"qwen3-coder-flash": {
		{MaxTokens: 32000, InputRatio: 0.001 * RMB, CompletionRatio: 4},      // ¥1/4
		{MaxTokens: 128000, InputRatio: 0.0015 * RMB, CompletionRatio: 4},    // ¥1.5/6
		{MaxTokens: 256000, InputRatio: 0.0025 * RMB, CompletionRatio: 4},    // ¥2.5/10
		{MaxTokens: 0, InputRatio: 0.005 * RMB, CompletionRatio: 5},          // ¥5/25
	},
}

// GetTieredPricing 根据模型名和输入 token 数量返回阶梯计价的倍率
// 返回: inputRatio, completionRatio, found
func GetTieredPricing(modelName string, promptTokens int) (float64, float64, bool) {
	tiers, ok := defaultTieredPricing[modelName]
	if !ok {
		return 0, 0, false
	}
	for _, tier := range tiers {
		if tier.MaxTokens == 0 || promptTokens <= tier.MaxTokens {
			return tier.InputRatio, tier.CompletionRatio, true
		}
	}
	// 如果都不匹配，使用最后一个阶梯
	last := tiers[len(tiers)-1]
	return last.InputRatio, last.CompletionRatio, true
}

// GetTieredPricingTiers 返回模型的所有阶梯配置（用于前端展示）
func GetTieredPricingTiers(modelName string) []PriceTier {
	tiers, ok := defaultTieredPricing[modelName]
	if !ok {
		return nil
	}
	return tiers
}

// InitRatioSettings initializes all model related settings maps
func InitRatioSettings() {
	// Initialize modelPriceMap
	modelPriceMapMutex.Lock()
	modelPriceMap = defaultModelPrice
	modelPriceMapMutex.Unlock()

	// Initialize modelRatioMap
	modelRatioMapMutex.Lock()
	modelRatioMap = defaultModelRatio
	modelRatioMapMutex.Unlock()

	// Initialize CompletionRatio
	CompletionRatioMutex.Lock()
	CompletionRatio = defaultCompletionRatio
	CompletionRatioMutex.Unlock()

	// Initialize cacheRatioMap
	cacheRatioMapMutex.Lock()
	cacheRatioMap = defaultCacheRatio
	cacheRatioMapMutex.Unlock()

	// initialize imageRatioMap
	imageRatioMapMutex.Lock()
	imageRatioMap = defaultImageRatio
	imageRatioMapMutex.Unlock()

	// initialize audioRatioMap
	audioRatioMapMutex.Lock()
	audioRatioMap = defaultAudioRatio
	audioRatioMapMutex.Unlock()

	// initialize audioCompletionRatioMap
	audioCompletionRatioMapMutex.Lock()
	audioCompletionRatioMap = defaultAudioCompletionRatio
	audioCompletionRatioMapMutex.Unlock()
}

func GetModelPriceMap() map[string]float64 {
	modelPriceMapMutex.RLock()
	defer modelPriceMapMutex.RUnlock()
	return modelPriceMap
}

func ModelPrice2JSONString() string {
	modelPriceMapMutex.RLock()
	defer modelPriceMapMutex.RUnlock()

	jsonBytes, err := common.Marshal(modelPriceMap)
	if err != nil {
		common.SysError("error marshalling model price: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateModelPriceByJSONString(jsonStr string) error {
	modelPriceMapMutex.Lock()
	defer modelPriceMapMutex.Unlock()
	modelPriceMap = make(map[string]float64)
	err := json.Unmarshal([]byte(jsonStr), &modelPriceMap)
	if err == nil {
		InvalidateExposedDataCache()
	}
	return err
}

// GetModelPrice 返回模型的价格，如果模型不存在则返回-1，false
func GetModelPrice(name string, printErr bool) (float64, bool) {
	modelPriceMapMutex.RLock()
	defer modelPriceMapMutex.RUnlock()

	name = FormatMatchingModelName(name)

	if strings.HasSuffix(name, CompactModelSuffix) {
		price, ok := modelPriceMap[CompactWildcardModelKey]
		if !ok {
			if printErr {
				common.SysError("model price not found: " + name)
			}
			return -1, false
		}
		return price, true
	}

	price, ok := modelPriceMap[name]
	if !ok {
		if printErr {
			common.SysError("model price not found: " + name)
		}
		return -1, false
	}
	return price, true
}

func UpdateModelRatioByJSONString(jsonStr string) error {
	modelRatioMapMutex.Lock()
	defer modelRatioMapMutex.Unlock()
	modelRatioMap = make(map[string]float64)
	err := common.Unmarshal([]byte(jsonStr), &modelRatioMap)
	if err == nil {
		InvalidateExposedDataCache()
	}
	return err
}

// 处理带有思考预算的模型名称，方便统一定价
func handleThinkingBudgetModel(name, prefix, wildcard string) string {
	if strings.HasPrefix(name, prefix) && strings.Contains(name, "-thinking-") {
		return wildcard
	}
	return name
}

func GetModelRatio(name string) (float64, bool, string) {
	modelRatioMapMutex.RLock()
	defer modelRatioMapMutex.RUnlock()

	name = FormatMatchingModelName(name)

	// 1. 精确匹配优先
	ratio, ok := modelRatioMap[name]
	if ok {
		return ratio, true, name
	}

	// 2. 通配符匹配 (支持 xxx* 格式)
	if wildcardRatio, matched := matchWildcardRatio(name, modelRatioMap); matched {
		return wildcardRatio, true, name
	}

	// 3. Compact 模型后缀匹配
	if strings.HasSuffix(name, CompactModelSuffix) {
		if wildcardRatio, ok := modelRatioMap[CompactWildcardModelKey]; ok {
			return wildcardRatio, true, name
		}
		return 0, true, name
	}

	// 4. 默认倍率
	return 37.5, operation_setting.SelfUseModeEnabled, name
}

// matchWildcardRatio 匹配通配符倍率配置
// 支持格式: "gpt-5*" 匹配 "gpt-5", "gpt-5.2", "gpt-5-mini" 等
// 返回: 倍率, 是否匹配
func matchWildcardRatio(name string, ratioMap map[string]float64) (float64, bool) {
	var bestMatch string
	var bestRatio float64

	for pattern, ratio := range ratioMap {
		if !strings.HasSuffix(pattern, "*") {
			continue
		}
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.HasPrefix(name, prefix) {
			// 选择最长匹配的前缀 (更精确的匹配优先)
			if len(prefix) > len(bestMatch) {
				bestMatch = prefix
				bestRatio = ratio
			}
		}
	}

	if bestMatch != "" {
		return bestRatio, true
	}
	return 0, false
}

func DefaultModelRatio2JSONString() string {
	jsonBytes, err := common.Marshal(defaultModelRatio)
	if err != nil {
		common.SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func GetDefaultModelRatioMap() map[string]float64 {
	return defaultModelRatio
}

func GetDefaultModelPriceMap() map[string]float64 {
	return defaultModelPrice
}

func GetDefaultImageRatioMap() map[string]float64 {
	return defaultImageRatio
}

func GetDefaultAudioRatioMap() map[string]float64 {
	return defaultAudioRatio
}

func GetDefaultAudioCompletionRatioMap() map[string]float64 {
	return defaultAudioCompletionRatio
}

func GetCompletionRatioMap() map[string]float64 {
	CompletionRatioMutex.RLock()
	defer CompletionRatioMutex.RUnlock()
	return CompletionRatio
}

func CompletionRatio2JSONString() string {
	CompletionRatioMutex.RLock()
	defer CompletionRatioMutex.RUnlock()

	jsonBytes, err := json.Marshal(CompletionRatio)
	if err != nil {
		common.SysError("error marshalling completion ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateCompletionRatioByJSONString(jsonStr string) error {
	CompletionRatioMutex.Lock()
	defer CompletionRatioMutex.Unlock()
	CompletionRatio = make(map[string]float64)
	err := common.Unmarshal([]byte(jsonStr), &CompletionRatio)
	if err == nil {
		InvalidateExposedDataCache()
	}
	return err
}

func GetCompletionRatio(name string) float64 {
	CompletionRatioMutex.RLock()
	defer CompletionRatioMutex.RUnlock()

	name = FormatMatchingModelName(name)

	// 1. 精确匹配（含/的模型名优先）
	if strings.Contains(name, "/") {
		if ratio, ok := CompletionRatio[name]; ok {
			return ratio
		}
	}

	// 2. 精确匹配
	if ratio, ok := CompletionRatio[name]; ok {
		return ratio
	}

	// 3. 通配符匹配
	if wildcardRatio, matched := matchWildcardRatio(name, CompletionRatio); matched {
		return wildcardRatio
	}

	// 4. 硬编码默认值
	hardCodedRatio, _ := getHardcodedCompletionModelRatio(name)
	return hardCodedRatio
}

func getHardcodedCompletionModelRatio(name string) (float64, bool) {

	isReservedModel := strings.HasSuffix(name, "-all") || strings.HasSuffix(name, "-gizmo-*")
	if isReservedModel {
		return 2, false
	}

	if strings.HasPrefix(name, "gpt-") {
		if strings.HasPrefix(name, "gpt-4o") {
			if name == "gpt-4o-2024-05-13" {
				return 3, true
			}
			if strings.HasPrefix(name, "gpt-4o-mini-tts") {
				return 20, false
			}
			return 4, false
		}
		// gpt-5 匹配
		if strings.HasPrefix(name, "gpt-5") {
			return 8, true
		}
		// gpt-4.5-preview匹配
		if strings.HasPrefix(name, "gpt-4.5-preview") {
			return 2, true
		}
		if strings.HasPrefix(name, "gpt-4-turbo") || strings.HasSuffix(name, "gpt-4-1106") || strings.HasSuffix(name, "gpt-4-1105") {
			return 3, true
		}
		// 没有特殊标记的 gpt-4 模型默认倍率为 2
		return 2, false
	}
	if strings.HasPrefix(name, "o1") || strings.HasPrefix(name, "o3") {
		return 4, true
	}
	if name == "chatgpt-4o-latest" {
		return 3, true
	}

	if strings.Contains(name, "claude-3") {
		return 5, true
	} else if strings.Contains(name, "claude-sonnet-4") || strings.Contains(name, "claude-opus-4") || strings.Contains(name, "claude-haiku-4") {
		return 5, true
	} else if strings.Contains(name, "claude-instant-1") || strings.Contains(name, "claude-2") {
		return 3, true
	}

	if strings.HasPrefix(name, "gpt-3.5") {
		if name == "gpt-3.5-turbo" || strings.HasSuffix(name, "0125") {
			// https://openai.com/blog/new-embedding-models-and-api-updates
			// Updated GPT-3.5 Turbo model and lower pricing
			return 3, true
		}
		if strings.HasSuffix(name, "1106") {
			return 2, true
		}
		return 4.0 / 3.0, true
	}
	if strings.HasPrefix(name, "mistral-") {
		return 3, true
	}
	if strings.HasPrefix(name, "gemini-") {
		if strings.HasPrefix(name, "gemini-1.5") {
			return 4, true
		} else if strings.HasPrefix(name, "gemini-2.0") {
			return 4, true
		} else if strings.HasPrefix(name, "gemini-2.5-pro") { // 移除preview来增加兼容性，这里假设正式版的倍率和preview一致
			return 8, false
		} else if strings.HasPrefix(name, "gemini-2.5-flash") { // 处理不同的flash模型倍率
			if strings.HasPrefix(name, "gemini-2.5-flash-preview") {
				if strings.HasSuffix(name, "-nothinking") {
					return 4, false
				}
				return 3.5 / 0.15, false
			}
			if strings.HasPrefix(name, "gemini-2.5-flash-lite") {
				return 4, false
			}
			return 2.5 / 0.3, false
		} else if strings.HasPrefix(name, "gemini-robotics-er-1.5") {
			return 2.5 / 0.3, false
		} else if strings.HasPrefix(name, "gemini-3-pro") {
			if strings.HasPrefix(name, "gemini-3-pro-image") {
				return 60, false
			}
			return 6, false
		}
		return 4, false
	}
	if strings.HasPrefix(name, "command") {
		switch name {
		case "command-r":
			return 3, true
		case "command-r-plus":
			return 5, true
		case "command-r-08-2024":
			return 4, true
		case "command-r-plus-08-2024":
			return 4, true
		default:
			return 4, false
		}
	}
	// hint 只给官方上4倍率，由于开源模型供应商自行定价，不对其进行补全倍率进行强制对齐
	if strings.HasPrefix(name, "ERNIE-Speed-") {
		return 2, true
	} else if strings.HasPrefix(name, "ERNIE-Lite-") {
		return 2, true
	} else if strings.HasPrefix(name, "ERNIE-Character") {
		return 2, true
	} else if strings.HasPrefix(name, "ERNIE-Functions") {
		return 2, true
	}
	switch name {
	case "llama2-70b-4096":
		return 0.8 / 0.64, true
	case "llama3-8b-8192":
		return 2, true
	case "llama3-70b-8192":
		return 0.79 / 0.59, true
	}
	return 1, false
}

func GetAudioRatio(name string) float64 {
	audioRatioMapMutex.RLock()
	defer audioRatioMapMutex.RUnlock()
	name = FormatMatchingModelName(name)
	if ratio, ok := audioRatioMap[name]; ok {
		return ratio
	}
	return 1
}

func GetAudioCompletionRatio(name string) float64 {
	audioCompletionRatioMapMutex.RLock()
	defer audioCompletionRatioMapMutex.RUnlock()
	name = FormatMatchingModelName(name)
	if ratio, ok := audioCompletionRatioMap[name]; ok {

		return ratio
	}
	return 1
}

func ContainsAudioRatio(name string) bool {
	audioRatioMapMutex.RLock()
	defer audioRatioMapMutex.RUnlock()
	name = FormatMatchingModelName(name)
	_, ok := audioRatioMap[name]
	return ok
}

func ContainsAudioCompletionRatio(name string) bool {
	audioCompletionRatioMapMutex.RLock()
	defer audioCompletionRatioMapMutex.RUnlock()
	name = FormatMatchingModelName(name)
	_, ok := audioCompletionRatioMap[name]
	return ok
}

func ModelRatio2JSONString() string {
	modelRatioMapMutex.RLock()
	defer modelRatioMapMutex.RUnlock()

	jsonBytes, err := common.Marshal(modelRatioMap)
	if err != nil {
		common.SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

var defaultImageRatio = map[string]float64{
	"gpt-image-1": 2,
}
var imageRatioMap map[string]float64
var imageRatioMapMutex sync.RWMutex
var (
	audioRatioMap      map[string]float64 = nil
	audioRatioMapMutex                    = sync.RWMutex{}
)
var (
	audioCompletionRatioMap      map[string]float64 = nil
	audioCompletionRatioMapMutex                    = sync.RWMutex{}
)

func ImageRatio2JSONString() string {
	imageRatioMapMutex.RLock()
	defer imageRatioMapMutex.RUnlock()
	jsonBytes, err := common.Marshal(imageRatioMap)
	if err != nil {
		common.SysError("error marshalling cache ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateImageRatioByJSONString(jsonStr string) error {
	imageRatioMapMutex.Lock()
	defer imageRatioMapMutex.Unlock()
	imageRatioMap = make(map[string]float64)
	return common.Unmarshal([]byte(jsonStr), &imageRatioMap)
}

func GetImageRatio(name string) (float64, bool) {
	imageRatioMapMutex.RLock()
	defer imageRatioMapMutex.RUnlock()
	ratio, ok := imageRatioMap[name]
	if !ok {
		return 1, false // Default to 1 if not found
	}
	return ratio, true
}

func AudioRatio2JSONString() string {
	audioRatioMapMutex.RLock()
	defer audioRatioMapMutex.RUnlock()
	jsonBytes, err := common.Marshal(audioRatioMap)
	if err != nil {
		common.SysError("error marshalling audio ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateAudioRatioByJSONString(jsonStr string) error {

	tmp := make(map[string]float64)
	if err := common.Unmarshal([]byte(jsonStr), &tmp); err != nil {
		return err
	}
	audioRatioMapMutex.Lock()
	audioRatioMap = tmp
	audioRatioMapMutex.Unlock()
	InvalidateExposedDataCache()
	return nil
}

func AudioCompletionRatio2JSONString() string {
	audioCompletionRatioMapMutex.RLock()
	defer audioCompletionRatioMapMutex.RUnlock()
	jsonBytes, err := common.Marshal(audioCompletionRatioMap)
	if err != nil {
		common.SysError("error marshalling audio completion ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateAudioCompletionRatioByJSONString(jsonStr string) error {
	tmp := make(map[string]float64)
	if err := common.Unmarshal([]byte(jsonStr), &tmp); err != nil {
		return err
	}
	audioCompletionRatioMapMutex.Lock()
	audioCompletionRatioMap = tmp
	audioCompletionRatioMapMutex.Unlock()
	InvalidateExposedDataCache()
	return nil
}

func GetModelRatioCopy() map[string]float64 {
	modelRatioMapMutex.RLock()
	defer modelRatioMapMutex.RUnlock()
	copyMap := make(map[string]float64, len(modelRatioMap))
	for k, v := range modelRatioMap {
		copyMap[k] = v
	}
	return copyMap
}

func GetModelPriceCopy() map[string]float64 {
	modelPriceMapMutex.RLock()
	defer modelPriceMapMutex.RUnlock()
	copyMap := make(map[string]float64, len(modelPriceMap))
	for k, v := range modelPriceMap {
		copyMap[k] = v
	}
	return copyMap
}

func GetCompletionRatioCopy() map[string]float64 {
	CompletionRatioMutex.RLock()
	defer CompletionRatioMutex.RUnlock()
	copyMap := make(map[string]float64, len(CompletionRatio))
	for k, v := range CompletionRatio {
		copyMap[k] = v
	}
	return copyMap
}

// 转换模型名，减少渠道必须配置各种带参数模型
func FormatMatchingModelName(name string) string {

	if strings.HasPrefix(name, "gemini-2.5-flash-lite") {
		name = handleThinkingBudgetModel(name, "gemini-2.5-flash-lite", "gemini-2.5-flash-lite-thinking-*")
	} else if strings.HasPrefix(name, "gemini-2.5-flash") {
		name = handleThinkingBudgetModel(name, "gemini-2.5-flash", "gemini-2.5-flash-thinking-*")
	} else if strings.HasPrefix(name, "gemini-2.5-pro") {
		name = handleThinkingBudgetModel(name, "gemini-2.5-pro", "gemini-2.5-pro-thinking-*")
	}

	if strings.HasPrefix(name, "gpt-4-gizmo") {
		name = "gpt-4-gizmo-*"
	}
	if strings.HasPrefix(name, "gpt-4o-gizmo") {
		name = "gpt-4o-gizmo-*"
	}
	return name
}

// result: 倍率or价格， usePrice， exist
func GetModelRatioOrPrice(model string) (float64, bool, bool) { // price or ratio
	price, usePrice := GetModelPrice(model, false)
	if usePrice {
		return price, true, true
	}
	modelRatio, success, _ := GetModelRatio(model)
	if success {
		return modelRatio, false, true
	}
	return 37.5, false, false
}
