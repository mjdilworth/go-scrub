package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	*http.Server
}

// NewServer creates and configures a server serving all application routes.
func NewServer(listenAddr string) (*Server, error) {

	api := newAPI()

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: api,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	return &Server{srv}, nil

}

// Routing
func newAPI() *http.ServeMux {

	mux := http.NewServeMux()

	//example to pass a variable to the handler. in this case a time format string
	th := TimeHandler(time.RFC1123)
	mux.Handle("/time", th)

	mux.HandleFunc("/health/", Health)
	mux.HandleFunc("/", Root)
	mux.HandleFunc("/secret/", Auth)
	mux.HandleFunc("/spacepeeps/", Spacepeeps)
	mux.HandleFunc("/help", Help)
	mux.HandleFunc("/scrub",Scrub)

	return mux
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) Start() {
	fmt.Println("Starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on %s\n", srv.Addr)
			log.Printf("%+v", err)
		}
	}()
	fmt.Println("Server is ready to handle requests")
	srv.gracefulShutdown()
}

// Start runs ListenAndServeTLS on the http.Server with graceful shutdown
func (srv *Server) StartTLS(certFile, keyFile string) {
	fmt.Println("Starting HTTPS server...")

	go func() {
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on %s\n", srv.Addr)
			log.Printf("%+v", err)
			os.Exit(-1)
		}
	}()
	fmt.Println("HTTPS Server is ready to handle requests")

	srv.gracefulShutdown()
}
func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	fmt.Printf("Server is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Could not gracefuly shutdown the server", err)
	}
	fmt.Println("Server stopped")
}
