package main

import (
	"errors"
	"fmt"
	"net/http"

	"quasar-ai.com/bast/marketing.landing/logger"

	"quasar-ai.com/bast/marketing.landing/helpers"

	"github.com/gin-gonic/gin"
	"quasar-ai.com/bast/marketing.landing/caching"
)

var PORT = helpers.EnvVar{Key: "PORT", DefaultValue: "8080"}.GetInt()
var INDEX_PATH = helpers.EnvVar{Key: "INDEX_PATH", DefaultValue: "index.html"}.GetString()

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", index)
	logger.WithTag("STARTUP").WriteInfo("server is listening on %d", PORT)
	r.Run(fmt.Sprintf(":%d", PORT))
}

func index(c *gin.Context) {
	if data, err := caching.GetCachedContent(INDEX_PATH); err != nil {
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
