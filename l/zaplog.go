package l

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/sunary/kitchen/id"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Short-hand functions for logging.
var (
	Any         = zap.Any
	Bool        = zap.Bool
	ByteString  = zap.ByteString
	ByteStrings = zap.ByteStrings
	Duration    = zap.Duration
	Float32     = zap.Float32
	Float64     = zap.Float64
	Int         = zap.Int
	Int32       = zap.Int32
	Int64       = zap.Int64
	Uint        = zap.Uint
	Uint32      = zap.Uint32
	Uint64      = zap.Uint64
	Uintptr     = zap.Uintptr
	Skip        = zap.Skip
	String      = zap.String
	Stringer    = zap.Stringer
	Time        = zap.Time
)

// ID wraps UUID.
func ID(id id.UUID) zapcore.Field {
	return String("id", id.String())
}

// Error wraps error for zap.Error.
func Error(err error) zapcore.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}

// Stack ...
func Stack() zapcore.Field {
	return zap.Stack("stack")
}

// Object ...
func Object(key string, val interface{}) zapcore.Field {
	return zap.Stringer(key, Dump(val))
}

type obj struct {
	v interface{}
}

// String ...
func (o obj) String() string {
	return pp.Sprint(o.v)
}

// Dump renders object for debugging
func Dump(v interface{}) fmt.Stringer {
	return obj{v}
}

// Interface ...
func Interface(key string, val interface{}) zapcore.Field {
	if val, ok := val.(fmt.Stringer); ok {
		return zap.Stringer(key, val)
	}
	return zap.Reflect(key, val)
}

// New returns new zap.Logger
func New() Logger {
	envLog := os.Getenv("LOG_LEVEL")
	if envLog == "" {
		envLog = "DEBUG"
	}

	var lv zapcore.Level
	err := lv.UnmarshalText([]byte(envLog))
	if err != nil {
		panic("go-common/l: " + err.Error())
	}

	cfg := zap.Config{
		Encoding:         ConsoleEncoderName,
		Level:            zap.NewAtomicLevelAt(lv),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",

			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}
	logger, _ := cfg.Build()
	return Logger{logger}
}

func init() {
	err := zap.RegisterEncoder(ConsoleEncoderName, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewConsoleEncoder(cfg), nil
	})
	if err != nil {
		panic(err)
	}
}
