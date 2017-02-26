package logger

import "go.uber.org/zap"

type logger struct {
	log     *zap.Logger
	options []option
}

type option struct {
	key, value string
}

var logr Logger

//Logger ...
type Logger interface {
	Info(string, map[string]interface{})
	Error(string, map[string]interface{})
}

func (l *logger) Info(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.log.Info(msg, zap.Int(k, v))

		case string:
			l.log.Info(msg, zap.String(k, v))

		case float32:
			l.log.Info(msg, zap.Float32(k, v))

		case float64:
			l.log.Info(msg, zap.Float64(k, v))

		}
	}
}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.log.Error(msg, zap.Int(k, v))

		case string:
			l.log.Error(msg, zap.String(k, v))

		case float32:
			l.log.Error(msg, zap.Float32(k, v))

		case float64:
			l.log.Error(msg, zap.Float64(k, v))

		}
	}
}

//GetInstance ...
func GetInstance() Logger {
	if logr == nil {
		logr = &logger{
			log: zap.NewNop(),
		}

	}
	return logr
}
