package app

import (
	"os"

	"go.uber.org/zap"
)

// Log is an implementaion of the zap logger 'go.uber.org/zap'
var Log *zap.Logger

func init() {
	switch os.Getenv("zap_logger_env") {
	case "Example":
		Log = zap.NewExample()
	case "Developer":
		Log, _ = zap.NewDevelopment()
	case "Sugar":
		sugar, _ := zap.NewDevelopment()
		slog := sugar.Sugar()
	default:
		Log, _ = zap.NewProduction()
	}
}
