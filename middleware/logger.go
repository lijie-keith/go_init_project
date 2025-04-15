package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lijie-keith/go_init_project/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

func LoggerToFile() gin.HandlerFunc {
	logFilePath := config.LOG_FILE_PATH
	logFileName := config.LOG_FILE_NAME

	fileName := path.Join(logFilePath, logFileName)
	_, err := PathExists(fileName)
	if err != nil {
		err := os.MkdirAll(fileName, os.ModePerm)
		if err != nil {
			fmt.Println("create log file error", err)
		}
	}
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("error", err)
	}

	writers := []io.Writer{
		src,
		os.Stdout}

	fileAndStdoutWriter := io.MultiWriter(writers...)

	config.SystemLogger.SetOutput(fileAndStdoutWriter)
	config.SystemLogger.SetLevel(getLevel())

	logWriter, err := rotatelogs.New(
		fileName+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationTime(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	ifHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	config.SystemLogger.AddHook(ifHook)
	config.SystemLogger.SetReportCaller(true)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		method := c.Request.Method
		reqURI := c.Request.RequestURI
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		// 日志格式
		config.SystemLogger.WithFields(logrus.Fields{
			"status_code":  status,
			"latency_time": latency,
			"client_ip":    clientIP,
			"req_method":   method,
			"req_uri":      reqURI,
		}).Info()
	}
}

func getLevel() logrus.Level {
	level := logrus.DebugLevel
	switch config.LOG_LEVEL {
	case 0:
		level = logrus.PanicLevel
	case 1:
		level = logrus.FatalLevel
	case 2:
		level = logrus.ErrorLevel
	case 3:
		level = logrus.WarnLevel
	case 4:
		level = logrus.InfoLevel
	case 5:
		level = logrus.DebugLevel
	case 6:
		level = logrus.TraceLevel
	}
	return level
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
