package helpers

import (
	"errors"
	"strings"
)

// ParseETag 解析ETag值，一般是传入从Headers["Etag"]中获取的字符串
func ParseETag(rawEtag string) (string, error) {
	etag := strings.Trim(rawEtag, "\"")
	if len(etag) == 0 {
		return "", errors.New("empty etag")
	}
	return etag, nil
}
