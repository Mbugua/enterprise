package logging

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"time"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type LogData struct {
	Time    time.Time
	Message string
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request
		end := time.Now()
		latency := end.Sub(start)

		logData := LogData{
			Time:    time.Now(),
			Message: fmt.Sprintf("[%s] %s %s %s", r.Method, r.RequestURI, r.Proto, latency),
		}

		// Write to log file
		writeToFile(logData)
	})
}

func writeToFile(logData LogData) {
	fileName := getLogFileName()
	file := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    getMaxFileSize(),
		MaxBackups: 3,    // Maximum number of old log files to retain
		Compress:   true, // Compress the rotated log files using gzip
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	logstashHook := getLogstashHook()
	if logstashHook != nil {
		logrus.AddHook(logstashHook)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			file := path.Base(f.File)
			return fmt.Sprintf("%s:%d", file, f.Line), ""
		},
	})

	logrus.WithFields(logrus.Fields{
		"level": "info",
		"time":  logData.Time.Format("2006-01-02 15:04:05"),
	}).Info(logData.Message)
}

func getLogFileName() string {
	logDir := godotenv.load("LOG_DIR")
	if logDir == "" {
		logDir = "../../logs/my_app"
	}

	logFile := godotenv.load("LOG_FILE_NAME")
	if logFile == "" {
		logFile = "myapp.log"
	}

	return fmt.Sprintf("%s/%s", logDir, logFile)
}

func getMaxFileSize() int {
	maxSizeStr := godotenv.load("MAX_LOG_FILE_SIZE_MB")
	if maxSizeStr == "" {
		return 256
	}

	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		return 256
	}

	return maxSize
}
func getLogstashHook() logrus.Hook {
	logstashURL := godotenv.load("LOGSTASH_URL")
	if logstashURL == "" {
		return nil
	}

	hook, err := logrustash.NewHook("tcp", logstashURL, "my_app ")
	if err != nil {
		fmt.Println("Failed to create Logstash hook:", err)
		return nil
	}

	return hook
}
