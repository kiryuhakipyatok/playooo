package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"io"
)

func NewLogger() *logrus.Logger{
	log:=logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.TextFormatter{})
	file, err := os.OpenFile("../../files/logging/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.WithError(err).Info("failed to open log file, using default stdout")
		log.SetOutput(os.Stdout)
	}else{
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	}
	return log
}