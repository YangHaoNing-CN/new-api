package ali

var ModelList = []string{
	// ========== 旗舰模型 ==========
	"qwen3-max",         // 通义千问3旗舰版，效果最好，适合复杂多步骤任务，阶梯计价
	"qwen3-max-preview", // 通义千问3旗舰版预览，支持思考模式，阶梯计价
	"qwen-max",          // 通义千问旗舰版，仅非思考模式，无阶梯计价
	"qwen-max-latest",   // 通义千问旗舰版最新快照
	// ========== 中端模型 ==========
	"qwen-plus",        // 通义千问Plus，能力均衡，推理/成本/速度适中，阶梯计价
	"qwen-plus-latest", // 通义千问Plus最新快照
	// ========== 轻量模型 ==========
	"qwen-flash",        // 通义千问Flash，速度最快成本极低，阶梯计价
	"qwen-turbo",        // 通义千问Turbo，建议替换为Flash，无阶梯计价
	"qwen-turbo-latest", // 通义千问Turbo最新快照
	// ========== 长文本 & 深度研究 ==========
	"qwen-long",          // 通义千问长文本模型，无阶梯计价
	"qwen-long-latest",   // 通义千问长文本最新快照
	"qwen-deep-research", // 通义千问深度研究，适合复杂分析任��
	// ========== Qwen3 开源系列 ==========
	"qwen3-235b-a22b", // Qwen3 MoE 235B (激活22B)，支持思考/非思考模式
	"qwen3-32b",       // Qwen3 32B Dense，支持思考/非思考模式
	"qwen3-30b-a3b",   // Qwen3 MoE 30B (激活3B)，轻量级
	"qwen3-14b",       // Qwen3 14B Dense
	"qwen3-8b",        // Qwen3 8B Dense
	"qwen3-4b",        // Qwen3 4B Dense
	"qwen3-1.7b",      // Qwen3 1.7B Dense
	"qwen3-0.6b",      // Qwen3 0.6B Dense，最轻量
	"qwq-32b",         // QwQ 32B，专注推理和思考
	// ========== Vision 主力模型 ==========
	"qwen3-vl-plus",     // 通义千问3视觉Plus，图文理解，阶梯计价
	"qwen3-vl-flash",    // 通义千问3视觉Flash，快速图文理解，阶梯计价
	"qwen-vl-max",       // 通义千问视觉旗舰版，无阶梯计价
	"qwen-vl-max-latest", // 通义千问视觉旗舰版最新快照
	"qwen-vl-plus",       // 通义千问视觉Plus，无阶梯计价
	"qwen-vl-plus-latest", // 通义千问视觉Plus最新快照
	"qwen-vl-ocr",         // 通义千问OCR，文字识别专用
	"qwen-vl-ocr-latest",  // 通义千问OCR最新快照
	// ========== Vision 开源系列 ==========
	"qwen3-vl-32b-instruct",    // Qwen3 VL 32B，视觉语言模型
	"qwen3-vl-8b-instruct",     // Qwen3 VL 8B，轻量视觉语言模型
	"qwen2.5-vl-72b-instruct",  // Qwen2.5 VL 72B，大型视觉语言模型
	"qwen2.5-vl-32b-instruct",  // Qwen2.5 VL 32B
	"qwen2.5-vl-7b-instruct",   // Qwen2.5 VL 7B
	// ========== Coder ==========
	"qwen3-coder-plus",    // 通义千问3代码Plus，代码生成与理解，阶梯计价
	"qwen3-coder-flash",   // 通义千问3代码Flash，快速代码生成，阶梯计价
	"qwen-coder-plus",     // 通义千问代码Plus，无阶梯计价
	"qwen-coder-plus-latest", // 通义千问代码Plus最新快照
	"qwen-coder-turbo",       // 通义千问代码Turbo，无阶梯计价
	"qwen-coder-turbo-latest", // 通义千问代码Turbo最新快照
	// ========== Math ==========
	"qwen-math-plus",  // 通义千问数学Plus，数学推理专用
	"qwen-math-turbo", // 通义千问数学Turbo
	// ========== Qwen2.5 开源系列 ==========
	"qwen2.5-72b-instruct", // Qwen2.5 72B，大型通用模型
	"qwen2.5-32b-instruct", // Qwen2.5 32B
	"qwen2.5-14b-instruct", // Qwen2.5 14B
	"qwen2.5-7b-instruct",  // Qwen2.5 7B
	// ========== Embedding模型 ==========
	"text-embedding-v1", // 通用文本向量模型
	// ========== Rerank模型 ==========
	"gte-rerank-v2", // GTE Rerank V2，文本重排序
	// ========== 文生图模型 (Text to Image) ==========
	"wanx2.1-t2i-turbo", // 万相2.1极速版，文生图
	"wanx2.1-t2i-plus",  // 万相2.1专业版，文生图
	"wanx-v1",           // 通义万相-文本生成图像
	"flux-schnell",      // FLUX极速版，文生图
	"flux-dev",          // FLUX开发版，文生图
	"flux-merged",       // FLUX融合版，文生图
	// ========== 图片编辑模型 (Image Edit) ==========
	"wanx-style-repaint-v1",   // 涂鸦作画/风格重绘
	"wanx-sketch-to-image-v1", // 线稿生成图像
}

var ChannelName = "ali"
