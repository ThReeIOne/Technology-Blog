package chatgpt

import (
	"Technology-Blog/Test/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	BaseUrl             = "https://api.ai.staringos.com"
	ModelChatGpt4       = "gpt-4"
	ArticleGeneratePath = "/docs-ai/chatByCorpus"
	ArticleAnswerPath   = "/message"
)

type Variables struct {
	ProductName    *string `json:"chanpinmingcheng,omitempty"` // 产品名称
	TargetAudience *string `json:"mubiaorenqun,omitempty"`     // 目标人群
	ArticleType    *string `json:"wenanleixing,omitempty"`     // 文案类型
	ArticleStyle   *string `json:"wenanfengge,omitempty"`      // 文案风格
	SellingPoint   *string `json:"chanpinmaidian,omitempty"`   // 卖点信息
	KeyWords       *string `json:"guanjianci,omitempty"`       // 关键字
	Long           *string `json:"wenanchangdu,omitempty"`     // 文案长短
}

type Params struct {
	Model     string     `json:"model"`               // 模型 gpt-4
	CorpusId  int        `json:"corpusId"`            // 语料库ID
	Prompt    *string    `json:"prompt,omitempty"`    // prompt
	Variables *Variables `json:"variables,omitempty"` // 模版参数
}

type Response struct {
	Id       string `json:"id"`
	Content  string `json:"content"`
	IsFinish bool   `json:"isFinish"`
	Model    string `json:"model,omitempty"`
	AppId    int    `json:"appId,omitempty"`
	At       uint64 `json:"at,omitempty"`
	T        string `json:"t,omitempty"`
	// ContentFrom string `json:"t,omitempty"` // LLM
}

func headers() http.Header {
	h := make(http.Header)
	h.Set("Host", "api.ai.staringos.com")
	h.Set("Content-Type", "application/json")
	h.Set("Authorization", fmt.Sprintf("Bearer %s", config.Get(config.ChatGptToken)))
	return h
}
func body(params Params) string {
	b, _ := json.Marshal(params)
	return string(b)
}

func Generate(params Params) (Response, error) {
	req, err := http.NewRequest("POST", BaseUrl+ArticleGeneratePath, strings.NewReader(body(params)))
	if err != nil {
		return Response{}, err
	}
	req.Header = headers()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}
	var r Response
	err = json.Unmarshal(data, &r)
	return r, err
}

func Answer(answerId string) (Response, error) {
	req, err := http.NewRequest("GET", BaseUrl+ArticleAnswerPath+"?id="+answerId, nil)
	if err != nil {
		return Response{}, err
	}
	req.Header = headers()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}
	var r Response
	err = json.Unmarshal(data, &r)
	return r, err
}
