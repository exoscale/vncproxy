package logger

import (
	"fmt"

	"gopkg.in/inconshreveable/log15.v2"
)

var defaultLogger Logger = &SimpleLogger{LogLevelInfo}

func GetLogLevel(logLevel string) LogLevel {
	switch logLevel {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	}
	return LogLevelInfo
}

type Logger interface {
	SetLogLevel(string)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type SimpleLogger struct {
	level LogLevel
}

func (sl *SimpleLogger) SetLogLevel(logLevel string) {
	level := GetLogLevel(logLevel)
	sl.level = level
}

func (sl *SimpleLogger) Debug(v ...interface{}) {
	if sl.level <= LogLevelDebug {
		arr := []interface{}{"[Debug]"}
		for _, item := range v {
			arr = append(arr, item)
		}

		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Debugf(format string, v ...interface{}) {
	if sl.level <= LogLevelDebug {
		fmt.Printf("[Debug] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Info(v ...interface{}) {
	if sl.level <= LogLevelInfo {
		arr := []interface{}{"[Info ]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Infof(format string, v ...interface{}) {
	if sl.level <= LogLevelInfo {
		fmt.Printf("[Info ] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Warn(v ...interface{}) {
	if sl.level <= LogLevelWarn {
		arr := []interface{}{"[Warn ]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Warnf(format string, v ...interface{}) {
	if sl.level <= LogLevelWarn {
		fmt.Printf("[Warn ] "+format+"\n", v...)
	}
}
func (sl *SimpleLogger) Error(v ...interface{}) {
	if sl.level <= LogLevelError {
		arr := []interface{}{"[Error]"}
		for _, item := range v {
			arr = append(arr, item)
		}
		fmt.Println(arr...)
	}
}
func (sl *SimpleLogger) Errorf(format string, v ...interface{}) {
	if sl.level <= LogLevelError {
		fmt.Printf("[Error] "+format+"\n", v...)
	}
}

type ExtendedLogger struct {
	logger log15.Logger
}

func NewExtendedLogger(l log15.Logger) *ExtendedLogger {
	if l == nil {
		l = log15.Root()
	}

	return &ExtendedLogger{logger: l}
}

func (el *ExtendedLogger) SetLogLevel(logLevel string) {
	level, _ := log15.LvlFromString(logLevel) // Falls back to level 'debug' in case of error
	el.logger.SetHandler(log15.LvlFilterHandler(level, log15.StdoutHandler))
}
func (el *ExtendedLogger) Debug(v ...interface{}) {
	el.logger.Debug(fmt.Sprint(v...))
}
func (el *ExtendedLogger) Debugf(format string, v ...interface{}) {
	el.logger.Debug(fmt.Sprintf(format, v...))
}
func (el *ExtendedLogger) Info(v ...interface{}) {
	el.logger.Info(fmt.Sprint(v...))
}
func (el *ExtendedLogger) Infof(format string, v ...interface{}) {
	el.logger.Info(fmt.Sprintf(format, v...))
}
func (el *ExtendedLogger) Warn(v ...interface{}) {
	el.logger.Warn(fmt.Sprint(v...))
}
func (el *ExtendedLogger) Warnf(format string, v ...interface{}) {
	el.logger.Warn(fmt.Sprintf(format, v...))
}
func (el *ExtendedLogger) Error(v ...interface{}) {
	el.logger.Error(fmt.Sprint(v...))
}
func (el *ExtendedLogger) Errorf(format string, v ...interface{}) {
	el.logger.Error(fmt.Sprintf(format, v...))
}

func SetLogger(i interface{}) {
	defaultLogger = i.(Logger)
}

func SetLogLevel(logLevel string) {
	defaultLogger.SetLogLevel(logLevel)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}
