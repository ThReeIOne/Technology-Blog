package upload

import "io"

type OSS interface {
	UploadFile(file io.Reader, objectKey string) (string, error)
	UploadBytes(file []byte, objectKey string) (string, error)
	DeleteFile(key string) error
}

func NewOss() OSS {
	return &TencentCOS{}
}
