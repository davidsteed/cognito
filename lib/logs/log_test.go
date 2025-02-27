package logs

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZap(t *testing.T) {
	_,err:=NewZap("dev")
	if err !=nil{
		panic("cannot initiate zap logger")
	}

	assert.IsType(t, &zap.Logger{}, Log)

	Log.Debug("should not show")
	Log.Info("info")
}

func TestNewZapWithLevel(t *testing.T) {
	NewZapWithLevel("asdas")

	Log.Debug("should not show!!!!")
	Log.Info("info OK")

	err := SetLevel("debug")
	require.NoError(t, err)

	Log.Debug("debug OK")
	Log.Info("info OK")
	Log.Error("error OK", zap.Error(fmt.Errorf("some error")))
}

func TestInfoLogOmitsStacktrace(t *testing.T) {
	var buf bytes.Buffer

	NewZapWithLevel("dev")

	logger := Log.Sugar().Desugar()

	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.DebugLevel,
		))
	}))

	logger.Info("This is an info message")

	if bytes.Contains(buf.Bytes(), []byte("stacktrace")) {
		t.Errorf("INFO log message includes a stacktrace, it should not")
	}
}

func TestErrorLogIncludeStacktrace(t *testing.T) {
	var buf bytes.Buffer

	NewZapWithLevel("dev")

	logger := Log.Sugar().Desugar()

	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.DebugLevel,
		))
	}))

	logger.Error("This is an error message")

	if !bytes.Contains(buf.Bytes(), []byte("stacktrace")) {
		t.Errorf("ERROR log message doesn't includes a stacktrace, it should")
	}
}
