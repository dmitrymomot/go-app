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
	lvl, err := logrus.ParseLevel(appLogLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)

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
