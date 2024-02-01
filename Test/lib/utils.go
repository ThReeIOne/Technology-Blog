package lib

import (
	"Technology-Blog/Test/config"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsDev() bool {
	debug, _ := strconv.ParseBool(config.Get(config.IsDev))
	return debug
}

func IsEnableNetwork() bool {
	enable, _ := strconv.ParseBool(config.Get(config.EnableNetwork))
	return enable
}

func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		day, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(day)
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}

func RandomId() string {
	id := time.Now().Format("20060102150405.000000")
	return strings.ReplaceAll(id, ".", "")
}

func StrToPtr(s string) *string {
	return &s
}

func ImgToBase64(url string) (string, error) {
	var bodyImg io.Reader
	req, err := http.NewRequest("GET", url, bodyImg)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	img := base64.StdEncoding.EncodeToString(b)
	mimeType := http.DetectContentType(b)
	switch mimeType {
	case "image/jpg":
		img = "data:image/jpg;base64," + img
	case "image/png":
		img = "data:image/png;base64," + img
	case "image/jpeg":
		img = "data:image/jpeg;base64," + img
	}
	return img, nil
}

func IntToBool(i int) bool {
	return i > 0
}
func BoolToInt(boo bool) int {
	if boo {
		return 1
	}
	return 0
}

func StringProcess(str string, processors ...func(string) string) string {
	for i := range processors {
		str = processors[i](str)
	}
	return str
}

func StringsProcess(ss []string, processors ...func(string) string) []string {
	result := []string{}
	for i := range ss {
		result = append(result, StringProcess(ss[i], processors...))
	}
	return result
}

// GetOffset 获取偏移量 X:100|Y:200|W:1290|H:1080
func GetOffset(offset string, key string) int {
	all := strings.Split(offset, "|")
	for _, v := range all {
		vl := strings.Split(v, ":")
		if strings.ToUpper(vl[0]) != strings.ToUpper(key) {
			continue
		}
		if len(vl) > 1 {
			value, _ := strconv.Atoi(vl[1])
			return value
		} else {
			return 0
		}
	}
	return 0
}

func GetImageBounds(imgUrl string) (int, int, error) {
	if !regexp.MustCompile(ImageUrlRegex).MatchString(imgUrl) {
		return 0, 0, errors.New(fmt.Sprintf("invalid url %s", imgUrl))
	}
	var bodyImg io.Reader
	req, err := http.NewRequest("GET", imgUrl, bodyImg)
	if err != nil {
		return 0, 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}
