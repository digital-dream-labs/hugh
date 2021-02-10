// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// setSyslogFormatter is nil if the target architecture does not support syslog.
var setSyslogFormatter func(logger, string, string) error

// setEventlogFormatter is nil if the target OS does not support Eventlog (i.e., is not Windows).
var setEventlogFormatter func(logger, string, bool) error

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

func setJSONFormatter() {
	origLogger.Formatter = &logrus.JSONFormatter{}
}

type loggerSettings struct {
	level  string
	format string
}

func (s *loggerSettings) apply(ctx *kingpin.ParseContext) error {
	err := baseLogger.SetLevel(s.level)
	if err != nil {
		return err
	}
	err = baseLogger.SetFormat(s.format)
	return err
}

// AddFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(a *kingpin.Application) {
	s := loggerSettings{}
	a.Flag("log.level", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]").
		Default(origLogger.Level.String()).
		StringVar(&s.level)
	defaultFormat := url.URL{Scheme: "logger", Opaque: "stderr"}
	a.Flag("log.format", `Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"`).
		Default(defaultFormat.String()).
		StringVar(&s.format)
	a.Action(s.apply)
}

// Logger is the interface for loggers used in the Prometheus components.
type Logger interface {
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	With(key string, value interface{}) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithStructs(args ...interface{}) Logger
	WithMetadata(args ...interface{}) Logger

	SetFormat(string) error
	SetLevel(string) error
}

type logger struct {
	entry *logrus.Entry
}

func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

func (l logger) WithField(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

func (l logger) WithFields(fields Fields) Logger {
	return logger{l.entry.WithFields(logrus.Fields(fields))}
}

func (l logger) WithMetadata(args ...interface{}) Logger {
	allFields := make(Fields)

	for in := range args {
		var fields Fields
		inrec, err := json.Marshal(&args[in])
		if err == nil {
			uerr := json.Unmarshal(inrec, &fields)
			if uerr == nil {
				for k, v := range fields {
					allFields[k] = v
				}
			}
		}
	}

	return logger{l.entry.WithFields(logrus.Fields(allFields))}
}

func (l logger) WithStructs(args ...interface{}) Logger {
	allFields := make(Fields)

	for in := range args {
		var fields Fields
		inrec, err := json.Marshal(&args[in])
		if err == nil {
			uerr := json.Unmarshal(inrec, &fields)
			if uerr == nil {
				for k, v := range fields {
					allFields[k] = v
				}
			}
		}
	}

	return logger{l.entry.WithFields(logrus.Fields(allFields))}
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l logger) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

func (l logger) SetLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	l.entry.Logger.Level = lvl
	return nil
}

func (l logger) SetFormat(format string) error {
	u, err := url.Parse(format)
	if err != nil {
		return err
	}
	if u.Scheme != "logger" {
		return fmt.Errorf("invalid scheme %s", u.Scheme)
	}
	jsonq := u.Query().Get("json")
	if jsonq == "true" {
		setJSONFormatter()
	}

	switch u.Opaque {
	case "syslog":
		if setSyslogFormatter == nil {
			return fmt.Errorf("system does not support syslog")
		}
		appname := u.Query().Get("appname")
		facility := u.Query().Get("local")
		return setSyslogFormatter(l, appname, facility)
	case "eventlog":
		if setEventlogFormatter == nil {
			return fmt.Errorf("system does not support eventlog")
		}
		name := u.Query().Get("name")
		debugAsInfo := false
		debugAsInfoRaw := u.Query().Get("debugAsInfo")
		if parsedDebugAsInfo, err := strconv.ParseBool(debugAsInfoRaw); err == nil {
			debugAsInfo = parsedDebugAsInfo
		}
		return setEventlogFormatter(l, name, debugAsInfo)
	case "stdout":
		l.entry.Logger.Out = os.Stdout
	case "stderr":
		l.entry.Logger.Out = os.Stderr
	default:
		return fmt.Errorf("unsupported logger %q", u.Opaque)
	}
	return nil
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	fn, _, _, ok := runtime.Caller(2)
	if !ok {
		return l.entry.WithField("source", "N/A")
	}
	frames := runtime.CallersFrames([]uintptr{fn})
	frame, _ := frames.Next()
	filidx := strings.LastIndex(frame.File, "/")
	file := frame.File[filidx+1:]
	funidx := strings.LastIndex(frame.Function, "/")
	fun := frame.Function[funidx+1:]

	return l.entry.WithFields(logrus.Fields{"source": fmt.Sprintf("%s:%d", file, frame.Line), "function": fun})
}

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// Base returns the default Logger logging to
func Base() Logger {
	return baseLogger
}

// NewLogger returns a new Logger logging to out.
func NewLogger(w io.Writer) Logger {
	l := logrus.New()
	l.Out = w
	return logger{entry: logrus.NewEntry(l)}
}

// NewNopLogger returns a logger that discards all log messages.
func NewNopLogger() Logger {
	l := logrus.New()
	l.Out = ioutil.Discard
	return logger{entry: logrus.NewEntry(l)}
}

// With adds a field to the logger.
func With(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// WithField Added to make swap from logrus easier.
func WithField(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// WithFields adds fields to a Logger.
func WithFields(fields Fields) Logger {
	return baseLogger.WithFields(fields)
}

// WithStructs adds all fields of structs to a Logger.
func WithStructs(args ...interface{}) Logger {
	return baseLogger.WithStructs(args...)
}

// WithMetadata adds all fields of structs to a Logger.
func WithMetadata(args ...interface{}) Logger {
	return baseLogger.WithMetadata(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.sourced().Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.sourced().Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.sourced().Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.sourced().Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.sourced().Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced().Fatalf(format, args...)
}

// AddHook adds hook to Prometheus' original logger.
func AddHook(hook logrus.Hook) {
	origLogger.Hooks.Add(hook)
}

type errorLogWriter struct{}

func (errorLogWriter) Write(b []byte) (int, error) {
	baseLogger.sourced().Error(string(b))
	return len(b), nil
}

// NewErrorLogger returns a log.Logger that is meant to be used
// in the ErrorLog field of an http.Server to log HTTP server errors.
func NewErrorLogger() *log.Logger {
	return log.New(&errorLogWriter{}, "", 0)
}

// SetTimestampFormat sets the timestamp format for the underlying logrus
// instance.
func SetTimestampFormat(ts string) {
	timeFormatter := new(logrus.TextFormatter)
	timeFormatter.TimestampFormat = ts
	timeFormatter.FullTimestamp = true
	origLogger.Formatter = timeFormatter
}

// SetJSONFormat sets the json format for the underlying logrus instance.
func SetJSONFormat(ts string) {
	jsonFormatter := new(logrus.JSONFormatter)
	jsonFormatter.TimestampFormat = ts
	origLogger.Formatter = jsonFormatter
}
