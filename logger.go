package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Values []Value

type Value struct {
	Key     string
	Payload interface{}
}

type Logger struct {
	logger zerolog.Logger
}

var Message = zerolog.MessageFieldName

func New() *Logger {
	return &Logger{
		logger: zerolog.New(
			zerolog.ConsoleWriter{
				Out: os.Stdout,
				FormatTimestamp: func(i interface{}) string {
					return time.Now().Format("2006-01-02 15:04:05")
				},
				FormatLevel: func(i interface{}) string {
					return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
				},
			},
		),
	}
}

func (l *Logger) Error(err error, values Values) {
	event := l.logger.Error().Err(err)
	event.Msg(l.enrichEventWithValues(event, values))
}

func (l *Logger) Panic(err error, values Values) {
	event := l.logger.Panic().Err(err)
	event.Msg(l.enrichEventWithValues(event, values))
}

func (l *Logger) Info(message string, values Values) {
	event := l.logger.Info()
	l.enrichEventWithValues(event, values)
	event.Msg(message)
}

func (l *Logger) enrichEventWithValues(event *zerolog.Event, values []Value) string {
	var msgPassed string

	for _, value := range values {
		if value.Key == zerolog.MessageFieldName {
			msgPassed = fmt.Sprintf("%v", value.Payload)
			continue
		}

		event = event.Interface(value.Key, fmt.Sprintf("%v", value.Payload))
	}

	return msgPassed
}
