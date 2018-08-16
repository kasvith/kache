package klogs

import (
	"bytes"
	"fmt"
	"github.com/kasvith/kache/internal/config"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var Logger *logrus.Entry

func InitLoggers(config config.AppConfig) {
	var logrusLogger = logrus.New()

	if config.Debug == true {
		logrusLogger.SetLevel(logrus.DebugLevel)
	} else if config.Verbose == true {
		logrusLogger.SetLevel(logrus.InfoLevel)
	} else {
		logrusLogger.SetLevel(logrus.WarnLevel)
	}

	logrusLogger.Formatter = &KacheFormatter{}
	Logger = logrusLogger.WithFields(logrus.Fields{"pid": os.Getpid()})

	// if we dont want logging, just discard all to a null device
	if config.Logging == false {
		logrusLogger.Out = ioutil.Discard
	}

	if config.Logging && config.Logfile != "" {
		fp, err := os.OpenFile(config.Logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			Logger.Warningf("%s cannot be opened, continue with stderr", config.Logfile)
		} else {
			multi := io.MultiWriter(os.Stderr, fp)
			logrusLogger.Out = multi
		}
	}

	// logrusLogger will output to stderr by default
	// logrusLogger will log the pid of the process

	Logger.Info("OK")
}

type KacheFormatter struct {
}

func (KacheFormatter) Format(e *logrus.Entry) ([]byte, error) {
	buffer := bytes.Buffer{}
	str := fmt.Sprintf("[%s] %s(%d): %s\n", strings.ToUpper(e.Level.String()[0:4]), e.Time.Format("2006-01-02 15:04:05"), e.Data["pid"], e.Message)
	buffer.WriteString(str)

	return buffer.Bytes(), nil
}

func PrintErrorAndExit(err error, exit int) {
	if os.Getenv("ENV") == "DEBUG" {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(exit)
}
