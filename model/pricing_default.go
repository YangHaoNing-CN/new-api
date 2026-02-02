package model

import (
	"strings"
)

// 默认模型描述，用于自动创建的模型元数据
// 价格数据来源：阿里云百炼官方文档（中国内地价格），标注为"官网价"
var defaultModelDescriptions = map[string]string{
	// ========== 通义千问Max ==========
	"qwen3-max":            "通义千问3旗舰版，效果最好，适合复杂多步骤任务。仅非思考模式，最大上下文252K。支持Batch调用(半价)和上下文缓存。官网价(阶梯)：≤32K ¥2.5/10, ≤128K ¥4/16, ≤252K ¥7/28 元/百万Token",
	"qwen3-max-2026-01-23": "通义千问3旗舰版 2026-01-23 快照。非思考和思考模式，最大上下文252K。官网价(���梯)：≤32K ¥2.5/10, ≤128K ¥4/16, ≤252K ¥7/28 元/百万Token",
	"qwen3-max-preview":    "通义千问3旗舰版预览。非思考和思考模式，最大上下文252K。支持上下文缓存。官网价(阶梯)：≤32K ¥6/24, ≤128K ¥10/40, ≤252K ¥15/60 元/百万Token",
	"qwen-max":             "通义千问旗舰版。仅非思考模式，无阶梯计价。支持Batch调用(半价)。官网价：输入¥2.4/输出¥9.6 元/百万Token",
	"qwen-max-latest":      "通义千问旗舰版最新快照。仅非思考模式。支持Batch调用(半价)。官网价：输入¥2.4/输出¥9.6 元/百万Token",
	// ========== 通义千问Plus ==========
	"qwen-plus":        "通义千问Plus，能力均衡，推理/成本/速度适中。支持Batch调用(半价)。官网价(阶梯)：≤128K ¥0.8/2, ≤256K ¥2.4/20, ≤1M ¥4.8/48 元/百万Token（非思考输出）",
	"qwen-plus-latest": "通义千问Plus最新快照。支持Batch调用(半价)。官网价(阶梯)：≤128K ¥0.8/2, ≤256K ¥2.4/20, ≤1M ¥4.8/48 元/百万Token",
	// ========== 通义千问Flash ==========
	"qwen-flash": "通义千问Flash，速度最快成本极低。非思考和思考模式。支持Batch调用(半价)和上下文缓存。官网价(阶梯)：≤128K ¥0.15/1.5, ≤256K ¥0.6/6, ≤1M ¥1.2/12 元/百万Token",
	// ========== 通义千问Turbo ==========
	"qwen-turbo":        "通义千问Turbo。非思考和思考模式，无阶梯计价。支持Batch调用(半价)。官网价：输入¥0.3/非思考输出¥0.6/思考输出¥3 元/百万Token",
	"qwen-turbo-latest": "通义千问Turbo最新快照。非思考和思考模式。支持Batch调用(半价)。官网价：输入¥0.3/非思考输出¥0.6/思考输出¥3 元/百万Token",
	// ========== 长文本 & 深度研究 ==========
	"qwen-long":          "通义千问长文本模型，无阶梯计价。支持Batch调用(半价)。官网价：输入¥0.5/输出¥2 元/百万Token",
	"qwen-long-latest":   "通义千问长文本最新快照。官网价：输入¥0.5/输出¥2 元/百万Token",
	"qwen-deep-research": "通义千问深度研究，适合复杂分析任务。仅支持中国内地。官网价：输入¥54/输出¥163 元/百万Token",
	// ========== Qwen3 开源系列 ==========
	"qwen3-235b-a22b": "Qwen3 MoE 235B (激活22B)，非思考和思考模式。官网价：输入¥2/非思考输出¥8/思考输出¥20 元/百万Token",
	"qwen3-32b":       "Qwen3 32B Dense，非思考和思考模式。官网价：输入¥2/非思考输出¥8/思考输出¥20 元/百万Token",
	"qwen3-30b-a3b":   "Qwen3 MoE 30B (激活3B)，非思考和思考模式。官网价：输入¥0.75/非思考输出¥3/思考输出¥7.5 元/百万Token",
	"qwen3-14b":       "Qwen3 14B Dense，非思考和思考模式。官网价：输入¥1/非思考输出¥4/思考输出¥10 元/百万Token",
	"qwen3-8b":        "Qwen3 8B Dense，非思考和思考模式。官网价：输入¥0.5/非思考输出¥2/思考输出¥5 元/百万Token",
	"qwen3-4b":        "Qwen3 4B Dense，非思考和思考模式。官网价：输入¥0.3/非思考输出¥1.2/思考输出¥3 元/百万Token",
	"qwen3-1.7b":      "Qwen3 1.7B Dense，非思考和思考模式。官网价：输入¥0.3/非思考输出¥1.2/思考输出¥3 元/百万Token",
	"qwen3-0.6b":      "Qwen3 0.6B Dense，最轻量，非思考和思考模式。官网价：输入¥0.3/非思考输出¥1.2/思考输出¥3 元/百万Token",
	"qwq-32b":         "QwQ 32B 开源版，专注推理和思考。仅支持中国内地。官网价：输入¥2/输出¥6 元/百万Token",
	// ========== Qwen3 thinking/instruct 变体 ==========
	"qwen3-235b-a22b-thinking-2507":  "Qwen3 235B 仅思考模式。官网价：输入¥2/输出¥20 元/百万Token",
	"qwen3-235b-a22b-instruct-2507":  "Qwen3 235B 仅非思考模式。官网价：输入¥2/输出¥8 元/百万Token",
	"qwen3-30b-a3b-thinking-2507":    "Qwen3 30B 仅思考模式。官网价：输入¥0.75/输出¥7.5 元/百万Token",
	"qwen3-30b-a3b-instruct-2507":    "Qwen3 30B 仅非思考模式。官网价：输入¥0.75/输出¥3 元/百万Token",
	"qwen3-next-80b-a3b-thinking":    "Qwen3 Next 80B 仅思考模式。官网价：输入¥1/输出¥10 元/百万Token",
	"qwen3-next-80b-a3b-instruct":    "Qwen3 Next 80B 仅非思考模式。官网价：输入¥1/输出¥4 元/百万Token",
	// ========== Vision 主力模型 ==========
	"qwen3-vl-plus":       "通义千问3视觉Plus，图文理解。非思考和思考模式。支持Batch调用(半价)和上下文缓存。官网价(阶梯)：≤32K ¥1/10, ≤128K ¥1.5/15, ≤256K ¥3/30 元/百万Token",
	"qwen3-vl-flash":      "通义千问3视觉Flash，快速图文理解。非思考和思考模式。支持Batch调用(半价)和上下文缓存。官网价(阶梯)：≤32K ¥0.15/1.5, ≤128K ¥0.3/3, ≤256K ¥0.6/6 元/百万Token",
	"qwen-vl-max":         "通义千问视觉旗舰版，无阶梯计价。支持Batch调用(半价)和上下文缓存。官网价：输入¥1.6/输出¥4 元/百万Token",
	"qwen-vl-max-latest":  "通义千问视觉旗舰版最新快照。支持Batch调用(半价)。官网价：输入¥1.6/输出¥4 元/百万Token",
	"qwen-vl-plus":        "通义千问视觉Plus，无阶梯计价。支持Batch调用(半价)和上下文缓存。官网价：输入¥0.8/输出¥2 元/百万Token",
	"qwen-vl-plus-latest": "通义千问视觉Plus最新快照。支持Batch调用(半价)。官网价：输入¥0.8/输出¥2 元/百万Token",
	"qwen-vl-ocr":         "通义千问OCR，文字识别专用。支持Batch调用(半价)。官网价：输入¥5/输出¥5 元/百万Token",
	"qwen-vl-ocr-latest":  "通义千问OCR最新快照。支持Batch调用(半价)。官网价：输入¥0.3/输出¥0.5 元/百万Token",
	// ========== Vision 开源系列 ==========
	"qwen3-vl-235b-a22b-thinking": "Qwen3 VL 235B 仅思考模式，视觉语言模型。官网价：输入¥2/输出¥20 元/百万Token",
	"qwen3-vl-235b-a22b-instruct": "Qwen3 VL 235B 仅非思考模式，视觉语言模型。官网价：输入¥2/输出¥8 元/百万Token",
	"qwen3-vl-32b-thinking":       "Qwen3 VL 32B 仅思考模式，视觉语言模型。官网价：输入¥2/输出¥20 元/百万Token",
	"qwen3-vl-32b-instruct":       "Qwen3 VL 32B 仅非思考模式，视觉语言模型。官网价：输入¥2/输出¥8 元/百万Token",
	"qwen3-vl-30b-a3b-thinking":   "Qwen3 VL 30B 仅思考模式，MoE视觉语言模型。官网价：输入¥0.75/输出¥7.5 元/百万Token",
	"qwen3-vl-30b-a3b-instruct":   "Qwen3 VL 30B 仅非思考模式，MoE视觉语言模型。官网价：输入¥0.75/输出¥3 元/百万Token",
	"qwen3-vl-8b-thinking":        "Qwen3 VL 8B 仅思考模式，轻量视觉语言模型。官网价：输入¥0.5/输出¥5 元/百万Token",
	"qwen3-vl-8b-instruct":        "Qwen3 VL 8B 仅非思考模式，轻量视觉语言模型。官网价：输入¥0.5/输出¥2 元/百万Token",
	"qwen2.5-vl-72b-instruct":     "Qwen2.5 VL 72B，大型视觉语言模型。官网价：输入¥16/输出¥48 元/百万Token",
	"qwen2.5-vl-32b-instruct":     "Qwen2.5 VL 32B 视觉语言模型。官网价：输入¥8/输出¥24 元/百万Token",
	"qwen2.5-vl-7b-instruct":      "Qwen2.5 VL 7B 视觉语言模型。官网价：输入¥2/输出¥5 元/百万Token",
	"qwen2.5-vl-3b-instruct":      "Qwen2.5 VL 3B 视觉语言模型。官网价：输入¥1.2/输出¥3.6 元/百万Token",
	// ========== Coder 主力模型 ==========
	"qwen3-coder-plus":  "通义千问3代码Plus，代码生成与理解。支持上下文缓存。官网价(阶梯)：≤32K ¥4/16, ≤128K ¥6/24, ≤256K ¥10/40, ≤1M ¥20/200 元/百万Token",
	"qwen3-coder-flash": "通义千问3代码Flash，快速代码生成。官网价(阶梯)：≤32K ¥1/4, ≤128K ¥1.5/6, ≤256K ¥2.5/10, ≤1M ¥5/25 元/百万Token",
	// ========== Coder 开源系列 ==========
	"qwen3-coder-480b-a35b-instruct": "Qwen3 Coder 480B 大型代码模型。官网价(阶梯)：≤32K ¥6/24, ≤128K ¥9/36, ≤200K ¥15/60 元/百万Token",
	"qwen3-coder-30b-a3b-instruct":   "Qwen3 Coder 30B 轻量代码模型。官网价(阶梯)：≤32K ¥1.5/6, ≤128K ¥2.25/9, ≤200K ¥3.75/15 元/百万Token",
	"qwen-coder-plus":                "通义千问代码Plus，无阶梯计价。官网价：输入¥3.5/输出¥7 元/百万Token",
	"qwen-coder-plus-latest":         "通义千问代码Plus最新快照。官网价：输入¥3.5/输出¥7 元/百万Token",
	"qwen-coder-turbo":               "通义千问代码Turbo，无阶梯计价。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen-coder-turbo-latest":        "通义千问代码Turbo最新快照。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen2.5-coder-32b-instruct":     "Qwen2.5 Coder 32B，无阶梯计价。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen2.5-coder-14b-instruct":     "Qwen2.5 Coder 14B，无阶梯计价。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen2.5-coder-7b-instruct":      "Qwen2.5 Coder 7B，无阶梯计价。官网价：输入¥1/输出¥2 元/百万Token",
	// ========== Math ==========
	"qwen-math-plus":            "通义千问数学Plus，数学推理专用。仅支持中国内地。官网价：输入¥4/输出¥12 元/百万Token",
	"qwen-math-turbo":           "通义千问数学Turbo。仅支持中国内地。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen2.5-math-72b-instruct": "Qwen2.5 Math 72B。仅支持中国内地。官网价：输入¥4/输出¥12 元/百万Token",
	"qwen2.5-math-7b-instruct":  "Qwen2.5 Math 7B。仅支持中国内地。官网价：输入¥1/输出¥2 元/百万Token",
	// ========== Qwen2.5 开源系列 ==========
	"qwen2.5-72b-instruct":    "Qwen2.5 72B，大型通用模型。官网价：输入¥4/输出¥12 元/百万Token",
	"qwen2.5-32b-instruct":    "Qwen2.5 32B 通用模型。官网价：输入¥2/输出¥6 元/百万Token",
	"qwen2.5-14b-instruct":    "Qwen2.5 14B 通用模型。官网价：输入¥1/输出¥3 元/百万Token",
	"qwen2.5-14b-instruct-1m": "Qwen2.5 14B 百万Token长文本版本。官网价：输入¥1/输出¥3 元/百万Token",
	"qwen2.5-7b-instruct":     "Qwen2.5 7B 通用模型。官网价：输入¥0.5/输出¥1 元/百万Token",
	"qwen2.5-7b-instruct-1m":  "Qwen2.5 7B 百万Token长文本版本。官网价：输入¥0.5/输出¥1 元/百万Token",
	"qwen2.5-3b-instruct":     "Qwen2.5 3B 轻量模型。官网价：输入¥0.3/输出¥0.9 元/百万Token",
	// ========== 图片/视频模型 ==========
	"wanx2.1-t2i-turbo":          "万相2.1极速版，文生图",
	"wanx2.1-t2i-plus":           "万相2.1专业版，文生图",
	"wanx-v1":                    "通义万相-文本生成图像",
	"wanx2.1-imageedit-v1":       "万相2.1图像编辑",
	"wanx2.0-imageedit-plus":     "万相2.0图像编辑",
	"wanx-style-repaint-v1":      "涂鸦作画/风格重绘",
	"wanx-background-generation": "图像背景生成",
	"wanx-sketch-to-image-v1":    "线稿生成图像",
	"wan2.6-i2v-flash":           "万相2.6图生视频极速版 (720P有声)",
	"wan2.6-i2v":                 "万相2.6图生视频 (720P有声)",
	"wan2.6-t2v":                 "万相2.6文生视频 (720P有声)",
	"flux-schnell":               "FLUX极速版，文生图",
	"flux-dev":                   "FLUX开发版，文生图",
	"flux-merged":                "FLUX融合版，文生图",
	"flux-inpaint":               "FLUX图像修复",
	"flux-redux":                 "FLUX图像变体",
	// ========== Embedding & Rerank ==========
	"text-embedding-v1": "通用文本向量模型",
	"gte-rerank-v2":     "GTE Rerank V2，文本重排序",
}

// 简化的供应商映射规则
var defaultVendorRules = map[string]string{
	"gpt":      "OpenAI",
	"dall-e":   "OpenAI",
	"whisper":  "OpenAI",
	"o1":       "OpenAI",
	"o3":       "OpenAI",
	"claude":   "Anthropic",
	"gemini":   "Google",
	"moonshot": "Moonshot",
	"kimi":     "Moonshot",
	"chatglm":  "智谱",
	"glm-":     "智谱",
	"qwen":     "阿里巴巴",
	"deepseek": "DeepSeek",
	"abab":     "MiniMax",
	"ernie":    "百度",
	"spark":    "讯飞",
	"hunyuan":  "腾讯",
	"command":  "Cohere",
	"@cf/":     "Cloudflare",
	"360":      "360",
	"yi":       "零一万物",
	"jina":     "Jina",
	"mistral":  "Mistral",
	"grok":     "xAI",
	"llama":    "Meta",
	"doubao":   "字节跳动",
	"kling":    "快手",
	"jimeng":   "即梦",
	"vidu":     "Vidu",
}

// 供应商默认图标映射
var defaultVendorIcons = map[string]string{
	"OpenAI":     "OpenAI",
	"Anthropic":  "Claude.Color",
	"Google":     "Gemini.Color",
	"Moonshot":   "Moonshot",
	"智谱":         "Zhipu.Color",
	"阿里巴巴":       "Qwen.Color",
	"DeepSeek":   "DeepSeek.Color",
	"MiniMax":    "Minimax.Color",
	"百度":         "Wenxin.Color",
	"讯飞":         "Spark.Color",
	"腾讯":         "Hunyuan.Color",
	"Cohere":     "Cohere.Color",
	"Cloudflare": "Cloudflare.Color",
	"360":        "Ai360.Color",
	"零一万物":       "Yi.Color",
	"Jina":       "Jina",
	"Mistral":    "Mistral.Color",
	"xAI":        "XAI",
	"Meta":       "Ollama",
	"字节跳动":       "Doubao.Color",
	"快手":         "Kling.Color",
	"即梦":         "Jimeng.Color",
	"Vidu":       "Vidu",
	"微软":         "AzureAI",
	"Microsoft":  "AzureAI",
	"Azure":      "AzureAI",
}

// initDefaultVendorMapping 简化的默认供应商映射
func initDefaultVendorMapping(metaMap map[string]*Model, vendorMap map[int]*Vendor, enableAbilities []AbilityWithChannel) {
	for _, ability := range enableAbilities {
		modelName := ability.Model
		if _, exists := metaMap[modelName]; exists {
			continue
		}

		// 匹配供应商
		vendorID := 0
		modelLower := strings.ToLower(modelName)
		for pattern, vendorName := range defaultVendorRules {
			if strings.Contains(modelLower, pattern) {
				vendorID = getOrCreateVendor(vendorName, vendorMap)
				break
			}
		}

		// 创建模型元数据
		desc := ""
		if d, ok := defaultModelDescriptions[modelName]; ok {
			desc = d
		}
		metaMap[modelName] = &Model{
			ModelName:   modelName,
			Description: desc,
			VendorID:    vendorID,
			Status:      1,
			NameRule:    NameRuleExact,
		}
	}
}

// 查找或创建供应商
func getOrCreateVendor(vendorName string, vendorMap map[int]*Vendor) int {
	// 查找现有供应商
	for id, vendor := range vendorMap {
		if vendor.Name == vendorName {
			return id
		}
	}

	// 创建新供应商
	newVendor := &Vendor{
		Name:   vendorName,
		Status: 1,
		Icon:   getDefaultVendorIcon(vendorName),
	}

	if err := newVendor.Insert(); err != nil {
		return 0
	}

	vendorMap[newVendor.Id] = newVendor
	return newVendor.Id
}

// 获取供应商默认图标
func getDefaultVendorIcon(vendorName string) string {
	if icon, exists := defaultVendorIcons[vendorName]; exists {
		return icon
	}
	return ""
}
