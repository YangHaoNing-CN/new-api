package ali

var ModelList = []string{
	// ========== 聊天模型 ==========
	"qwen-turbo",
	"qwen-plus",
	"qwen-max",
	"qwen-max-longcontext",
	"qwq-32b",
	"qwen3-235b-a22b",
	// ========== Embedding模型 ==========
	"text-embedding-v1",
	// ========== Rerank模型 ==========
	"gte-rerank-v2",
	// ========== 文生图模型 (Text to Image) ==========
	"wanx2.1-t2i-turbo", // 万相2.1极速版
	"wanx2.1-t2i-plus",  // 万相2.1专业版
	"wanx-v1",           // 通义万相-文本生成图像
	"flux-schnell",      // FLUX极速版
	"flux-dev",          // FLUX开发版
	"flux-merged",       // FLUX融合版
	// ========== 图片编辑模型 (Image Edit) ==========
	"wanx2.1-imageedit-v1",       // 万相2.1图像编辑
	"wanx2.0-imageedit-plus",     // 万相2.0图像编辑
	"wanx-style-repaint-v1",      // 涂鸦作画/风格重绘
	"wanx-background-generation", // 图像背景生成
	"wanx-sketch-to-image-v1",    // 线稿生成图像
	"flux-inpaint",               // FLUX图像修复
	"flux-redux",                 // FLUX图像变体
}

var ChannelName = "ali"
