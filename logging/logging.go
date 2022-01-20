package logging

import (
	"log"
	"strings"
)

const (
	errorPrefix   string = "[ERROR]"
	warningPrefix string = "[WARN]"
	infoPrefix    string = "[INFO]"
	debugPrefix   string = "[DEBUG]"
	tracePrefix   string = "[TRACE]"
)

func init() {
	log.SetFlags(0) // Removes all logger prefixes
}

func prefixLogArgs(prefix string, args ...interface{}) []interface{} {
	return append([]interface{}{prefix}, args...)
}

func prefixLogFormat(prefix, format string) string {
	var builder strings.Builder
	builder.WriteString(prefix)
	builder.WriteString(" ")
	builder.WriteString(format)
	return builder.String()
}

func Error(args ...interface{}) {
	log.Println(prefixLogArgs(errorPrefix, args...)...)
}

func Errorf(format string, args ...interface{}) {
	log.Printf(prefixLogFormat(errorPrefix, format), args...)
}

func Warning(args ...interface{}) {
	log.Println(prefixLogArgs(warningPrefix, args...)...)
}

func Warningf(format string, args ...interface{}) {
	log.Printf(prefixLogFormat(warningPrefix, format), args...)
}

func Info(args ...interface{}) {
	log.Println(prefixLogArgs(infoPrefix, args...)...)
}

func Infof(format string, args ...interface{}) {
	log.Printf(prefixLogFormat(infoPrefix, format), args...)
}

func Debug(args ...interface{}) {
	log.Println(prefixLogArgs(debugPrefix, args...)...)
}

func Debugf(format string, args ...interface{}) {
	log.Printf(prefixLogFormat(debugPrefix, format), args...)
}

func Trace(args ...interface{}) {
	log.Println(prefixLogArgs(tracePrefix, args...)...)
}

func Tracef(format string, args ...interface{}) {
	log.Printf(prefixLogFormat(tracePrefix, format), args...)
}
