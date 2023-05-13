package log

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/hr3lxphr6j/bililive-go/src/instance"
	"github.com/hr3lxphr6j/bililive-go/src/interfaces"
)

func New(ctx context.Context) *interfaces.Logger {
	inst := instance.GetInstance(ctx)
	logLevel := logrus.InfoLevel
	if inst.Config.Debug {
		logLevel = logrus.DebugLevel
	}
	config := inst.Config
	writers := []io.Writer{os.Stderr}
	outputFolder := config.Log.OutPutFolder
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		log.Fatalf("err: \"%s\", Failed to determine log output folder: %s", err, outputFolder)
	} else {
		if config.Log.SaveEveryLog {
			runID := time.Now().Format("run-2006-01-02-15-04-05")
			logLocation := filepath.Join(outputFolder, runID+".log")
			logFile, err := os.OpenFile(logLocation, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Failed to open log file %s for output: %s", logLocation, err)
			} else {
				writers = append(writers, logFile)
			}
		}
		if config.Log.SaveLastLog {
			logLocation := filepath.Join(outputFolder, "bililive-go.log")
			logFile, err := os.OpenFile(logLocation, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalf("Failed to open default log file %s for output: %s", logLocation, err)
			} else {
				writers = append(writers, logFile)
			}
		}
	}
	logger := &interfaces.Logger{Logger: &logrus.Logger{
		Out: io.MultiWriter(writers...),
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
		Hooks: make(logrus.LevelHooks),
		Level: logLevel,
	}}

	inst.Logger = logger

	return logger
}
