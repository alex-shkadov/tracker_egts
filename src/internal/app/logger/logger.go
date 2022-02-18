package logger

import (
	"log"
	"os"
	"regexp"
	"time"
)

const DIRECTORY = "logs"

func logFileName() string {
	return time.Now().Format("2006-01-02_15") + ":00.log"
}

func logFilePath() string {
	return DIRECTORY + "/" + logFileName()
}

func GetFiles() []string {
	regex := `^\d{4}\-\d{2}\-\d{2}_\d{2}\:00\.log$`
	files, err := os.ReadDir(DIRECTORY)
	if err != nil {
		log.Fatalln(err)
	}

	result := []string{}

	for _, file := range files {
		matched, err := regexp.Match(regex, []byte(file.Name()))

		if err != nil {
			log.Fatalln(err)
		}

		if matched || true {
			result = append(result, DIRECTORY+"/"+file.Name())
		}
	}

	return result
}

func LogETGSConnectionData(bytes []byte, in bool) {
	logFile := logFilePath()

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	prefix := ">> "
	if !in {
		prefix = "<< "
	}

	log.Println("LOG TO FILE", logFile, "DATA:", bytes)
	dataLog := log.New(file, prefix, log.Ldate|log.Ltime)

	dataLog.Println(bytes)
}
