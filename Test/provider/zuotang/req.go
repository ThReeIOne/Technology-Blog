package zuotang

import (
	"Technology-Blog/Test/config"
	"bytes"

	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type VisualSegmentationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TaskId      string  `json:"task_id"`
		Image       string  `json:"image"`
		Mask        string  `json:"mask"`
		MaskObj     string  `json:"mask_obj"`
		ImageObj    string  `json:"image_obj"`
		ReturnType  uint    `json:"return_type"`
		OutputType  uint    `json:"output_type"`
		Type        string  `json:"type"`
		ResultType  string  `json:"result_type"`
		TimeElapsed float64 `json:"time_elapsed"`
		Progress    uint    `json:"progress"`
		State       int     `json:"state"`
		ImageWidth  uint    `json:"image_width"`
		ImageHeight uint    `json:"image_height"`
	} `json:"data"`
}

func Cutout(imageUrl string, tpe string) (string, error) {
	url := "https://techsz.aoscdn.com/api/tasks/visual/segmentation"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("sync", "1")
	_ = writer.WriteField("image_url", imageUrl)
	_ = writer.WriteField("type", tpe)
	_ = writer.WriteField("return_type", "1")
	_ = writer.WriteField("output_type", "2")
	_ = writer.WriteField("crop", "1")
	_ = writer.WriteField("format", "png")
	err := writer.Close()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Get(config.ZuoTangApiKey))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var taskResp *VisualSegmentationResponse
	err = json.Unmarshal(respBody, &taskResp)
	if err != nil {
		return "", err
	}

	if taskResp.Status != http.StatusOK {
		return "", fmt.Errorf(string(respBody))
	} else {
		if taskResp.Data.State == 1 {
			return taskResp.Data.Image, nil
		} else {
			return "", fmt.Errorf(string(respBody))
		}
	}
}
