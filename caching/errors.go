package caching

// ContentFetchError 内容获取失败的错误，如找不到页面、资源站无法访问等
type ContentFetchError string

func (e ContentFetchError) Error() string {
	return string(e)
}

// InternalError 服务器内部错误
type InternalError string

func (e InternalError) Error() string {
	return string(e)
}
