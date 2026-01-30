package ali

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"

	"github.com/gin-gonic/gin"
)

// FluxImageInput FLUX 文生图输入参数
type FluxImageInput struct {
	Prompt string `json:"prompt"`
}

// FluxImageEditInput FLUX 图片编辑输入参数
type FluxImageEditInput struct {
	Prompt      string `json:"prompt,omitempty"`       // 文本提示词
	ImageURL    string `json:"image_url,omitempty"`    // 原始图片URL (flux-inpaint, flux-redux)
	MaskURL     string `json:"mask_url,omitempty"`     // 遮罩图片URL (flux-inpaint)
	RefImageURL string `json:"ref_image_url,omitempty"` // 参考图片URL (flux-redux)
}

// FluxImageParameters FLUX 图片生成参数
type FluxImageParameters struct {
	Size             string  `json:"size,omitempty"`             // 图片尺寸: "1024*1024", "512*1024" 等
	N                int     `json:"n,omitempty"`                // 生成数量 1-4
	Steps            int     `json:"steps,omitempty"`            // 推理步数 (flux-dev: 1-50)
	Guidance         float64 `json:"guidance,omitempty"`         // 引导强度 (flux-dev: 1.5-5.0)
	Seed             int     `json:"seed,omitempty"`             // 随机种子
	PromptUpsampling bool    `json:"prompt_upsampling,omitempty"` // 提示词优化
}

// FluxImageRequest FLUX 图片生成请求
type FluxImageRequest struct {
	Model      string              `json:"model"`
	Input      interface{}         `json:"input"`
	Parameters FluxImageParameters `json:"parameters,omitempty"`
}

// oaiImage2FluxImageRequest 将 OpenAI 格式图片请求转换为 FLUX 格式
func oaiImage2FluxImageRequest(info *relaycommon.RelayInfo, request dto.ImageRequest) (*FluxImageRequest, error) {
	fluxRequest := &FluxImageRequest{
		Model: request.Model,
		Input: FluxImageInput{
			Prompt: request.Prompt,
		},
	}

	// 设置参数
	params := FluxImageParameters{}

	// 处理尺寸
	if request.Size != "" {
		// 将 "1024x1024" 格式转换为 "1024*1024"
		params.Size = strings.Replace(request.Size, "x", "*", -1)
	}

	// 处理数量
	if request.N > 0 {
		params.N = int(request.N)
		info.PriceData.AddOtherRatio("n", float64(params.N))
	}

	// 处理 extra 参数
	if request.Extra != nil {
		if val, ok := request.Extra["parameters"]; ok {
			// 将 json.RawMessage 解析为 map
			var paramsMap map[string]interface{}
			if err := common.Unmarshal(val, &paramsMap); err == nil {
				if steps, ok := paramsMap["steps"].(float64); ok {
					params.Steps = int(steps)
				}
				if guidance, ok := paramsMap["guidance"].(float64); ok {
					params.Guidance = guidance
				}
				if seed, ok := paramsMap["seed"].(float64); ok {
					params.Seed = int(seed)
				}
				if promptUpsampling, ok := paramsMap["prompt_upsampling"].(bool); ok {
					params.PromptUpsampling = promptUpsampling
				}
				if size, ok := paramsMap["size"].(string); ok {
					params.Size = size
				}
				if n, ok := paramsMap["n"].(float64); ok {
					params.N = int(n)
					info.PriceData.AddOtherRatio("n", float64(params.N))
				}
			}
		}

		// 检查 input 字段
		if val, ok := request.Extra["input"]; ok {
			var inputMap map[string]interface{}
			if err := common.Unmarshal(val, &inputMap); err == nil {
				if prompt, ok := inputMap["prompt"].(string); ok && prompt != "" {
					fluxRequest.Input = FluxImageInput{Prompt: prompt}
				}
			}
		}
	}

	fluxRequest.Parameters = params
	return fluxRequest, nil
}

// oaiImage2FluxImageEditRequest 将 OpenAI 格式图片编辑请求转换为 FLUX 格式
func oaiImage2FluxImageEditRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (*FluxImageRequest, error) {
	fluxRequest := &FluxImageRequest{
		Model: request.Model,
	}

	// 根据模型类型构建不同的 input
	modelName := strings.ToLower(request.Model)

	if strings.Contains(modelName, "inpaint") {
		// flux-inpaint 需要 image_url 和 mask_url
		input := FluxImageEditInput{
			Prompt: request.Prompt,
		}

		// 尝试从 extra 获取图片URL
		if request.Extra != nil {
			if val, ok := request.Extra["input"]; ok {
				var inputMap map[string]interface{}
				if err := common.Unmarshal(val, &inputMap); err == nil {
					if imageURL, ok := inputMap["image_url"].(string); ok {
						input.ImageURL = imageURL
					}
					if maskURL, ok := inputMap["mask_url"].(string); ok {
						input.MaskURL = maskURL
					}
				}
			}
		}

		// 如果使用表单上传
		if strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
			imageBase64s, err := getImageBase64sFromForm(c, "image")
			if err == nil && len(imageBase64s) > 0 {
				input.ImageURL = imageBase64s[0]
			}
			// 尝试获取 mask
			maskBase64s, err := getImageBase64sFromForm(c, "mask")
			if err == nil && len(maskBase64s) > 0 {
				input.MaskURL = maskBase64s[0]
			}
		}

		fluxRequest.Input = input
	} else if strings.Contains(modelName, "redux") {
		// flux-redux 需要 ref_image_url
		input := FluxImageEditInput{
			Prompt: request.Prompt,
		}

		if request.Extra != nil {
			if val, ok := request.Extra["input"]; ok {
				var inputMap map[string]interface{}
				if err := common.Unmarshal(val, &inputMap); err == nil {
					if refImageURL, ok := inputMap["ref_image_url"].(string); ok {
						input.RefImageURL = refImageURL
					}
					if imageURL, ok := inputMap["image_url"].(string); ok {
						input.ImageURL = imageURL
					}
				}
			}
		}

		// 如果使用表单上传
		if strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
			imageBase64s, err := getImageBase64sFromForm(c, "image")
			if err == nil && len(imageBase64s) > 0 {
				input.RefImageURL = imageBase64s[0]
			}
		}

		fluxRequest.Input = input
	} else {
		// 其他 FLUX 模型
		fluxRequest.Input = FluxImageInput{
			Prompt: request.Prompt,
		}
	}

	// 设置参数
	params := FluxImageParameters{}
	if request.Size != "" {
		params.Size = strings.Replace(request.Size, "x", "*", -1)
	}
	if request.N > 0 {
		params.N = int(request.N)
		info.PriceData.AddOtherRatio("n", float64(params.N))
	}

	// 处理 extra 参数
	if request.Extra != nil {
		if val, ok := request.Extra["parameters"]; ok {
			var paramsMap map[string]interface{}
			if err := common.Unmarshal(val, &paramsMap); err == nil {
				if steps, ok := paramsMap["steps"].(float64); ok {
					params.Steps = int(steps)
				}
				if guidance, ok := paramsMap["guidance"].(float64); ok {
					params.Guidance = guidance
				}
				if seed, ok := paramsMap["seed"].(float64); ok {
					params.Seed = int(seed)
				}
				if size, ok := paramsMap["size"].(string); ok {
					params.Size = size
				}
			}
		}
	}

	fluxRequest.Parameters = params
	return fluxRequest, nil
}
