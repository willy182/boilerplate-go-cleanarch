package utils

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	logrusSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

const (
	// TOPIC for setting topic of log
	TOPIC = "my-project-log"
	// LogTag default log tag
	LogTag = "my-project"
)

// LogContext function for logging the context of echo
// c string context
// s string scope
func LogContext(c string, s string) *log.Entry {
	return log.WithFields(log.Fields{
		"topic":      TOPIC,
		"context":    c,
		"scope":      s,
		"server_env": os.Getenv("SERVER_ENV"),
	})
}

// Log function for returning entry type
// level log.Level
// message string message of log
// context string context of log
// scope string scope of log
func Log(level log.Level, message string, context string, scope string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	log.SetFormatter(&log.JSONFormatter{})
	syslogOutput, err := logrusSyslog.NewSyslogHook("", "", syslog.LOG_INFO, LogTag)
	log.AddHook(syslogOutput)

	if err != nil {
		return
	}
	defer syslogOutput.Writer.Close()

	entry := LogContext(context, scope)
	switch level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.PanicLevel:
		entry.Panic(message)
	}
}

// LogError logging error
func LogError(err error, context string, messageData interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	log.SetFormatter(&log.JSONFormatter{})
	syslogOutput, errSys := logrusSyslog.NewSyslogHook("", "", syslog.LOG_INFO, LogTag)
	log.AddHook(syslogOutput)
	if errSys != nil {
		return
	}
	defer syslogOutput.Writer.Close()

	entry := log.WithFields(log.Fields{
		"topic":      TOPIC,
		"context":    context,
		"error":      err,
		"line_code":  TraceLineCode(),
		"server_env": os.Getenv("SERVER_ENV"),
	})

	jsonStr, _ := json.Marshal(messageData)
	entry.Error(string(jsonStr))
}

// TraceLineCode detect caller runtime
func TraceLineCode() string {
	var name string
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		name = fn.Name()
	}

	var githubLink string
	branch := os.Getenv("SERVER_ENV")
	if branch == "production" {
		branch = "master"
	}

	var prefix, suffix string
	sign := "boilerplate-go-cleanarch"
	i := strings.Index(name, sign)
	if i > 0 {
		prefix = name[:i+len(sign)]
	}
	i = strings.Index(file, sign)
	if i > 0 {
		suffix = file[i+len(sign):]
	}

	if prefix != "" && suffix != "" {
		githubLink = fmt.Sprintf("https://%s/blob/%s%s#L%d", prefix, branch, suffix, line)
	}
	return githubLink
}
