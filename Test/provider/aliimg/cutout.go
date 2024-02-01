package aliimg

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/lib"
	"Technology-Blog/Test/provider/httprequest"
	"bytes"

	"fmt"
	"io"
	"net/http"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	imageseg "github.com/alibabacloud-go/imageseg-20191230/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

var client *imageseg.Client

func getClient(ak *string, sk *string) (*imageseg.Client, error) {
	if client == nil {
		host := config.Get(config.AliImageSegHost)
		cfg := &openapi.Config{
			AccessKeyId:     ak,
			AccessKeySecret: sk,
			Endpoint:        &host,
		}
		// Endpoint 请参考 https://api.aliyun.com/product/imageseg
		newClient, err := imageseg.NewClient(cfg)
		if err != nil {
			return nil, err
		}
		client = newClient
	}
	return client, nil
}

func imageToFile(imageUrl string) (io.Reader, error) {
	resp, code, err := httprequest.Get(imageUrl, map[string]interface{}{}, map[string]string{})
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf(string(resp))
	}
	return bytes.NewReader(resp), nil
}

func Cutout(imageUrl string) (string, error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html
	ak := config.Get(config.AliAk)
	sk := config.Get(config.AliSk)
	client, err := getClient(&ak, &sk)
	if err != nil {
		return "", err
	}

	form := "crop"
	imageObj, err := imageToFile(imageUrl)
	if err != nil {
		return "", err
	}
	if bs, err := io.ReadAll(imageObj); err == nil {
		imageObj = bytes.NewReader(lib.ImgCompress(bs))
	}
	segmentCommodityRequest := &imageseg.SegmentCommodityAdvanceRequest{
		ImageURLObject: imageObj,
		ReturnForm:     &form,
	}
	runtime := &util.RuntimeOptions{}
	resp, err := client.SegmentCommodityAdvance(segmentCommodityRequest, runtime)
	if err != nil {
		return "", err
	}
	errorRespPtr := util.ToJSONString(resp)
	if resp != nil && resp.Body != nil && resp.Body.Data != nil && resp.Body.Data.ImageURL != nil {
		return *resp.Body.Data.ImageURL, nil
	}

	if errorRespPtr != nil {
		return "", fmt.Errorf(*errorRespPtr)
	}
	return "", fmt.Errorf("failed to request ali cloud image segmentation api")
}
