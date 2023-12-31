package logger

import (
	"os"

	"test/conf"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger_instance = log.New()

var filepath = "./log/logrus.log"

func SetPath(path string) {

	filepath = path
}

func Init() {

	//getting the credentials in the conf file
	logLevel := viper.Get(conf.LogLevel)

	//setting the format of the logs to be a JSON one
	logger_instance.SetFormatter(&log.JSONFormatter{
		DataKey:     "Assessment",
		PrettyPrint: true,
	})

	//getting the log level set in the configuration file
	LogLevel, err := log.ParseLevel(logLevel.(string))
	//If the log level in conf file can't be parsed, log level should be the default info level
	if err != nil {
		logLevel = log.InfoLevel
	}
	//setting the log level
	logger_instance.SetLevel(LogLevel)
	//If we want to throw logs into a local file

	logger_instance.SetOutput(os.Stdout)
	//setting it to a file writer
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if err == nil {
		logger_instance.Out = file
	} else {
		logger_instance.Info("Failed to log to file, using default stderr")
	}
	logger_instance.Info("Logrus has been initiated")

}

func Instance() *log.Logger {

	return logger_instance
}
