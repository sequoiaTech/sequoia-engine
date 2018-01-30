package logs

import (
	"os"
	"log"
)

const (
	debugLevel = 0
	infoLevel  = 1
	errorLevel = 2
	fatalLevel = 3
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

func New(filePath string) (*Logger, error) {

	//var logger = new(Logger)
	var logger = &Logger{}

	if filePath != "" {
		//file, err := os.Create(filePath)
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		logger.file = file
		logger.logLogger = log.New(file, "", log.LstdFlags)
	} else {
		logger.logLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return logger, nil
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
	l.format(data, infoLevel)
}

func (l *Logger) Debug(data string) {
	l.format(data, debugLevel)
}

func (l *Logger) Error(data string) {
	l.format(data, errorLevel)
}

func (l *Logger) Fatal(data string) {
	l.format(data, fatalLevel)
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
