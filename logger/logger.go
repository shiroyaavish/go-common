package logger

import (
	"encoding/json"
	"fmt"
	"github.com/IntelXLabs-LLC/go-common/utils"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

var Logger *zerolog.Logger

// init initialised the logger
func init() {
	Logger = utils.Pointer(zerolog.New(os.Stdout).With().Timestamp().Logger())
}

// Trace provides stack strace
func Trace(debugStack, format string, v ...interface{}) {
	Logger.Trace().Str("stack", debugStack).Msgf(format, v...)
}

// Info provides info log
func Info(format string, v ...interface{}) {
	Logger.Info().Msgf(format, v...)
}

// Errorf provides error level log
func Errorf(format string, v ...interface{}) {
	Logger.Error().Msgf(format, v...)
}

// Error provides error level log
func Error(err error, msg ...string) {
	Logger.Error().Msgf(fmt.Sprintf("%s\n %s", strings.Join(msg, "\n"), err.Error()))
}

func Fatal(err error, msg ...string) {
	Error(err, msg...)
	os.Exit(1)
}

func Data(data interface{}) {
	v, _ := json.Marshal(data)
	Logger.Info().Msgf(fmt.Sprintf("Data: %s", string(v)))
}

// Warn function logs a warning message using the provided format and arguments.
// It formats the message using the format and v arguments and then logs the message at the warning level
func Warn(format string, v ...interface{}) {
	Logger.Warn().Msgf(format, v...)
}

// Debug prints a debug-level message with formatting.
// It uses the globally defined Logger to output the message with the desired format.
// The format parameter specifies the layout of the message.
// The v parameter is an optional variadic argument that can be used
func Debug(format string, v ...interface{}) {
	Logger.Debug().Msgf(format, v...)
}
