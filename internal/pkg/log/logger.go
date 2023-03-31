package log

import (
	domLogger "article-dispatcher/internal/domain/adaptors/logger"
	"fmt"
	"log"
)

const (
	FATAL string = `FATAL`
	ERROR string = `ERROR`
	WARN  string = `WARN`
	INFO  string = `INFO`
	DEBUG string = `DEBUG`
	TRACE string = `TRACE`
)

// LevelMap in-memory map was kept to track the levels of the logs
// string type was used to increase the visibility of the log type
var LevelMap = map[string]int{
	FATAL: 6,
	ERROR: 5,
	WARN:  4,
	INFO:  3,
	DEBUG: 2,
	TRACE: 1,
}

type logger struct {
	Level string
}

// NewLogger create a new logger with several levels of logging
// FATAL being the highest level and TRACE being the lowest
func NewLogger(level string) (domLogger.Logger, error) {
	_, ok := LevelMap[level]
	if !ok {
		return nil, fmt.Errorf("invalid log level received [%s]", level)
	}

	l := &logger{Level: level}
	return l, nil
}

// Fatal only fatal logs will be logged
func (l *logger) Fatal(message string) {
	log.Fatalln("[FATAL]: ", message)
}

// Error both fatal and error logs will be logged
func (l *logger) Error(message string) {
	if LevelMap[l.Level] <= LevelMap[ERROR] {
		log.Println("[ERROR]: ", message)
	}
}

// Warn fatal,error and warn logs will be logged
func (l *logger) Warn(message string) {
	if LevelMap[l.Level] <= LevelMap[WARN] {
		log.Println("[WARN]: ", message)
	}
}

// Debug fatal,error,warn and debug logs will be logged
func (l *logger) Debug(message string) {
	if LevelMap[l.Level] <= LevelMap[DEBUG] {
		log.Println("[DEBUG]: ", message)
	}
}

// Info fatal,error,warn,debug and info logs will be logged
func (l *logger) Info(message string) {
	if LevelMap[l.Level] <= LevelMap[INFO] {
		log.Println("[INFO]: ", message)
	}
}

// Trace fatal,error,warn,debug,info and trace logs will be logged
func (l *logger) Trace(message string) {
	if LevelMap[l.Level] <= LevelMap[TRACE] {
		log.Println("[TRACE]: ", message)
	}
}
