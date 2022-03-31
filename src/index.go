package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"tracker/internal/app"
	"tracker/internal/app/logger"
	"tracker/internal/app/server"
)

func main() {

	var (
		port int
		err  error
	)

	startServer := flag.Bool("s", true, "Start server")
	startParser := flag.Bool("p", false, "Start parsing")

	flag.Parse()

	_port, ok := os.LookupEnv("port")
	if !ok {
		log.Println("Введите номер порта или настройте его в окружении")
		return
	}

	port, err = strconv.Atoi(_port)

	if err != nil {
		log.Println(err)
		return
	}

	if *startServer {
		server.Start(port)
	}

	if *startParser {
		files := logger.GetFiles()

		app.TestParser(files)
	}

}
