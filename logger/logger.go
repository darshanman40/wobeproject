package logger

import (
	"io"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logInfo  *zap.Logger
	logWarn  *zap.Logger
	logErr   *zap.Logger
	logPanic *zap.Logger
	options  []field
	f        func()
	io.Writer
}

type field struct {
	key, value string
}

var logr Logger

//Logger ...
type Logger interface {
	Info(string, map[string]interface{})
	Warning(string, map[string]interface{})
	Error(string, map[string]interface{})
	Panic(string, map[string]interface{})
	CloseAll()
}

func (l *logger) Info(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.logInfo.Info(msg, zap.Int(k, v))

		case string:
			l.logInfo.Info(msg, zap.String(k, v))

		case float32:
			l.logInfo.Info(msg, zap.Float32(k, v))

		case float64:
			l.logInfo.Info(msg, zap.Float64(k, v))

		}
	}

}

func (l *logger) Warning(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.logWarn.Warn(msg, zap.Int(k, v))

		case string:
			l.logWarn.Warn(msg, zap.String(k, v))

		case float32:
			l.logWarn.Warn(msg, zap.Float32(k, v))

		case float64:
			l.logWarn.Warn(msg, zap.Float64(k, v))

		}
	}
}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.logErr.Error(msg, zap.Int(k, v))

		case string:
			l.logErr.Error(msg, zap.String(k, v))

		case float32:
			l.logErr.Error(msg, zap.Float32(k, v))

		case float64:
			l.logErr.Error(msg, zap.Float64(k, v))

		}
	}
}

func (l *logger) Panic(msg string, opts map[string]interface{}) {
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			l.logPanic.Panic(msg, zap.Int(k, v))

		case string:
			l.logPanic.Panic(msg, zap.String(k, v))

		case float32:
			l.logPanic.Panic(msg, zap.Float32(k, v))

		case float64:
			l.logPanic.Panic(msg, zap.Float64(k, v))

		}
	}
}

func (l *logger) CloseAll() {
	l.f()
}

//GetInstance ...
func GetInstance() Logger {
	if logr == nil {
		logr = NewLogger("")
	}
	return logr
}

//NewLogger ...
func NewLogger(flag string) Logger {

	//f, _ := os.Create("data/app.log")
	//ws, err := os.Open("data/app.json")
	ws, f, err := zap.Open("logs/app.json", "stderr")
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}
	var logInfo, logWarn, logErr, logPanic *zap.Logger
	switch flag {
	case "prod":
		// log = zap.New(
		// 	zapCoreConfig(ws, zapcore.DebugLevel),
		// 	//	zap.AddStacktrace(zapcore.DebugLevel),
		// 	zap.AddCallerSkip(1),
		// )

		logInfo = zap.New(
			zapCoreConfig(ws, zapcore.InfoLevel),
			//zap.ErrorOutput(zapcore.AddSync(ws)),
			//zap.AddCallerSkip(1),
		)

		logWarn = zap.New(
			zapCoreConfig(ws, zapcore.WarnLevel),
			// zap.ErrorOutput(zapcore.AddSync(ws)),
			// zap.AddStacktrace(zapcore.ErrorLevel),
			//zap.AddCallerSkip(1),
		)

		logErr = zap.New(
			zapCoreConfig(ws, zapcore.ErrorLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(2),
		)

		logPanic = zap.New(
			zapCoreConfig(ws, zapcore.PanicLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(4),
			//zap.AddStacktrace(zapcore.PanicLevel),
		)

	case "debug":

		logInfo = zap.New(
			zapCoreConfig(ws, zapcore.InfoLevel),
			//zap.AddCallerSkip(1),
			//			zap.ErrorOutput(zapcore.AddSync(ws)),
		)

		logWarn = zap.New(
			zapCoreConfig(ws, zapcore.WarnLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			// zap.ErrorOutput(zapcore.AddSync(ws)),
			// zap.AddStacktrace(zapcore.ErrorLevel),
		)

		logErr = zap.New(
			zapCoreConfig(ws, zapcore.ErrorLevel),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.AddCallerSkip(2),
		)

		logPanic = zap.New(
			zapCoreConfig(ws, zapcore.PanicLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(4),
		)
	case "dev":

		logInfo, _ = zap.NewDevelopment(zap.Development(),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)

		logWarn, _ = zap.NewDevelopment(zap.Development(),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)

		logErr, _ = zap.NewDevelopment(zap.Development(),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.ErrorOutput(ws),
			zap.AddCaller(),
			zap.AddCallerSkip(2),
		)

		logPanic, _ = zap.NewDevelopment(zap.Development(),
			zap.AddStacktrace(zapcore.PanicLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(4),
		)
	default:

		logInfo = zap.NewNop()
		logErr = zap.NewNop()
		logWarn = zap.NewNop()
		logPanic = zap.NewNop()
	}
	logr = &logger{
		logInfo:  logInfo,
		logWarn:  logWarn,
		logErr:   logErr,
		logPanic: logPanic,
		f:        f,
	}
	return logr
}

func zapCoreConfig(ws zapcore.WriteSyncer, l zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig()),
		ws,
		l,
	)

}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

}
