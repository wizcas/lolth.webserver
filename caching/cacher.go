package caching

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"quasar-ai.com/bast/marketing.landing/helpers"

	"quasar-ai.com/bast/marketing.landing/logger"
)

// 环境变量
var (
	FetchTimeout    = helpers.EnvVar{Key: "FETCH_TIMEOUT", DefaultValue: "3"}.GetInt()
	RefreshInterval = helpers.EnvVar{Key: "REFRESH_INTERVAL", DefaultValue: "30"}.GetInt()
)

type contentCache struct {
	sync.WaitGroup
	url   string
	hash  string
	data  []byte
	timer *helpers.Timer
}
type tplData struct {
	// BaseURL 站点根级URL
	BaseURL string
}

var mapCache = make(map[string]*contentCache)

// GetCachedContent 根据相对路径获取缓存的内容，如果内容有更新，将会在刷新间隔之后更新缓存
func GetCachedContent(relpath string) (data []byte, err error) {
	url := helpers.ResolveURL(relpath)
	cache, ok := mapCache[url]
	if !ok {
		cache = &contentCache{url: url, timer: helpers.NewTimer(time.Duration(RefreshInterval) * time.Second)}
		mapCache[url] = cache
		logger.WriteDebug("new cache: %s", url)
	}
	err = nil
	if cache.isObsoleted() {
		logger.WriteDebug("refreshing cache: %s", cache.url)
		err = cache.refresh()
	}
	data = cache.data
	return
}

func (c *contentCache) getRemoteHash() (string, error) {
	res, err := http.Head(c.url)
	if err != nil {
		return "", err
	}
	return helpers.ParseETag(res.Header.Get("Etag"))
}

func (c *contentCache) isObsoleted() bool {
	c.Wait()
	// ignore check when in refresh intervals
	if c.timer.IsEnabled() && !c.timer.IsTimeUp() {
		return false
	}
	if len(c.hash) <= 0 {
		// newly created cache
		return true
	}
	c.Add(1)
	defer c.Done()
	remoteHash, err := c.getRemoteHash()
	if err != nil {
		logger.WriteError("get remote hash failed: %v", err)
		return false
	}
	logger.WriteDebug("local hash: %s, remote hash: %s", c.hash, remoteHash)
	sameHash := remoteHash == c.hash
	if sameHash {
		c.timer.Renew()
	}
	return !sameHash
}

func (c *contentCache) refresh() error {
	c.Add(1)
	chErr := make(chan error, 1)
	go c.doRefresh(chErr)
	var err error
	select {
	case err = <-chErr:
		break
	case <-time.After(time.Duration(FetchTimeout) * time.Second):
		err = ContentFetchError("time out")
	}
	if err != nil {
		logger.WriteError("refresh cache error: %v\n(@ %s)", err, c.url)
	}
	return err
}

func (c *contentCache) doRefresh(chErr chan error) {
	defer c.Done()
	res, err := http.Get(c.url)
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("GET failed: %v", err))
		return
	}
	logger.WriteDebug("header: %v\n", res.Header)
	hash, err := helpers.ParseETag(res.Header.Get("Etag"))
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("parse ETag failed: %v", err))
		return
	}
	c.hash = hash
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		chErr <- InternalError(fmt.Sprintf("ready body data failed: %v", err))
		return
	}
	tpl, err := template.New("remote").Parse(string(data))
	var buf bytes.Buffer
	tpl.Execute(&buf, tplData{helpers.BaseURL})
	c.data = buf.Bytes()
	c.timer.Renew()
	chErr <- nil
	close(chErr)
}
