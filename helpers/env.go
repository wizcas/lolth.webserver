package helpers

import (
	"fmt"
	"os"
	"strconv"

	"quasar-ai.com/lolth/server.static/logger"
)

// EnvVar 用于获取环境变量的数据结构
type EnvVar struct {
	// Key 环境变量名称
	Key string
	// Required 是否必须通过环境变量提供，为`true`时如果找不到会panic，为`false`时则使用`DefaultValue`的值来解析
	Required bool
	// NotEmpty 是否必须提供变量值，为`true`时会在变量值为空字符串时panic
	NotEmpty bool
	// Default Value 若`Required`为`false`，系统中没找到该环境变量时所使用的默认值
	DefaultValue string
}

// GetString 获取该环境变量的字符串值
func (ev EnvVar) GetString() string {
	val, ok := os.LookupEnv(ev.Key)
	if !ok {
		if ev.Required {
			panic(ev.errNotFound())
		} else {
			val = ev.DefaultValue
		}
	}
	logger.WithTag("ENV").WriteInfo("%s = %s", ev.Key, val)
	if ev.NotEmpty && len(val) == 0 {
		panic(ev.errEmpty())
	}
	return val
}

// GetInt 获取该环境变量的整数值。如果环境变量不能被转换成整数（空、含有非数字字符等)，则会panic
func (ev EnvVar) GetInt() int {
	strval := ev.GetString()
	if len(strval) == 0 {
		panic(ev.errInvalidValue("cannot convert empty string into an integer"))
	}
	if val, err := strconv.ParseInt(strval, 10, 32); err != nil {
		panic(ev.errInvalidValue(err.Error()))
	} else {
		return int(val)
	}
}

func (ev EnvVar) errNotFound() error {
	return fmt.Errorf("environment variable is required: %s", ev.Key)
}
func (ev EnvVar) errEmpty() error {
	return fmt.Errorf("environment variable '%s' must not be empty", ev.Key)
}
func (ev EnvVar) errInvalidValue(err string) error {
	return fmt.Errorf("Invalid value of environment variable '%s': %v", ev.Key, err)
}
