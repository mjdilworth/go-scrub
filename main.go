package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	//server *http.Server
	server, err := NewServer(*serverPort)
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
