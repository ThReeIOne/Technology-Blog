package httprequest

import (
	"fmt"
	"io"
	"net/http"
)

func Get(url string, queryParams map[string]interface{}, headers map[string]string) (responsePayload []byte, statusCode int, err error) {
	return GetWithClient(GetDefaultClient(), url, queryParams, headers)
}

func GetWithClient(
	client *http.Client,
	url string,
	queryParams map[string]interface{},
	headers map[string]string,
) (
	responsePayload []byte,
	statusCode int,
	err error,
) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	query := req.URL.Query()
	for k, v := range queryParams {
		query.Add(k, fmt.Sprint(v))
	}
	req.URL.RawQuery = query.Encode()

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
