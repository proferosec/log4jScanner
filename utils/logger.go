package utils

import (
    joonix "github.com/joonix/log"
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
)

var Logger *logger

type logger struct {
    Format string
    File *os.File
    Level log.Level
}

func InitLogger() {
    if Logger == nil {
        Logger = &logger{
            File:  nil,
            Level: log.InfoLevel,
            Format: "text",
        }
        Logger.SetFormatter(Logger.Format)
        Logger.SetLevel(Logger.Level)
    }
}

func GetLogger() *logger{
    return Logger
}

func (l *logger) Close() error {
    err := l.File.Close()
    if err != nil {
        return err
    }
    return nil
}

func (l *logger) SetLevel(level log.Level){
    l.Level = level
    log.SetLevel(level)
}
func (l *logger) SetFormatter(format string) {
    l.Format = format
    l.setupFormat()
}
func (l *logger) SetFile(file *os.File) {
    l.File = file
    log.SetOutput(file)
}


func (l *logger) setupFormat() {
    switch strings.ToLower(l.Format) {
    case "text":
        {
            log.SetFormatter(&log.TextFormatter{}) // normal output
        }
    case "json":
        {
            log.SetFormatter(&log.JSONFormatter{}) // simple json output
        }
    case "fluentd":
        {
            log.SetFormatter(joonix.NewFormatter()) //Fluentd compatible
        }
    default:
        log.SetFormatter(joonix.NewFormatter()) //Fluentd compatible
    }
}