package lib

import (
	"Technology-Blog/Test/config"
	"fmt"
	"strings"
)

func CdnPrefix() string {
	return config.Get(config.CDNPrefix)
}

func CdnUserPrefix() string {
	return fmt.Sprintf("%s/%s", CdnPrefix(), OssUserObjectPath)
}
func CdnSystemPrefix() string {
	return fmt.Sprintf("%s/%s", CdnPrefix(), OssSystemObjectPath)
}

func CdnSuffixPath(url string) string {
	prefix := fmt.Sprintf(`%s/`, CdnUserPrefix())
	if strings.HasPrefix(url, prefix) {
		return url[len(prefix):]
	}
	return url
}

func CdnParseUserUrl(path string) string {
	return fmt.Sprintf(`%s/%s`, CdnUserPrefix(), path)
}

func CdnParseSystemUrl(path string) string {
	return fmt.Sprintf(`%s/%s`, CdnSystemPrefix(), path)
}

func CdnUserObjectKey(userId int, path string) string {
	return fmt.Sprintf(`%s/%d/%s`, OssUserObjectPath, userId, path)
}
