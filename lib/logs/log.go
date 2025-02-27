package logs

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log     *zap.Logger
	atom    zap.AtomicLevel
	version string
)

// NewZapWithLevel set logger with initial fields and level
func NewZapWithLevel(appVersion string) {
	version = appVersion
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	if appVersion == "dev" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	stacktraceLevel := zap.NewAtomicLevelAt(zapcore.ErrorLevel)

	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"

	core := zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stderr),
		atom,
	)

	option := zap.AddStacktrace(stacktraceLevel)

	Log = zap.New(core, option).With(zap.String("version", appVersion))
}

// Deprecated: use NewZapWithLevel: see examples/lambda-openapi/cmd/lambda-openapi/main.go
func NewZap(appVersion string) (*zap.Logger, error) {
	NewZapWithLevel(appVersion)

	return Log, nil
}

func SetLevel(level string) error {
	var logLevel zapcore.Level

	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("failed to unmarshal log level")
	}

	atom.SetLevel(logLevel)

	return nil
}

// SilentLogs is useful for test functions - Sample code below:
//
//	if testing.Verbose() {
//		logs.NewZapWithLevel(Version)
//	} else {
//
//		logs.SilentLogs()
//	}
func SilentLogs() {
	Log = zap.NewNop()
}

// Version returns the appVersion set in NewZapWithLevel
func AppVersion() string {
	return version
}
