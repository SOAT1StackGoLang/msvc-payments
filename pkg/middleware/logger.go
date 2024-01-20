package logger

import (
	"log"
	"os"
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

// initializeLogger initializes the logger using go-kit.
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

	if os.Getenv("APP_LOG_LEVEL") == "debug" {
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

	DebugLogger.Log("message", "debug logger initialized")
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
