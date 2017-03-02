package logger_test

import (
	"errors"
	"log"
	"testing"

	"github.com/wobeproject/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_GetFields(t *testing.T) {
	err := errors.New("custom error")
	originalFields := []zapcore.Field{
		zap.Int("int", 9),
		zap.Int8("int8", 1),
		zap.Int16("int16", 16),
		zap.Int32("int32", 32),
		zap.Int64("int64", 64),
		zap.String("string", "stringdata"),
		zap.Float32("float32", 2.56455),
		zap.Float64("float64", 3.5565645),
		zap.Error(err),
	}

	m := map[string]interface{}{
		"string":  "stringdata",
		"int":     int(9),
		"int8":    int8(1),
		"int16":   int16(16),
		"int32":   int32(32),
		"int64":   int64(64),
		"float32": float32(2.56455),
		"float64": float64(3.5565645),
		"err":     err,
	}
	fields := logger.GetFields(m)
	log.Println("fields: ", fields)
	log.Println("originalFields: ", originalFields)
	if !testEq(fields, originalFields) {
		t.Fatalf("list mismatch")
	}
}

func testEq(a, b []zapcore.Field) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	temp := false
	for _, i := range a {
		for _, j := range b {
			if i != j {
				temp = false
				continue
			}
			temp = true
			break
		}
		if !temp {
			return temp
		}
	}

	return temp
}
