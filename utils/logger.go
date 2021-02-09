package utils

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is an implementaion of the zap logger 'go.uber.org/zap'
var Log *zap.SugaredLogger

func init() {

	writerSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	var err error

	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
	if Log != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer Log.Sync()

}

func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	switch os.Getenv("go_logger") {
	case "Example":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	case "Developer":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	default:
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./hoa.log",
		MaxSize:    50,
		MaxBackups: 12,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
