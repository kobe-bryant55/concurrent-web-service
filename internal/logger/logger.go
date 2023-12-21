package logger

import (
	"log"
	"os"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func Init() {
	c, err := os.OpenFile("critical.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	i, err := os.OpenFile("info.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	wr, err := os.OpenFile("warning.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(i, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(wr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(c, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
