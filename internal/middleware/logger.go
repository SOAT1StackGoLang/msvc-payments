package logger

import (
	"log"
	"os"
	"strings"
	"time"

	kitlog "github.com/go-kit/log"
	kitloglevel "github.com/go-kit/log/level"
	kitlogterm "github.com/go-kit/log/term"
)

// create Debug, Info and error loggers
var (
	DebugLogger kitlog.Logger
	InfoLogger  kitlog.Logger
	ErrorLogger kitlog.Logger
)

// InitializeLogger initializes the logger for the microservice.
// It sets the logger level based on the environment variable "APP_LOG_LEVEL".
// If the log level is set to "debug", it creates a logger that allows debug logs.
// Otherwise, it creates a logger that allows error, warn and info level log events to pass.
// The logger is configured with a custom timestamp format and caller information.
// It also creates separate loggers for debug, info, and error logs.
// Finally, it sets the logger as the output for the standard library log package.
func InitializeLogger() {
	// set logger level
	var logger kitlog.Logger

	colorFn := func(keyvals ...interface{}) kitlogterm.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] != "level" {
				continue
			}
			switch keyvals[i+1] {
			case "debug":
				return kitlogterm.FgBgColor{Fg: kitlogterm.DarkGray}
			case "info":
				return kitlogterm.FgBgColor{Fg: kitlogterm.Blue}
			case "warn":
				return kitlogterm.FgBgColor{Fg: kitlogterm.Yellow}
			case "error":
				return kitlogterm.FgBgColor{Fg: kitlogterm.Red}
			case "crit":
				return kitlogterm.FgBgColor{Fg: kitlogterm.Gray, Bg: kitlogterm.DarkRed}
			default:
				// to fix, some cases are using default instead of the correct level
				// light blue is the default color
				return kitlogterm.FgBgColor{Fg: kitlogterm.Cyan}
			}
		}
		return kitlogterm.FgBgColor{}
	}

	if strings.ToLower(os.Getenv("APP_LOG_LEVEL")) == "debug" {
		logger = kitlogterm.NewLogger(kitlog.NewSyncWriter(os.Stdout), kitlog.NewLogfmtLogger, colorFn)
		logger = kitloglevel.NewFilter(logger, kitloglevel.AllowDebug())
	} else {
		logger = kitlogterm.NewLogger(kitlog.NewSyncWriter(os.Stdout), kitlog.NewLogfmtLogger, colorFn)
		logger = kitloglevel.NewFilter(logger, kitloglevel.AllowInfo())
	}

	// set logger timestamp and caller
	customTimestampFormat := kitlog.TimestampFormat(func() time.Time {
		return time.Now().UTC()
	}, "2006/01/02-15:04:05")
	logger = kitlog.With(logger, "ts", customTimestampFormat)
	// logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.Caller(4))

	// create Debug, Info and error loggers
	DebugLogger = kitloglevel.Debug(logger)
	InfoLogger = kitloglevel.Info(logger)
	ErrorLogger = kitloglevel.Error(logger)

	log.SetOutput(kitlog.NewStdlibAdapter(logger))

	DebugLogger.Log("message", "DEBUG MODE ON: initialized")
}

// Helper functions
func Debug(msg string) {
	DebugLogger.Log("message", msg)
}

func Info(msg string) {
	InfoLogger.Log("message", msg)
}

func Error(msg string) {
	ErrorLogger.Log("message", msg)
}
