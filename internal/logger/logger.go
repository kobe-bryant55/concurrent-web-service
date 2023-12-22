package logger

import (
	"log"
	"os"
)

type Logger struct {
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
}

func Init() *Logger {
	c, err := os.OpenFile("logs/critical.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	i, err := os.OpenFile("logs/info.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	wr, err := os.OpenFile("logs/warning.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	return &Logger{
		InfoLog:    log.New(i, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLog: log.New(wr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLog:   log.New(c, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
