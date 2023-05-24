package logging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type LogData struct {
	Time    time.Time
	Message string
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("unable to load env vars. Err: %s", err)
	}
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

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

	// logstashHook := getLogstashHook()
	// if logstashHook != nil {
	// 	logrus.AddHook(logstashHook)
	// }

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
	appName := os.Getenv("APP_NAME")
	fmt.Printf("app name >>  %s", appName)
	if appName == "" {
		appName = "enterprise"
	}
	logDir := os.Getenv("LOG_DIR")
	fmt.Printf("logs dir >>  %s", logDir)
	if logDir == "" {
		logDir = "../../logs/" + appName + "/"
	}

	logFile := os.Getenv("LOG_FILE_NAME")
	if logFile == "" {
		logFile = "app.log"
	}

	return fmt.Sprintf("%s/%s", logDir, logFile)
}

func getMaxFileSize() int {
	maxSizeStr := os.Getenv("MAX_LOG_FILE_SIZE_MB")
	if maxSizeStr == "" {
		return 256
	}

	maxSize, err := strconv.Atoi(maxSizeStr)
	if err != nil {
		return 256
	}

	return maxSize
}

// func getLogstashHook() logrus.Hook {
// 	logstashURL := os.Getenv("LOGSTASH_URL")
// 	if logstashURL == "" {
// 		return nil
// 	}

// 	hook, err := logrustash.NewHook("tcp", logstashURL, "my_app ")
// 	if err != nil {
// 		fmt.Println("Failed to create Logstash hook:", err)
// 		return nil
// 	}

// 	return hook
// }
