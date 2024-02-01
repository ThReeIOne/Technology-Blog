package httprequest

import (
	"net/http"
	"time"
)

var client *http.Client

// GetDefaultClient 单例获取http client
func GetDefaultClient() *http.Client {
	if client == nil {
		t := http.DefaultTransport.(*http.Transport).Clone()

		t.MaxIdleConns = 100        // 最大空闲连接数
		t.MaxConnsPerHost = 100     // 每个host可建立的最大连接数
		t.MaxIdleConnsPerHost = 100 // 每个目标host最大空闲连接数

		client = &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}
	}
	return client
}
