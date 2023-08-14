package logger

import "github.com/sirupsen/logrus"

func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	logrus.SetLevel(logrus.InfoLevel)
}
