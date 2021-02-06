package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
)

// Log is an implementaion of the zap logger 'go.uber.org/zap'
var Log *zap.Logger

func init() {
	logconfig := zap.NewDevelopmentConfig()
	var err error

	switch os.Getenv("zap_logger") {
	case "Example":
		Log = zap.NewExample()
	case "Developer":
		Log, err = zap.NewDevelopment()
	default:
		Log, err = zap.NewProduction()
		logconfig = zap.NewProductionConfig()
	}

	logconfig.OutputPaths = []string{"zap-demo.log"}
	Log, err = logconfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer Log.Sync()

}
