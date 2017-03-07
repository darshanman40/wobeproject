package logger_test

import (
	"errors"
	"testing"

	"github.com/wobeproject/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_GetFields(t *testing.T) {
	err := errors.New("custom error")
	originalFields := []zapcore.Field{
		zap.Int("intKey", 9),
		zap.Int8("int8Key", 1),
		zap.Int16("int16Key", 16),
		zap.Int32("int32Key", 32),
		zap.Int64("int64Key", 64),
		zap.String("stringKey", "stringdata"),
		zap.Float32("float32Key", 2.56455),
		zap.Float64("float64Key", 3.5565645),
		zap.Error(err),
	}

	m := map[string]interface{}{
		"stringKey":  "stringdata",
		"intKey":     int(9),
		"int8Key":    int8(1),
		"int16Key":   int16(16),
		"int32Key":   int32(32),
		"int64Key":   int64(64),
		"float32Key": float32(2.56455),
		"float64Key": float64(3.5565645),
		"errKey":     err,
	}
	fields := logger.GetFields(m)
	if !testEqField(fields, originalFields) {
		t.Fatalf("list mismatch")
	}
}

//testEqField is created to compare unorganized slice of zapcore.Field
func testEqField(a, b []zapcore.Field) bool {

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
