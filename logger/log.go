package logger

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type logLevel interface {
	GetCode() int
	GetLiteral() string
}

var levelLiteral = map[int]string{
	int(Debug):     "Debug",
	int(Info):      "Info",
	int(Notice):    "Notice",
	int(Warning):   "Warning",
	int(Error):     "Error",
	int(Critical):  "Critical",
	int(Alert):     "Alert",
	int(Emergency): "Emergency",
}

//RFC5424SeverityLevel  RFC5424SeverityLevel
type RFC5424SeverityLevel int

//GetCode logLevel interface
func (s RFC5424SeverityLevel) GetCode() int {
	return int(s)
}

//GetLiteral logLevel interface
func (s RFC5424SeverityLevel) GetLiteral() string {
	if literal, ok := levelLiteral[int(s)]; ok {
		return literal
	}
	return "Unknown"
}

var _ logLevel = RFC5424SeverityLevel(0)

// 日志的错误级别
const (
	Debug     = RFC5424SeverityLevel(7)
	Info      = RFC5424SeverityLevel(6)
	Notice    = RFC5424SeverityLevel(5)
	Warning   = RFC5424SeverityLevel(4)
	Error     = RFC5424SeverityLevel(3)
	Critical  = RFC5424SeverityLevel(2)
	Alert     = RFC5424SeverityLevel(1)
	Emergency = RFC5424SeverityLevel(0)
)

// Log 日志结构体
type Log struct {
	Level logLevel
	Tag   string
	Msg   string
	File  string
	Line  int
}

//WithTag 带tag 的日志记录
func WithTag(tag string) *Log {
	return &Log{
		Tag: tag,
	}
}

//WriteDebug WriteDebug
func WriteDebug(format string, v ...interface{}) {
	Write(Debug, true, format, v...)
}

//WriteDebug WriteDebug
func (l *Log) WriteDebug(format string, v ...interface{}) {
	l.Write(Debug, true, format, v...)
}

//WriteInfo WriteInfo
func WriteInfo(format string, v ...interface{}) {
	Write(Info, false, format, v...)
}

//WriteInfo WriteInfo
func (l *Log) WriteInfo(format string, v ...interface{}) {
	l.Write(Info, false, format, v...)
}

//WriteNotice WriteNotice
func WriteNotice(format string, v ...interface{}) {
	Write(Notice, true, format, v...)
}

//WriteNotice WriteNotice
func (l *Log) WriteNotice(format string, v ...interface{}) {
	l.Write(Notice, true, format, v...)
}

//WriteWarning WriteWarning
func WriteWarning(format string, v ...interface{}) {
	Write(Warning, true, format, v...)
}

//WriteWarning WriteWarning
func (l *Log) WriteWarning(format string, v ...interface{}) {
	l.Write(Warning, true, format, v...)
}

//WriteError WriteError
func WriteError(format string, v ...interface{}) {
	Write(Error, true, format, v...)
}

//WriteError WriteError
func (l *Log) WriteError(format string, v ...interface{}) {
	l.Write(Error, true, format, v...)
}

//WriteCritical WriteCritical
func WriteCritical(format string, v ...interface{}) {
	Write(Critical, true, format, v...)
}

//WriteCritical WriteCritical
func (l *Log) WriteCritical(format string, v ...interface{}) {
	l.Write(Critical, true, format, v...)
}

//WriteAlert WriteAlert
func WriteAlert(format string, v ...interface{}) {
	Write(Alert, true, format, v...)
}

//WriteAlert WriteAlert
func (l *Log) WriteAlert(format string, v ...interface{}) {
	l.Write(Alert, true, format, v...)
}

//WriteEmergency WriteEmergency
func WriteEmergency(format string, v ...interface{}) {
	Write(Emergency, true, format, v...)
}

//WriteEmergency WriteEmergency
func (l *Log) WriteEmergency(format string, v ...interface{}) {
	l.Write(Emergency, true, format, v...)
}

//Write write log
func (l Log) Write(level logLevel, trace bool, format string, v ...interface{}) {
	var file string
	var line int
	if trace {
		// 这里需要判断调用的层级关系,否则直接调用Write 报错位置是错误的
		// TODO:应该想办法通过包名来判断是否,而且最好从2开始倒着遍历效率更高
		_, currentfile, _, _ := runtime.Caller(0)
		for i := 1; i < 3; i++ {
			_, file, line, _ = runtime.Caller(i)
			if file != currentfile {
				break
			}
		}
		l.File = file
		l.Line = line
	}
	l.Level = level
	l.Msg = fmt.Sprintf(format, v...)
	doWrite(&l, trace)
}

//Write 记录日志
func Write(level logLevel, trace bool, format string, v ...interface{}) {
	var file string
	var line int
	if trace {
		// 这里需要判断调用的层级关系,否则直接调用Write 报错位置是错误的
		// TODO:应该想办法通过包名来判断是否,而且最好从2开始倒着遍历效率更高
		_, currentfile, _, _ := runtime.Caller(0)
		for i := 1; i < 3; i++ {
			_, file, line, _ = runtime.Caller(i)
			if file != currentfile {
				break
			}
		}
	}
	l := &Log{
		Level: level,
		Msg:   fmt.Sprintf(format, v...),
		File:  file,
		Line:  line,
	}
	doWrite(l, trace)
}

func doWrite(l *Log, trace bool) {
	if trace {
		log.Println(fmt.Sprintf("%s [%s] [%s] %s --- %s:%d", time.Now().Format(time.RFC3339), l.Level.GetLiteral(), l.Tag, l.Msg, l.File, l.Line))
	} else {
		log.Println(fmt.Sprintf("%s [%s] [%s] %s", time.Now().Format(time.RFC3339), l.Level.GetLiteral(), l.Tag, l.Msg))
	}
}
