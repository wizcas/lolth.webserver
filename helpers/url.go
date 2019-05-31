package helpers

import (
	"net/url"
	"path"

	"quasar-ai.com/bast/marketing.landing/logger"
)

// BaseURL 站点根级URL
var BaseURL = EnvVar{Key: "BASE_URL", Required: true, NotEmpty: true}.GetString()

// ResolveURL 基于环境变量BASE_URL，解析相对路径的完整URL
func ResolveURL(relpath string) string {
	u, err := url.Parse(BaseURL)
	if err != nil {
		panic("invalid base URL")
	}
	u.Path = path.Join(u.Path, relpath)
	return u.String()
}

func init() {
	_ = ResolveURL("")
	logger.WithTag("PRE-FLIGHT").WriteInfo("BASE URL accepted.")
}
