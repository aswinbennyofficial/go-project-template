package logs

import (
	"io"
	"log"
	"myapp/src/config"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(config config.LogConfig) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	
	var output io.Writer = os.Stdout

	if config.Output == "file" {
		// Ensure the directory exists
		dir := filepath.Dir(config.File.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			// Use log package to print errors before zerolog is set up
			log.Printf("Failed to create log directory: %v\n", err)
		}
		
		fileOutput := &lumberjack.Logger{
			Filename:   config.File.Path,
			MaxSize:    config.File.MaxSize,    // Max megabytes before file is rotated
			MaxAge:     config.File.MaxAge,     // Max days to retain old log files
			MaxBackups: config.File.MaxBackups, // Max number of old log files to retain
			Compress:   true,                   // Compress log files
		}
		
		// Combine stdout and file output
		output = zerolog.MultiLevelWriter(os.Stdout, fileOutput)
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	
	// Set the log level
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		log.Printf("Invalid log level '%s', falling back to Info: %v\n", config.Level, err)
		level = zerolog.InfoLevel
	}

	return logger.Level(level)
}