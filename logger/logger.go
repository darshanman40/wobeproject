package logger

import (
	"io"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logInfo   *zap.Logger
	logWarn   *zap.Logger
	logErr    *zap.Logger
	logPanic  *zap.Logger
	fileClose func()
	io.Writer
}

type logMessages struct {
	log    *zap.Logger
	msg    string
	fields map[string]interface{}
}

var logInfoChan, logWarnChan, logErrChan, logPanicChan chan logMessages

// var mutex sync.Mutex

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
	// zapOpts := GetFields(opts)
	// l.logInfo.Info(msg, zapOpts...)
	logInfoChan <- logMessages{log: l.logInfo, msg: msg, fields: opts}
}

func (l *logger) Warning(msg string, opts map[string]interface{}) {
	// zapOpts := GetFields(opts)
	// l.logWarn.Warn(msg, zapOpts...)
	logWarnChan <- logMessages{log: l.logWarn, msg: msg, fields: opts}
}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	// zapOpts := GetFields(opts)
	// l.logErr.Error(msg, zapOpts...)
	logErrChan <- logMessages{log: l.logErr, msg: msg, fields: opts}
}

func (l *logger) Panic(msg string, opts map[string]interface{}) {
	// zapOpts := GetFields(opts)
	// l.logPanic.Panic(msg, zapOpts...)
	logPanicChan <- logMessages{log: l.logPanic, msg: msg, fields: opts}
}

//GetFields ...
func GetFields(opts map[string]interface{}) []zapcore.Field {
	zapOptions := make([]zapcore.Field, len(opts))
	i := 0
	for k, v := range opts {
		switch v := v.(type) {
		case int:
			zapOptions[i] = zap.Int(k, v)
		case int8:
			zapOptions[i] = zap.Int8(k, v)
		case int16:
			zapOptions[i] = zap.Int16(k, v)
		case int32:
			zapOptions[i] = zap.Int32(k, v)
		case int64:
			zapOptions[i] = zap.Int64(k, v)
		case string:
			zapOptions[i] = zap.String(k, v)
		case float32:
			zapOptions[i] = zap.Float32(k, v)
		case float64:
			zapOptions[i] = zap.Float64(k, v)
		case error:
			zapOptions[i] = zap.Error(v)

		}
		i++
	}
	return zapOptions
}

func (l *logger) CloseAll() {
	l.fileClose()
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
	var ws zapcore.WriteSyncer
	var f func()
	var err error

	var logInfo, logWarn, logErr, logPanic *zap.Logger
	// mutex.Lock()

	switch flag {
	case "prod":
		ws, f, err = zap.Open("logs/app.json", "stderr")
		logInfo = zap.New(
			zapCoreConfig(ws, zapcore.InfoLevel),
		)

		logWarn = zap.New(
			zapCoreConfig(ws, zapcore.WarnLevel),
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
		)

	case "debug":
		ws, f, err = zap.Open("logs/app.json", "stderr")
		logInfo = zap.New(
			zapCoreConfig(ws, zapcore.InfoLevel),
		)

		logWarn = zap.New(
			zapCoreConfig(ws, zapcore.WarnLevel),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
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
		ws, f, err = zap.Open("logs/app.json", "stderr")
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
		_, f, err = zap.Open("stderr")
		logInfo = zap.NewNop()
		logErr = zap.NewNop()
		logWarn = zap.NewNop()
		logPanic = zap.NewNop()
	}
	if err != nil {
		log.Fatal("Can't open file at logger ", err)
	}
	logr = &logger{
		logInfo:   logInfo,
		logWarn:   logWarn,
		logErr:    logErr,
		logPanic:  logPanic,
		fileClose: f,
	}

	logInfoChan = make(chan logMessages)
	logWarnChan = make(chan logMessages)
	logErrChan = make(chan logMessages)
	logPanicChan = make(chan logMessages)
	go LogRoutine()
	// defer mutex.Unlock()
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

//LogRoutine ...
func LogRoutine() {
	// go func() {
	for {
		// mutex.Lock()
		select {
		case l := <-logInfoChan:
			opts := GetFields(l.fields)
			l.log.Info(l.msg, opts...)

		case l := <-logWarnChan:
			opts := GetFields(l.fields)
			l.log.Warn(l.msg, opts...)

		case l := <-logErrChan:
			opts := GetFields(l.fields)
			l.log.Error(l.msg, opts...)

		case l := <-logPanicChan:
			opts := GetFields(l.fields)
			l.log.Panic(l.msg, opts...)

		}
		// mutex.Unlock()
	}
	// }()
}
