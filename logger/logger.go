package logger

import (
	"errors"
	"strings"
	"fmt"
	"time"
	"io"
)
var ErrLevelNotRecognize = errors.New("unable to recognize level settings")
var LogLevels = [...]string {
	"trace",
	"debug",
	"info",
	"warning",
	"error",
	"fatal",
}
var wrlvl int
var wrlow io.Writer
var whigh io.Writer
func Init(level interface{}, w1 io.Writer, w2 io.Writer) error {
	switch v := level.(type) {
	case string:
		v = strings.ToLower(v)
		i:=0
		for i=0; i<len(LogLevels); i++ {
			if v == LogLevels[i] {
				break
			}
		}
		if i<len(LogLevels) {
			wrlvl = i
			goto INIT_SUCCESS
		}
	case int:
		if v >= 0 && v < len(LogLevels) {
			wrlvl = v
			goto INIT_SUCCESS
		}
	}
	return ErrLevelNotRecognize
INIT_SUCCESS:
	wrlow = w1
	whigh = w2
	return nil
}
func Trace(argv ...interface{}) {
	if wrlvl > 0 {
		return
	}
	fmt.Fprintf(wrlow, "[%s] (trace)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(wrlow, " %v", v)
	}
	fmt.Fprintf(wrlow, "\n")
}
func Debug(argv ...interface{}) {
	if wrlvl > 1 {
		return
	}
	fmt.Fprintf(wrlow, "[%s] (trace)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(wrlow, " %v", v)
	}
	fmt.Fprintf(wrlow, "\n")
}
func Info(argv ...interface{}) {
	if wrlvl > 2 {
		return
	}
	fmt.Fprintf(wrlow, "[%s] (trace)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(wrlow, " %v", v)
	}
	fmt.Fprintf(wrlow, "\n")
}
func Warning(argv ...interface{}) {
	if wrlvl > 3 {
		return
	}
	fmt.Fprintf(whigh, "[%s] (warning)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(whigh, " %v", v)
	}
	fmt.Fprintf(whigh, "\n")
}
func Warn(argv ...interface{}) {
	if wrlvl > 3 {
		return
	}
	fmt.Fprintf(whigh, "[%s] (warning)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(whigh, " %v", v)
	}
	fmt.Fprintf(whigh, "\n")
}
func Error(argv ...interface{}) {
	if wrlvl > 4 {
		return
	}
	fmt.Fprintf(whigh, "[%s] (error)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(whigh, " %v", v)
	}
	fmt.Fprintf(whigh, "\n")
}
func Fatal(argv ...interface{}) {
	if wrlvl > 5 {
		return
	}
	fmt.Fprintf(whigh, "[%s] (fatal)", time.Now().Format("2006-01-02 15:04:05"))
	for _, v := range argv {
		fmt.Fprintf(whigh, " %v", v)
	}
	fmt.Fprintf(whigh, "\n")
}
