package helpers

import (
	"errors"
	"strings"
	"time"
)

// ParseETag 解析ETag值，一般是传入从Headers["Etag"]中获取的字符串
func ParseETag(rawEtag string) (string, error) {
	etag := strings.Trim(rawEtag, "\"")
	if len(etag) == 0 {
		return "", errors.New("empty etag")
	}
	return etag, nil
}

// ParseHTTPDateTime 解析HTTP Header中的datetime数据
func ParseHTTPDateTime(raw string) (time.Time, error) {
	return time.Parse(time.RFC1123, raw)
}
