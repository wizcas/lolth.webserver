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
	url        string
	modifiedAt time.Time
	data       []byte
	timer      *helpers.Timer
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

func (c *contentCache) getLastModified() (time.Time, error) {
	res, err := http.Head(c.url)
	if err != nil {
		return time.Time{}, err
	}
	return helpers.ParseHTTPDateTime(res.Header.Get("Last-Modified"))
}

func (c *contentCache) isObsoleted() bool {
	c.Wait()
	// ignore check when in refresh intervals
	if c.timer.IsEnabled() && !c.timer.IsTimeUp() {
		return false
	}
	if c.modifiedAt.IsZero() {
		// newly created cache
		return true
	}
	c.Add(1)
	defer c.Done()
	remoteModified, err := c.getLastModified()
	if err != nil {
		logger.WriteError("get remote header failed: %v", err)
		return false
	}
	logger.WriteDebug("local modified: %s, remote modified: %s", c.modifiedAt, remoteModified)
	sameTime := remoteModified.Equal(c.modifiedAt)
	if sameTime {
		c.timer.Renew()
	}
	return !sameTime
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
	defer close(chErr)
	res, err := http.Get(c.url)
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("GET failed: %v", err))
		return
	}
	remoteModified, err := helpers.ParseHTTPDateTime(res.Header.Get("Last-Modified"))
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("parse Last-Modified header failed: %v", err))
		return
	}
	c.modifiedAt = remoteModified
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		chErr <- InternalError(fmt.Sprintf("ready body data failed: %v", err))
		return
	}
	tpl, err := template.New("remote").Parse(string(data))
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("invalid template: %v", err))
		return
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, tplData{helpers.BaseURL})
	if err != nil {
		chErr <- ContentFetchError(fmt.Sprintf("render template error: %v", err))
		return
	}
	c.data = buf.Bytes()
	c.timer.Renew()
	chErr <- nil
}
