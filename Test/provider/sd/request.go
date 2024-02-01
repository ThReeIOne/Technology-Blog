package sd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	Options Options
}
type Options struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}
type Params struct {
	SceneId       string `json:"sence_id"`
	PositionX     int    `json:"position_x"`
	PositionY     int    `json:"position_y"`
	WbImg         string `json:"wb_img"`
	WhRatio       string `json:"whratio"`
	SceneImg      string `json:"sence_img"`
	Count         int    `json:"count"`
	ProductWidth  int    `json:"product_width"`
	ProductHeight int    `json:"product_height"`
}
type Response struct {
	SdResult string  `json:"sd_result"`
	Progress float64 `json:"progress"`
	Text     string  `json:"text"`
}

func New(opt Options) *Client {
	return &Client{
		Options: Options{
			Path:   opt.Path,
			Method: opt.Method,
		},
	}
}

func (c *Client) TextToImage(host string, params Params) (Response, error) {
	req, err := http.NewRequest(c.Options.Method, host+c.Options.Path, bytes.NewReader(body(params)))
	if err != nil {
		return Response{}, err
	}
	req.Header = headers(host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	var r Response
	err = json.Unmarshal(data, &r)
	return r, err
}

func (c *Client) GetProgress(host string) (string, error) {

	req, err := http.NewRequest(c.Options.Method, host+c.Options.Path, nil)
	if err != nil {
		return "0.00", err
	}

	req.Header = headers(host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "0.00", err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "0.00", err
	}
	defer res.Body.Close()

	var r Response
	err = json.Unmarshal(data, &r)
	return strconv.FormatFloat(r.Progress, 'f', 2, 64), err
}

func (c *Client) GetText(answerId string) (string, error) {
	return "", nil
}
func (c *Client) GetStatus(host string) (bool, error) {
	req, err := http.NewRequest(c.Options.Method, host+c.Options.Path, nil)
	if err != nil {
		return false, err
	}

	req.Header = headers(host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil

}
func headers(host string) http.Header {
	h := make(http.Header)
	h.Set("Host", host)
	h.Set("Content-Type", "application/json")
	h.Set("Keep-Alive", "true")
	return h
}
func body(params Params) []byte {
	b, _ := json.Marshal(params)
	return b
}
