package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger  *logrus.Logger
	logFile *os.File
)

func Init() {
	logFilePath := "C:/Users/Admin/Desktop/test/logs/stack.log"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Không thể mở file log: %v", err)
	}

	// Tạo logger và chuyển sang JSON formatter
	logger := logrus.New()
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{})
	Logger = logger
	logFile = file
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}
