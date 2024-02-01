// Copyright 2015 Eleme Inc. All rights reserved.

// Package log implements leveled logging.
package log

import (
	"Technology-Blog/Test/lib"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

// Level
const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
)

// Level name
var levelNames = [4]string{"debug", "info", "warning", "error"}

// Logging runtime
var (
	level             = INFO
	w       io.Writer = os.Stderr
	colored           = true
	enabled           = true
)

// colors to ansi code map
var colors = map[string]int{
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
}

// levelColors
var levelColors = map[int]string{
	DEBUG: "blue",
	INFO:  "green",
	WARN:  "yellow",
	ERROR: "red",
}

// SetColored sets the color enability.
func SetColored(b bool) {
	colored = b
}

func GetLevel(lName string) int {
	for i, v := range levelNames {
		if v == lName {
			return i
		}
	}
	return INFO
}

// SetLevel sets the logging level.
func SetLevel(l int) {
	level = l % len(levelNames)
}

// SetWriter sets the writer.
func SetWriter(writer io.Writer) {
	w = writer
}

// Disable the logging.
func Disable() {
	enabled = false
}

// Enable the logging.
func Enable() {
	enabled = true
}

// Debug logs message with level DEBUG.
func Debug(a ...interface{}) error {
	return log(DEBUG, fmt.Sprint(a...))
}

// Info logs message with level INFO.
func Info(a ...interface{}) error {
	return log(INFO, fmt.Sprint(a...))
}

// Warn logs message with level WARN.
func Warn(a ...interface{}) error {
	return log(WARN, fmt.Sprint(a...))
}

// Error logs message with level ERROR.
func Error(a ...interface{}) error {
	return log(ERROR, fmt.Sprint(a...))
}

// Fatal and logs message with level FATAL.
func Fatal(a ...interface{}) {
	_ = log(ERROR, fmt.Sprint(a...))
	os.Exit(1)
}

// Debugf formats and logs message with level DEBUG.
func Debugf(format string, a ...interface{}) error {
	return log(DEBUG, fmt.Sprintf(format, a...))
}

// Infof formats and logs message with level INFO.
func Infof(format string, a ...interface{}) error {
	disable := os.Getenv("DISABLE_INFO")
	if len(disable) > 0 {
		return nil
	}
	return log(INFO, fmt.Sprintf(format, a...))
}

// Warnf formats and logs message with level WARN.
func Warnf(format string, a ...interface{}) error {
	return log(WARN, fmt.Sprintf(format, a...))
}

// Errorf formats and logs message with level ERROR.
func Errorf(format string, a ...interface{}) error {
	return log(ERROR, fmt.Sprintf(format, a...))
}

// Fatalf formats and logs message with level FATAL.
func Fatalf(format string, a ...interface{}) {
	_ = log(ERROR, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Colored returns text in color.
func Colored(color string, text string) string {
	return fmt.Sprintf("\033[3%dm%s\033[0m", colors[color], text)
}

// log dose logging.
func log(l int, msg string) error {
	if enabled && l >= level {
		// Caller pkg.
		_, fileName, line, _ := runtime.Caller(2)
		pkgName := path.Base(path.Dir(fileName))
		filepath := path.Join(pkgName, path.Base(fileName))
		// Datetime and pid.
		now := time.Now().Format(time.RFC3339)
		// Message
		level := levelNames[l]
		header := Colored(levelColors[l], fmt.Sprintf("time=\"%s\" level=%s method=%s:%d", now, level, filepath, line))
		_, err := fmt.Fprintf(w, "%s msg=\"%s\"\n", header, msg)
		if l >= WARN && !lib.IsDev() {
			alter(l, fmt.Sprintf("time=\"%s\" level=%s method=%s:%d msg=\"%s\n", now, level, filepath, line, msg))
		}
		return err
	}
	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func alter(l int, msg string) {
	title := levelNames[l]
	content := []*feishu.FsSegment{
		{Tag: feishu.FsTagA, Text: msg},
	}
	c := feishu.NewFsRichContent(&title, [][]*feishu.FsSegment{content})
	if err := feishu.SendMessage(feishu.SaasBackendErrorBotId, c); err != nil {
		fmt.Println(err)
	}
}
