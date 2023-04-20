package main

import (
	"encoding/gob"
	"net/url"
	"os"

	_ "github.com/lib/pq" // init pg driver
	"github.com/sirupsen/logrus"
)

func init() {
	// SetLevel sets the global log level used by the standard logger.
	logrus.SetLevel(getLogrusLogLevel(appLogLevel))

	// SetReportCaller sets whether the standard logger will include the calling
	// method as a field.
	logrus.SetReportCaller(true)

	// // Default formatter is the standard logger.
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	ForceColors: true,
	// })
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Register the types for gob
	gob.Register(url.Values{})
	gob.Register(map[string]string{})
	gob.Register(map[string][]string{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[string][]interface{}{})
	gob.Register(map[interface{}]interface{}{})
}

// getLogrusLogLevel returns logrus log level by string.
func getLogrusLogLevel(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
