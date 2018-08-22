/*
 * MIT License
 *
 * Copyright (c)  2018 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package klogs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/sys"
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

	fields := logrus.Fields{"pid": os.Getpid()}

	switch strings.ToLower(config.LogType) {
	case "json":
		logrusLogger.Formatter = &logrus.JSONFormatter{}
		break
	case "logfmt":
		logrusLogger.Formatter = &logrus.TextFormatter{DisableColors: true, ForceColors: false}
		break
	case "default":
		logrusLogger.Formatter = &kacheFormatter{}
		break
	default:
		logrusLogger.Formatter = &kacheFormatter{}
		logrusLogger.WithFields(fields).Warnf("%s format is unknown, continuing with default", config.LogType)
		break
	}

	Logger = logrusLogger.WithFields(fields)

	// if we dont want logging, just discard all to a null device
	if config.Logging == false {
		logrusLogger.Out = ioutil.Discard
	}

	if config.Logging && config.Logfile != "" {
		// try to create folder path if not exists
		err := sys.AutoCreateSubDirs(config.Logfile)

		// if failed, we can skip for logging to a file, warn user and continue
		if err != nil {
			Logger.Warningf("%s cannot be opened, continue with stderr", config.Logfile)
			return
		}

		// try to create the file
		fp, err := os.OpenFile(config.Logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

		// if failed, skip and warn
		if err != nil {
			Logger.Warningf("%s cannot be opened, continue with stderr", config.Logfile)
			return
		}

		// info about log file
		path, err := filepath.Abs(config.Logfile)

		if err != nil {
			Logger.Errorf("cannot resolve absolute path for %s", config.Logfile)
		} else {
			Logger.Infof("log file is %s", path)
		}

		// use two writers
		multi := io.MultiWriter(os.Stderr, fp)
		logrusLogger.Out = multi
	}
}

type kacheFormatter struct {
}

func (kacheFormatter) Format(e *logrus.Entry) ([]byte, error) {
	buffer := bytes.Buffer{}
	lvl := strings.ToUpper(e.Level.String()[0:4])
	t := e.Time.Format("2006-01-02 15:04:05")
	str := fmt.Sprintf("[%s] %s(%d): %s\n", lvl, t, e.Data["pid"], e.Message)
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
