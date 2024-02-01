package upload

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/lib"
	"bytes"

	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TencentCOS struct {
	client *cos.Client
}

func (t *TencentCOS) GetClient() *cos.Client {
	if t.client == nil {
		urlStr := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Get(config.TencentCosBucket), config.Get(config.TencentCosRegion))
		u, _ := url.Parse(urlStr)
		baseURL := &cos.BaseURL{BucketURL: u}
		t.client = cos.NewClient(baseURL, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.Get(config.TencentCosSecretId),
				SecretKey: config.Get(config.TencentCosSecretKey),
			},
		})
	}
	return t.client
}

func (t *TencentCOS) UploadFile(file io.Reader, objectKey string) (string, error) {
	objectKey = strings.TrimPrefix(objectKey, "/")
	_, err := t.GetClient().Object.Put(context.Background(), objectKey, file, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", lib.CdnPrefix(), objectKey), nil
}

func (t *TencentCOS) UploadBytes(file []byte, objectKey string) (string, error) {
	f := bytes.NewReader(file)
	return t.UploadFile(f, objectKey)
}

func (t *TencentCOS) DeleteFile(key string) error {
	name := config.Get(config.TencentCosBasePath) + "/" + key
	_, err := t.GetClient().Object.Delete(context.Background(), name)
	if err != nil {
		return err
	}
	return nil
}
