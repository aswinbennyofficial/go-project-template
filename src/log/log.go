package log

import (
	"io"
	"myapp/src/utils"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)



func NewLogger(config utils.LogConfig) zerolog.Logger {
    zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

    var output io.Writer = os.Stdout
    if config.Output == "file" {
        output = &lumberjack.Logger{
            Filename:   config.File.Path,
            MaxSize:    config.File.MaxSize,
            MaxAge:     config.File.MaxAge,
            MaxBackups: config.File.MaxBackups,
            Compress:   true,
        }
    }

    logger := zerolog.New(output).With().Timestamp().Logger()
    level, _ := zerolog.ParseLevel(config.Level)
    return logger.Level(level)
}
