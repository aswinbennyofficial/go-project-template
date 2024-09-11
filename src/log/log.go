package logs

import (
	"io"
	"log"
	"myapp/src/utils"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(config utils.LogConfig) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	
	var output io.Writer = os.Stdout

	if config.Output == "file" {
		log.Printf("Configuring file and console output: %s\n", config.File.Path)
		
		// Ensure the directory exists
		dir := filepath.Dir(config.File.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create log directory: %v\n", err)
		}
		
		fileOutput := &lumberjack.Logger{
			Filename:   config.File.Path,
			MaxSize:    config.File.MaxSize,
			MaxAge:     config.File.MaxAge,
			MaxBackups: config.File.MaxBackups,
			Compress:   true,
		}
		
		// Try to open the file to check if we have write permissions
		f, err := os.OpenFile(config.File.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open log file: %v\n", err)
		} else {
			f.Close()
		}
		
		// Use MultiWriter to write to both console and file
		output = zerolog.MultiLevelWriter(os.Stdout, fileOutput)
	} else {
		log.Println("Configuring console output only")
	}

	logger := zerolog.New(output).With().Timestamp().Logger()
	
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		log.Printf("Failed to parse log level '%s': %v\n", config.Level, err)
		level = zerolog.InfoLevel
	}
	
	return logger.Level(level)
}