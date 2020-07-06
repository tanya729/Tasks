package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger
var once sync.Once

// Init logger
func Init(logLevel, filePath string) {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{
			filePath,
			"stderr",
		}
		config.ErrorOutputPaths = []string{
			filePath,
			"stderr",
		}
		config.EncoderConfig = zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		}

		var level zapcore.Level
		err := level.UnmarshalText([]byte(logLevel))
		if err != nil {
			log.Fatalf("can't marshal level string: %v", logLevel)
		}
		config.Level = zap.NewAtomicLevelAt(level)

		logger, err := config.Build()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
		defer logger.Sync()
		sugar = logger.Sugar()
	})
}

// Debug logs a debug message with the given fields
func Debug(message string, fields ...interface{}) {
	sugar.Debugw(message, fields...)
}

// Info logs a debug message with the given fields
func Info(message string, fields ...interface{}) {
	sugar.Infow(message, fields...)
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...interface{}) {
	sugar.Warnw(message, fields...)
}

// Error logs a debug message with the given fields
func Error(message string, fields ...interface{}) {
	sugar.Errorw(message, fields...)
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string, fields ...interface{}) {
	sugar.Fatalw(message, fields...)
}
