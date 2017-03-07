package logger

import (
	"log"
	"sync"

	"github.com/wobeproject/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logInfo   *zap.Logger
	logWarn   *zap.Logger
	logErr    *zap.Logger
	logPanic  *zap.Logger
	logDebug  *zap.Logger
	fileClose func()
}

type logMessages struct {
	log    *zap.Logger
	msg    string
	fields map[string]interface{}
}

var (
	logInfoChan, logWarnChan, logErrChan, logPanicChan chan logMessages

	env string

	mutex sync.Mutex

	logr Logger
)

//Logger to implement zap logger
type Logger interface {
	Info(string, map[string]interface{})
	Warning(string, map[string]interface{})
	Error(string, map[string]interface{})
	Panic(string, map[string]interface{})
	Debug(string, map[string]interface{})
	CloseAll()
}

func (l *logger) Info(msg string, opts map[string]interface{}) {
	zapOpts := GetFields(opts)
	l.logInfo.Info(msg, zapOpts...)
}

func (l *logger) Warning(msg string, opts map[string]interface{}) {
	zapOpts := GetFields(opts)
	l.logWarn.Warn(msg, zapOpts...)
}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	zapOpts := GetFields(opts)
	l.logErr.Error(msg, zapOpts...)
}

func (l *logger) Panic(msg string, opts map[string]interface{}) {
	zapOpts := GetFields(opts)
	l.logPanic.Panic(msg, zapOpts...)
}

func (l *logger) Debug(msg string, opts map[string]interface{}) {
	zapOpts := GetFields(opts)
	l.logDebug.Debug(msg, zapOpts...)
}

//GetFields conver maps into zapcore.Fields
func GetFields(fields map[string]interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for k, v := range fields {
		switch v := v.(type) {
		case int:
			zapFields[i] = zap.Int(k, v)
		case int8:
			zapFields[i] = zap.Int8(k, v)
		case int16:
			zapFields[i] = zap.Int16(k, v)
		case int32:
			zapFields[i] = zap.Int32(k, v)
		case int64:
			zapFields[i] = zap.Int64(k, v)
		case string:
			zapFields[i] = zap.String(k, v)
		case float32:
			zapFields[i] = zap.Float32(k, v)
		case float64:
			zapFields[i] = zap.Float64(k, v)
		case error:
			zapFields[i] = zap.Error(v)
		}
		i++
	}
	return zapFields
}

func (l *logger) CloseAll() {
	l.fileClose()
}

//GetInstance to retrieve single instance of Logger
func GetInstance() Logger {
	if logr == nil {
		logr = NewLogger(nil)
	}
	return logr
}

//NewLogger Get new instance of Logger
func NewLogger(l map[string]config.Log) Logger {

	ws, f, err := zap.Open("stderr")
	if err != nil {
		log.Fatal("Error at logger, ", err)
	}
	if l == nil {
		logr = &logger{
			logInfo:   zap.NewNop(),
			logWarn:   zap.NewNop(),
			logErr:    zap.NewNop(),
			logPanic:  zap.NewNop(),
			logDebug:  zap.NewNop(),
			fileClose: f,
		}

	} else {

		logr = &logger{
			logInfo:   GetZapLogger(ws, l["info"]),
			logWarn:   GetZapLogger(ws, l["warn"]),
			logErr:    GetZapLogger(ws, l["err"]),
			logPanic:  GetZapLogger(ws, l["panic"]),
			logDebug:  GetZapLogger(ws, l["debug"]),
			fileClose: f,
		}
	}

	logInfoChan = make(chan logMessages)
	logWarnChan = make(chan logMessages)
	logErrChan = make(chan logMessages)
	logPanicChan = make(chan logMessages)
	LogRoutine()
	return logr
}

//GetZapLogger ..
func GetZapLogger(ws zapcore.WriteSyncer, l config.Log) *zap.Logger {
	var z []zap.Option

	var zc zapcore.Core
	var zl zapcore.Level
	if l.Tracelevel != "" {
		zl = getLevel(l.Tracelevel)
	}

	if l.Erroroutput {
		z = append(z, zap.ErrorOutput(ws))
	}

	if l.Stacktrace {
		z = append(z, zap.AddStacktrace(zl))
	}

	if l.Caller {
		z = append(z, zap.AddCaller())
	}

	z = append(z, zap.AddCallerSkip(l.CallerSkip))

	zc = zapCoreConfig(ws, zl)
	newZap := zap.New(zc, z...)
	return newZap

}

func getLevel(s string) zapcore.Level {
	switch s {
	case "infolevel":
		return zapcore.InfoLevel
	case "warnlevel":
		return zapcore.WarnLevel
	case "errorlevel":
		return zapcore.ErrorLevel
	case "paniclevel":
		return zapcore.PanicLevel
	case "debuglevel":
		return zapcore.DebugLevel
	}
	return 10
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
	go func() {
		for {
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
		}
	}()
}
