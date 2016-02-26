package wlog

import (
	"fmt"
	"os"
)

type Redirector int
type LogLevel int

const (
	Stdout Redirector = iota
	Stderr
)

const (
	Debug LogLevel = iota
	Info
	Notice
	Warn
	Error
	Crit
)

type Logger struct {
	Rdr   Redirector
	Level LogLevel
}

func New(r Redirector, level string) Logger {
	return Logger{
		Rdr:   r,
		Level: string2Level(level),
	}
}

func AccessLogger(level string) Logger {
	return New(Stdout, level)
}

func ErrorLogger(level string) Logger {
	return New(Stderr, level)
}

func (l *Logger) Out(level LogLevel, msg string, args ...interface{}) {
	if l.Rdr == Stderr {
		if level >= l.Level {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(fmtMessage(level, msg), args...))
		}
	} else {
		fmt.Fprintln(os.Stdout, fmt.Sprintf(fmtMessage(level, msg), args...))
	}
}

func string2Level(s string) LogLevel {
	var result LogLevel
	switch s {
	case "debug":
		result = Debug
	case "notice":
		result = Notice
	case "warn":
		result = Warn
	case "error":
		result = Error
	case "crit":
		result = Crit
	case "info":
		fallthrough
	default:
		result = Info
	}
	return result
}

func level2String(level LogLevel) string {
	var result string
	switch level {
	case Debug:
		result = "debug"
	case Notice:
		result = "notice"
	case Warn:
		result = "warn"
	case Error:
		result = "error"
	case Crit:
		result = "crit"
	case Info:
		fallthrough
	default:
		result = "info"
	}
	return result
}

func fmtMessage(level LogLevel, msg string) string {
	return fmt.Sprintf("level:%s\t%s", level2String(level), msg)
}
