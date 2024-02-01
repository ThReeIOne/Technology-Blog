package Test_module1

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
)

func TestLogin(t *testing.T) {
	//打桩操作，凡是调用time.Now，都返回now，打桩在自动化测试很常用
	patch := gomonkey.ApplyFunc(time.Now, func() time.Time { return now })
	defer patch.Reset()

	//后端需要的前端出入的参数
	type params struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	//返回给前端的
	type res struct {
		Id          int        `json:"id"`
		Phone       string     `json:"phone"`
		Balance     int        `json:"balance"`
		Token       string     `json:"token"`
		LastLoginAt *time.Time `json:"lastLoginAt"`
	}
	tests := []struct {
		name   string
		method string
		path   string
		in     params
		want   res
		status int
		sql    []string
	}{
		{
			name:   "正常登陆",
			method: "POST",
			path:   "/user/login",
			in: params{
				Phone: user.Phone,
				Code:  code,
			},
			want: res{
				Id:          1,
				Phone:       user.Phone,
				Balance:     50,
				Token:       token,
				LastLoginAt: &now,
			},
			status: http.StatusOK,
			sql:    []string{},
		},
		{
			name:   "无效的验证码",
			method: "POST",
			path:   "/user/login",
			in: params{
				Phone: user.Phone,
				Code:  "654321",
			},
			want:   res{},
			status: http.StatusBadRequest,
			sql:    []string{}, //在里面写入插入语句，执行操作数据库
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setup()
			before(c, tt.method, tt.path, tt.in, tt.sql)
			server.User.Login(c)

			assert.Nil(t, c.Err())
			assert.Equal(t, tt.status, w.Code)

			var data res
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &data))
			assert.Equal(t, data, tt.want)
		})
	}
}
