package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// set up my log levels
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {

	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	//keep main minimal
	var err error

	//flags
	serverPort := flag.String("port", ":8080", "specify the port the server listens on")
	certFile := flag.String("certfile", "self-signed-cert/cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "self-signed-cert/key.pem", "key PEM file")
	tls := flag.Bool("tls", false, "run HTTPS")
	//username := *flag.String("user", "foobar", "username used for basic auth")
	//password := *flag.String("pass", "x!gH1as", "passowrd used for basic auth")

	flag.Parse()

	//start log goroutine
	command := make(chan string)
	go logging(command)
	command <- "Pause"

	//server *http.Server
	server, err := NewServer(*serverPort, command)
	if err != nil {
		panic(err)
	}

	if !*tls {
		server.Start()
	} else {
		server.StartTLS(*certFile, *keyFile)
	}

	fmt.Println("ended")

}

func logging(command <-chan string) {
	var status = "Play"
	var logLevel = "info"
	count := 0
	for {
		select {
		case cmd := <-command:
			fmt.Println(cmd)
			switch cmd {
			case "stop":
				return
			case "pause":
				status = "pause"
			case "info":
				logLevel = "info"
			case "warn":
				logLevel = "warn"
			case "error":
				logLevel = "error"
			default:
				status = "play"
			}
		case <-time.After(1 * time.Second):
			if status == "play" {
				logwork(count, logLevel)
				count = count + 1
			}
		}
	}
}

func logwork(i int, llevel string) {
	//time.Sleep(250 * time.Millisecond)
	switch llevel {
	case "info":
		InfoLogger.Printf("log something %d", i)
	case "warn":
		WarningLogger.Printf("log something %d", i)
	case "error":
		ErrorLogger.Printf("log something %d", i)
	default:
		InfoLogger.Printf("log something %d", i)
	}
}

