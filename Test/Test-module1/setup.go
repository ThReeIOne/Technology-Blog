package Test_module1

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/engine"
	"Technology-Blog/Test/lib"
	"Technology-Blog/Test/log"
	"Technology-Blog/Test/provider"
	"bytes"

	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	code       = "1234456"
	token      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiUGhvbmUiOiIxMzY4ODg4OTk5OSIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJDZU1ldGEiLCJhdWQiOlsiR1ZBIl0sImV4cCI6MTcwMjAwNDQwMCwibmJmIjoxNzAxMzk5NTk5fQ.AA52LG45PzcvI2vbzuOywQhM6CZDfBN0ybrQpsCvMGs"
	adminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiUGhvbmUiOiIiLCJVc2VybmFtZSI6ImFkbWluIiwiUGFzc3dvcmQiOiIiLCJCdWZmZXJUaW1lIjo4NjQwMCwiaXNzIjoiQ2VNZXRhIiwiYXVkIjpbIkdWQSJdLCJleHAiOjE3MDIwMDQ0MDAsIm5iZiI6MTcwMTM5OTU5OX0.xGhMORpS_sfS8s5WWt6JIIIe8yAqSSQCXHN51PsbIOw"
)

var (
	user = model.User{
		Id:    1,
		Phone: "13688889999",
	}
	admin = model.Admin{
		Id:       1,
		Username: "admin",
		Password: lib.MD5("123456"),
	}
	assert = Assert{}
	jwt    = lib.NewJWT()
	now    = time.Date(2023, 12, 1, 3, 0, 0, 0, time.UTC)

	server = struct {
		User  controller.User
		Admin controller.Admin
		Role  controller.Role
	}{
		User:  controller.User{},
		Admin: controller.Admin{},
		Role:  controller.Role{},
	}

	// 初始化
	execute = []string{}
)

func init() {
	loadEnv()
	os.Setenv(config.EnableNetwork, "false")
	log.SetLevel(log.GetLevel(config.Get(config.LogLevel)))

	engine.Init()
	engine.Start()
}

func loadEnv() {
	configPath := os.Getenv("CONFIG")
	if len(configPath) == 0 {
		configPath = "../.env"
	}
	if err := godotenv.Load(configPath); err != nil {
		log.Fatal(err)
	}
}

func truncate() {
	if len(config.Get(config.IsDev)) == 0 {
		return
	}
	// 清库
	var all []string
	_ = provider.Database.DB.Select(&all, `
	SELECT Concat('TRUNCATE TABLE ',table_schema,'.',TABLE_NAME, ';')
	FROM INFORMATION_SCHEMA.TABLES where  table_schema in ('saas');
	`)
	for _, q := range all {
		_, err := provider.Database.DB.Exec(q)
		if err != nil {
			_ = log.Errorf("Execute fail %v %v", q, err)
		}
	}
	// 初始化
	for _, q := range execute {
		_, err := provider.Database.DB.Exec(q)
		if err != nil {
			_ = log.Errorf("Execute fail %v %v", q, err)
		}
	}
}

func setup() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	truncate()
	setCode()
	login(c, user.Id)
	return c, w
}

func login(c *gin.Context, uid int) {
	claims := jwt.CreateClaims(lib.BaseClaims{
		ID:    user.Id,
		Phone: user.Phone,
	})

	c.Set("claims", &claims)
	c.Set("x-user-id", strconv.Itoa(uid))
	c.Set(lib.HeaderXToken, token)
}

func before(c *gin.Context, method string, path string, data any, sql []string) {
	for _, q := range sql {
		_, _ = provider.Database.DB.Exec(q)
	}
	setRequest(c, method, path, data)
}

func setRequest(c *gin.Context, method string, path string, data any) {
	uid, _ := c.Get("x-user-id")
	token, _ := c.Get(lib.HeaderXToken)
	_data, _ := json.Marshal(data)
	setParams(c, path, data)
	if strings.ToUpper(method) == "GET" {
		path += "?" + params(data)
	}

	c.Request = httptest.NewRequest(method, path, bytes.NewBuffer(_data))
	c.Request.Header.Add("Content-Type", gin.MIMEJSON)
	c.Request.Header.Add("x-user-id", uid.(string))
	c.Request.Header.Add(lib.HeaderXToken, token.(string))
}

func params(data any) string {
	var params string
	d := reflect.ValueOf(data)
	for i := 0; i < d.NumField(); i++ {
		name := d.Type().Field(i).Tag.Get("json")
		value := d.Field(i)
		var param string
		switch value.Kind() {
		case reflect.Bool:
			param = strconv.FormatBool(value.Bool())
		case reflect.Int:
			param = strconv.Itoa(int(value.Int()))
		case reflect.String:
			param = value.String()
		}
		params += name + "=" + param + "&"
	}
	params, _ = strings.CutSuffix(params, "&")
	return params
}
func setParams(c *gin.Context, path string, data any) {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		paramKey := match[0][1:]
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			if tag == paramKey {
				paramValue := v.Field(i).Interface()
				strValue, ok := paramValue.(string)
				if ok {
					c.Params = append(c.Params, gin.Param{Key: paramKey, Value: strValue})
				} else {
					strValue = fmt.Sprintf("%v", paramValue)
					c.Params = append(c.Params, gin.Param{Key: paramKey, Value: strValue})
				}
				break
			}
		}
	}
}
func setCode() {
	provider.Cache.Pool.Get().Do("SET", fmt.Sprintf("sms_code_%s", user.Phone), code, "EX", 600)
}
