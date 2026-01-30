package ali

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"

	"github.com/gin-gonic/gin"
)

func oaiFormEdit2WanxImageEdit(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (*AliImageRequest, error) {
	var err error
	var imageRequest AliImageRequest
	imageRequest.Model = request.Model
	imageRequest.ResponseFormat = request.ResponseFormat
	wanInput := WanImageInput{
		Prompt: request.Prompt,
	}

	if err := common.UnmarshalBodyReusable(c, &wanInput); err != nil {
		return nil, err
	}
	if wanInput.Images, err = getImageBase64sFromForm(c, "image"); err != nil {
		return nil, fmt.Errorf("get image base64s from form failed: %w", err)
	}
	//wanParams := WanImageParameters{
	//	N: int(request.N),
	//}
	imageRequest.Input = wanInput
	imageRequest.Parameters = AliImageParameters{
		N: int(request.N),
	}
	info.PriceData.AddOtherRatio("n", float64(imageRequest.Parameters.N))

	return &imageRequest, nil
}

func isOldWanModel(modelName string) bool {
	return strings.Contains(modelName, "wan") && !strings.Contains(modelName, "wan2.6") && !isFluxModel(modelName)
}

func isWanModel(modelName string) bool {
	return strings.Contains(modelName, "wan") && !isFluxModel(modelName)
}

// isFluxModel 判断是否为 FLUX 系列模型
func isFluxModel(modelName string) bool {
	return strings.HasPrefix(modelName, "flux-")
}

// isWanT2IModel 判断是否为万相文生图模型
func isWanT2IModel(modelName string) bool {
	return strings.HasPrefix(modelName, "wanx2.1-t2i-") || modelName == "wanx-v1"
}

// isWanImageEditModel 判断是否为万相图片编辑模型
func isWanImageEditModel(modelName string) bool {
	return strings.Contains(modelName, "imageedit") ||
		modelName == "wanx-style-repaint-v1" ||
		modelName == "wanx-background-generation" ||
		modelName == "wanx-sketch-to-image-v1"
}
