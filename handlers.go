package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var logging bool = false

// a channel to tell it to stop
var stoplogchan = make(chan struct{})

// a channel to signal that it's stopped
var stoppedlogchan = make(chan struct{})

// placeholder
func verifyUserPass(user string, pass string) bool {

	usernameHash := sha256.Sum256([]byte(user))
	passwordHash := sha256.Sum256([]byte(pass))
	expectedUsernameHash := sha256.Sum256([]byte(user))
	expectedPasswordHash := sha256.Sum256([]byte(pass))

	usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
	passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

	if usernameMatch && passwordMatch {
		return true
	} else {
		return false
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "I am healthy"}`))
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "HTTP Served by GO"}`))
}

func Help(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "commands to use: play, pause, stop, info, warn, error "}`))
}

func Auth(w http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()
	if ok && verifyUserPass(user, pass) {
		w.Write([]byte(`{"message": "You get to see the secret"}`))
		//fmt.Fprintf(w, "You get to see the secret\n")
	} else {
		// i should redirect to login page
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// Used for matching variables in a request URL.
//var reResVars = regexp.MustCompile(`\\\{[^{}]+\\\}`)

// Log handler to sopt and start
func LogOrig(w http.ResponseWriter, r *http.Request) {

	//this retrives the last elementin URI
	request := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	//request := r.URL.Path

	if request == "start" {
		if logging {
			return
		}

		//start the logging go routine
		logging = true
		go func() { // work in background
			// close the stoppedchan when this func
			// exits
			defer close(stoppedlogchan)
			//defer stoppedlogchan = make(chan struct{})
			// TODO: do setup work
			defer func() {
				// TODO: do teardown work
				fmt.Println("Graceful handler exit using stoplogchan")
			}()
			for {
				select {
				// TODO: do a bit of the work
				case <-time.After(1 * time.Second):
					fmt.Println("log something")
				case <-stoplogchan:
					fmt.Println("stopping")
					// stop
					return
				}
			}
		}()

	} else if request == "stop" {
		if !logging {
			return
		}
		//stop the logging go routine
		log.Println("stopping...")
		close(stoplogchan) // tell it to stop
		<-stoppedlogchan   // wait for it to have stopped
		logging = false
		//recrreact channels
		stoplogchan = make(chan struct{})
		log.Println("Stopped.")

	} else {
		//do nothing
		fmt.Println("unknown request")
	}

	w.WriteHeader(http.StatusOK)
}

func TimeHandler(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}

func LogHandler(command chan<- string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//this retrives the last elementin URI
		request := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		switch request {
		case "play":
			command <- "play"
			log.Println("Playing logs")
		case "pause":
			command <- "pause"
			log.Println("Pausing logs")
		case "stop":
			command <- "stop"
			log.Println("Stopping")
		case "info":
			command <- "info"
			log.Println("Setting to INFO")
		case "warn":
			command <- "warn"
			log.Println("Setting to WARN")
		case "error":
			command <- "error"
			log.Println("Setting to ERROR")
		default:
			log.Printf("Unkown command %s : send help for ... help", request)
		}
	}
	return http.HandlerFunc(fn)
}

type person struct {
	Name string `json:"name"`
}
type people struct {
	Number int      `json:"number"`
	Person []person `json:"people"`
}

//This functions goes off to the net to find peopel currenlty in sapce - i can inject delay in this
func Spacepeeps(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	//w.Write([]byte(`{"message": "I am healthy"}`))
	apiURL := "http://api.open-notify.org/astros.json"

	people, err := getAstros(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)

	w.Header().Set("AtTheEnd1", "Mikes value 1")
	io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
	
	w.Header().Set("AtTheEnd2", "Mikes value 2")

	
	sout := fmt.Sprintf("%d people found in space.\n", people.Number)
	io.WriteString(w, sout)
	for _, p := range people.Person {

		sout = fmt.Sprintf("Hola to: %s\n", p.Name)
		io.WriteString(w, sout)
	}

}
func getAstros(apiURL string) (people, error) {
	p := people{}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return p, err
	}
	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return p, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return p, err
	}

	if err := json.Unmarshal(body, &p); err != nil {
		return p, err
	}

	return p, nil
}
