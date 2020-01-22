package grok

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Init setup default loggin configuration
func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
}
