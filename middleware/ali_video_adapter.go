package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
)

// AliVideoInput 阿里云原生输入格式
type AliVideoInput struct {
	Prompt         string `json:"prompt,omitempty"`
	ImgURL         string `json:"img_url,omitempty"`
	FirstFrameURL  string `json:"first_frame_url,omitempty"`
	LastFrameURL   string `json:"last_frame_url,omitempty"`
	AudioURL       string `json:"audio_url,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	Template       string `json:"template,omitempty"`
}

// AliVideoParameters 阿里云原生参数格式
type AliVideoParameters struct {
	Resolution   string `json:"resolution,omitempty"`
	Size         string `json:"size,omitempty"`
	Duration     int    `json:"duration,omitempty"`
	PromptExtend *bool  `json:"prompt_extend,omitempty"`
	ShotType     string `json:"shot_type,omitempty"`
	Watermark    *bool  `json:"watermark,omitempty"`
	Audio        *bool  `json:"audio,omitempty"`
	Seed         int    `json:"seed,omitempty"`
}

// AliVideoRequest 阿里云原生请求格式
type AliVideoRequest struct {
	Model      string              `json:"model"`
	Input      *AliVideoInput      `json:"input,omitempty"`
	Parameters *AliVideoParameters `json:"parameters,omitempty"`
}

// AliVideoRequestConvert 转换阿里云原生请求格式为统一格式
func AliVideoRequestConvert() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 检查是否是查询请求
		if c.Request.Method == "GET" {
			// 从路径中提取 task_id
			path := c.Request.URL.Path
			if strings.Contains(path, "/tasks/") {
				parts := strings.Split(path, "/tasks/")
				if len(parts) > 1 {
					taskID := strings.Split(parts[1], "/")[0]
					c.Request.URL.Path = "/v1/video/generations/" + taskID
				}
			}
			c.Next()
			return
		}

		var aliReq AliVideoRequest
		if err := common.UnmarshalBodyReusable(c, &aliReq); err != nil {
			c.Next()
			return
		}

		// 构建统一请求格式
		unifiedReq := map[string]interface{}{
			"model": aliReq.Model,
		}

		// 处理 input
		if aliReq.Input != nil {
			if aliReq.Input.Prompt != "" {
				unifiedReq["prompt"] = aliReq.Input.Prompt
			}
			if aliReq.Input.ImgURL != "" {
				unifiedReq["input_reference"] = aliReq.Input.ImgURL
			}
		}

		// 处理 parameters 中的基本字段
		if aliReq.Parameters != nil {
			if aliReq.Parameters.Resolution != "" {
				unifiedReq["size"] = aliReq.Parameters.Resolution
			} else if aliReq.Parameters.Size != "" {
				unifiedReq["size"] = aliReq.Parameters.Size
			}
			if aliReq.Parameters.Duration > 0 {
				unifiedReq["duration"] = aliReq.Parameters.Duration
			}
		}

		// 将完整的原始请求作为 metadata 传递
		metadata := map[string]interface{}{}
		if aliReq.Input != nil {
			inputMap := map[string]interface{}{}
			if aliReq.Input.AudioURL != "" {
				inputMap["audio_url"] = aliReq.Input.AudioURL
			}
			if aliReq.Input.FirstFrameURL != "" {
				inputMap["first_frame_url"] = aliReq.Input.FirstFrameURL
			}
			if aliReq.Input.LastFrameURL != "" {
				inputMap["last_frame_url"] = aliReq.Input.LastFrameURL
			}
			if aliReq.Input.NegativePrompt != "" {
				inputMap["negative_prompt"] = aliReq.Input.NegativePrompt
			}
			if aliReq.Input.Template != "" {
				inputMap["template"] = aliReq.Input.Template
			}
			if len(inputMap) > 0 {
				metadata["input"] = inputMap
			}
		}

		if aliReq.Parameters != nil {
			paramsMap := map[string]interface{}{}
			if aliReq.Parameters.PromptExtend != nil {
				paramsMap["prompt_extend"] = *aliReq.Parameters.PromptExtend
			}
			if aliReq.Parameters.ShotType != "" {
				paramsMap["shot_type"] = aliReq.Parameters.ShotType
			}
			if aliReq.Parameters.Watermark != nil {
				paramsMap["watermark"] = *aliReq.Parameters.Watermark
			}
			if aliReq.Parameters.Audio != nil {
				paramsMap["audio"] = *aliReq.Parameters.Audio
			}
			if aliReq.Parameters.Seed > 0 {
				paramsMap["seed"] = aliReq.Parameters.Seed
			}
			if len(paramsMap) > 0 {
				metadata["parameters"] = paramsMap
			}
		}

		if len(metadata) > 0 {
			unifiedReq["metadata"] = metadata
		}

		jsonData, err := json.Marshal(unifiedReq)
		if err != nil {
			c.Next()
			return
		}

		// 重写请求体和路径
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
		c.Request.URL.Path = "/v1/video/generations"

		// 保存请求体供后续处理
		c.Set(common.KeyRequestBody, jsonData)
		c.Next()
	}
}
