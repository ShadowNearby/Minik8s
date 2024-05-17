package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var defaultLogFormat = "[%lvl%] [%time% | %func% | %file%]: %msg%\n"
var defaultTimestampFormat = "15:04:05"

type CustomFormatter struct {
	TimestampFormat string
	LogFormat       string
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		fileVal = strings.Replace(fileVal, RootPath+"/", "", 1)
		output = strings.Replace(output, "%func%", funcVal, 1)
		output = strings.Replace(output, "%file%", fileVal, 1)
	} else {
		output = strings.Replace(output, "%func%", "", 1)
		output = strings.Replace(output, "%file%", "", 1)
	}

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}
