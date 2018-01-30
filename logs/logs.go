package logs

import (
	"os"
	"log"
)

const (
	debug = 0
	info  = 1
	error = 2
	fatal = 3
)

const (
	debugTag = "[debug] "
	infoTag  = "[info ] "
	errorTag = "[error] "
	fatalTag = "[fatal] "
)

type Logger struct {
	file      *os.File
	logLogger *log.Logger
}

func New(filePath string) *Logger {

	//var logger = new(Logger)
	var logger = &Logger{}

	if filePath != "" {
		//file, err := os.Create(filePath)
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil
		}

		logger.file = file
		logger.logLogger = log.New(file, "", log.LstdFlags)
	} else {
		logger.logLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return logger
}

func (l *Logger) format(data string, flag int) {
	//l.logLogger.SetFlags(log.LstdFlags)
	switch flag {
	case 0:
		l.logLogger.SetPrefix(debugTag)
	case 1:
		l.logLogger.SetPrefix(infoTag)
	case 2:
		l.logLogger.SetPrefix(errorTag)
	case 3:
		l.logLogger.SetPrefix(fatalTag)
	default:
		l.logLogger.SetPrefix(infoTag)
	}
	l.logLogger.Println(data)

	if flag == 3 {
		os.Exit(1)
	}
}

func (l *Logger) Info(data string) {
	l.format(data, info)
}

func (l *Logger) Debug(data string) {
	l.format(data, debug)
}

func (l *Logger) Error(data string) {
	l.format(data, error)
}

func (l *Logger) Fatal(data string) {
	l.format(data, fatal)
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
