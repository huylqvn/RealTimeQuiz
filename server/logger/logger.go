package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.elastic.co/apm/v2"
)

type Log struct {
	l           log.Logger
	isCaller    bool
	isTimestamp bool
	isFile      bool
	level       string
	source      string
}

var ll *Log
var once sync.Once

func Get() *Log {
	once.Do(func() {
		l := log.NewJSONLogger(os.Stdout)
		ll = &Log{
			l:           l,
			isCaller:    false,
			isTimestamp: false,
			isFile:      false,
		}
	})
	return ll
}

func (l *Log) GetLogger() log.Logger {
	return l.l
}

func (l *Log) SetCaller() *Log {
	l.isCaller = true
	l.l = log.With(l.l, "caller", log.Caller(6))
	return l
}

func (l *Log) SetTimestamp() *Log {
	l.isTimestamp = true
	l.l = log.With(l.l, "timestamp", log.TimestampFormat(time.Now, time.RFC3339))
	return l
}

func (l *Log) SetLevel(lvl string) *Log {
	l.level = lvl
	keyLevel := level.ParseDefault(lvl, level.InfoValue())
	l.l = level.NewFilter(l.l, level.Allow(keyLevel))
	return l
}

func (l *Log) SetSource(source string) *Log {
	l.source = source
	l.l = log.With(l.l, "source", source)
	return l
}

func (l *Log) SetFile() *Log {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "./logs"
	}
	l.isFile = true
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/%s.log", logPath, "%Y-%m-%d"),
		rotatelogs.WithMaxAge(time.Hour*24*10),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		l.l.Log("message", "Failed to Initialize Log File")
	}

	w := log.NewSyncWriter(writer)
	mw := io.MultiWriter(os.Stdout, w)
	l.l = log.NewJSONLogger(mw)
	if l.isCaller {
		l.SetCaller()
	}
	if l.isTimestamp {
		l.SetTimestamp()
	}
	if l.source != "" {
		l.SetSource(l.source)
	}
	l.SetLevel(l.level)

	return l
}

func (l *Log) Log(keyvals ...interface{}) error {
	return l.l.Log(keyvals...)
}

func (l *Log) Info(traceKey, message string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "trace_key", traceKey)
	keyvals = append(keyvals, "message", message)
	return level.Info(l.l).Log(keyvals...)
}

func (l *Log) Error(traceKey, message string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "trace_key", traceKey)
	keyvals = append(keyvals, "message", message)
	return level.Error(l.l).Log(keyvals...)
}

func (l *Log) Debug(traceKey, message string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "trace_key", traceKey)
	keyvals = append(keyvals, "message", message)
	return level.Debug(l.l).Log(keyvals...)
}

func (l *Log) Warn(traceKey, message string, keyvals ...interface{}) error {
	keyvals = append(keyvals, "trace_key", traceKey)
	keyvals = append(keyvals, "message", message)
	return level.Warn(l.l).Log(keyvals...)
}

func (l *Log) Trace(traceKey string, timeStart time.Time, metadata map[string]interface{}, err error, span *apm.Span) error {
	defer span.End()
	keyvals := []interface{}{
		"trace_key", traceKey,
		"time_used", fmt.Sprintf("%v", time.Since(timeStart)),
		"metadata", metadata,
	}
	span.Context.SetLabel("trace_key", traceKey)
	if len(metadata) > 0 {
		for k, v := range metadata {
			span.Context.SetLabel(k, v)
		}
	}
	if err != nil {
		keyvals = append(keyvals, "error", err)
		keyvals = append(keyvals, "message", "failed")
		e := apm.DefaultTracer().NewError(err)
		e.SetSpan(span)
		e.Send()
	} else {
		keyvals = append(keyvals, "message", "success")
	}
	return l.l.Log(keyvals...)
}
