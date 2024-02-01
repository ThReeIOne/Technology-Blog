package httprequest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func Post(url string, jsonPayload any, headers map[string]string) (responsePayload []byte, statusCode int, err error) {
	return PostWithClient(GetDefaultClient(), url, jsonPayload, headers)
}

func PostWithClient(
	client *http.Client,
	url string,
	jsonPayload any,
	headers map[string]string,
) (
	responsePayload []byte,
	statusCode int,
	err error,
) {
	bs, err := json.Marshal(jsonPayload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	statusCode = response.StatusCode
	responsePayload, err = io.ReadAll(response.Body)
	return
}
