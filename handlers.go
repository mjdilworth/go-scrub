package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mjdilworth/go-scrub/httpreq"
)

// Req is the request query struct.
// 5105 1051 0510 5100 example master debit card format
// curl 'http://localhost:8080?timestamp=1437743020&limit=10&card=5105%201051%200510%205100&fields=foo,bar,badger'
type Req struct {
	Fields    []string
	Limit     int
	Card      string
	Timestamp time.Time
}

func Scrub(w http.ResponseWriter, r *http.Request) {
	//what is the query string i need to scrub
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := &Req{}
	if err := httpreq.NewParsingMap().
		Add("limit", httpreq.ToInt, &data.Limit).
		Add("card", httpreq.ToString, &data.Card).
		Add("fields", httpreq.ToCommaList, &data.Fields).
		Add("timestamp", httpreq.ToTSTime, &data.Timestamp).
		Parse(r.Form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	_ = json.NewEncoder(w).Encode(data)

}


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

func TimeHandler(format string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
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

// This functions goes off to the net to find peopel currenlty in sapce - i can inject delay in this
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
