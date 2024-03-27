package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SaveToLogInfo(data any) {
	log := logrus.New()
	file, err := os.OpenFile("../application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	HandleErrorWithPanic(err)
	defer file.Close()
	log.SetOutput(file)
	log.Info(data)
}

func SaveToLogError(data any) {
	log := logrus.New()
	file, err := os.OpenFile("../application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	HandleErrorWithPanic(err)
	defer file.Close()
	log.SetOutput(file)
	log.Error(data)
}