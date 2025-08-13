package logger

import (
	"log"
	"os"
	"workluv/pkg/config"
)

// SetupLogger configures the logger based on the configuration
func SetupLogger(cfg *config.Config) {
	// Set log level based on configuration
	switch cfg.Log.Level {
	case "debug":
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	case "info":
		log.SetFlags(log.LstdFlags)
	case "warn", "warning":
		log.SetFlags(log.LstdFlags)
	case "error":
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	default:
		log.SetFlags(log.LstdFlags)
	}

	// Set output to stdout for better container logging
	log.SetOutput(os.Stdout)
}
