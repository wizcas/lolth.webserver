package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"quasar-ai.com/lolth/server.static/logger"

	"quasar-ai.com/lolth/server.static/helpers"

	"github.com/gin-gonic/gin"
	"quasar-ai.com/lolth/server.static/caching"
)

// Port 服务器监听端口
var Port = helpers.EnvVar{Key: "PORT", DefaultValue: "8080"}.GetInt()

// IndexPath 首页文件相对路径(相对于BaseUrl)
var IndexPath = helpers.EnvVar{Key: "INDEX_PATH", DefaultValue: "index.html"}.GetString()

var assetSuffixes = []string{
	".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico",
	".js", ".css", ".txt",
	".woff2",
	".webmanifest",
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(subpathMiddleware(r))
	r.GET("/", index)
	logger.WithTag("STARTUP").WriteInfo("server is listening on %d", Port)
	r.Run(fmt.Sprintf(":%d", Port))
}

func index(c *gin.Context) {
	if data, err := caching.GetCachedContent(IndexPath); err != nil {
		var status int
		switch err.(type) {
		case caching.ContentFetchError:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		c.AbortWithError(status, err)
	} else if data == nil {
		c.AbortWithError(http.StatusNoContent, errors.New("cached data is empty"))
	} else {
		c.Data(200, "text/html", data)
	}
}

func subpathMiddleware(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		relpath := c.Request.URL.Path
		if redirectToAsset(relpath, c) {
			return
		}
		if relpath == "/" {
			logger.WithTag("REQ").WriteDebug("ROOT")
			c.Next()
		} else {
			logger.WithTag("REQ").WriteDebug("REDIRECT TO ROOT")
			c.Request.URL.Path = "/"
			r.HandleContext(c)
		}
	}
}

func redirectToAsset(relpath string, c *gin.Context) bool {
	logger.WithTag("REQ").WriteDebug("URL: %s", relpath)
	for _, suffix := range assetSuffixes {
		if strings.HasSuffix(relpath, suffix) {
			c.Redirect(http.StatusMovedPermanently, helpers.ResolveURL(relpath))
			return true
		}
	}
	return false
}
