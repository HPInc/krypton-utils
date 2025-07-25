package logging

import (
	"log"
	"net/http"
)

// level declaration
type Level uint32

const (
	Info Level = iota
	Verbose
)

type Log struct {
	docProvider DocProvider
	level       Level
	logger      *log.Logger
	plainLogger *log.Logger
}

var logger Log

func InitLogger(l Level) *Log {
	logger = Log{
		level:       l,
		logger:      log.Default(),
		plainLogger: log.New(log.Default().Writer(), "", 0),
	}
	return &logger
}

// provide a log util
func GetLogger() *Log {
	return &logger
}

// parse and set doc type
func (l *Log) SetDocType(dt string) {
	l.docProvider = NewDocProvider(dt)
}

func (l *Log) SetLevel(value Level) {
	l.level = value
}

func (l *Log) IsVerbose() bool {
	return l.level >= Verbose
}

func (l *Log) IsInfo() bool {
	return l.level >= Info
}

func (l *Log) Debug(v ...any) {
	if l.IsVerbose() {
		l.logger.Println(v...)
	}
}

func (l *Log) Debugf(format string, v ...any) {
	if l.IsVerbose() {
		l.logger.Printf(format, v...)
	}
}

func (l *Log) Info(v ...any) {
	if l.IsInfo() {
		l.logger.Println(v...)
	}
}

func (l *Log) FatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (l *Log) Fatal(v ...any) {
	log.Fatal(v...)
}

func (l *Log) Fatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}

func (l *Log) Errorf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

// this function is here for a possible decoration for errors
// right now, its the same as Println
func (l *Log) Error(v ...any) {
	log.Println(v...)
}

func (l *Log) Println(v ...any) {
	log.Println(v...)
}

func (l *Log) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func (l *Log) GetPlainLogger() *log.Logger {
	return l.plainLogger
}

// some specialized logging for http requests and responses
func (l *Log) HttpRequest(req *http.Request) {
	l.docProvider.HttpRequest(req)
}

func (l *Log) HttpResponse(resp *http.Response, body []byte) {
	l.docProvider.HttpResponse(resp, body)
}
