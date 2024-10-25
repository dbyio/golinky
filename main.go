package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"internal/callback"
	"internal/urlshortener"

	"github.com/gorilla/mux"
)

var shortener *urlshortener.Shortener

type shortenParams struct {
	Url string `json:"url"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
	var params shortenParams
	json.NewDecoder(r.Body).Decode(&params)

	shortUrl, err := shortener.CreateShortUrl(params.Url)
	if err != nil {
		log.Fatalln(err)
	}

	body := map[string]string{"shortUrl": shortUrl}
	json_body, _ := json.Marshal(body)
	w.Write(json_body)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	clientIp := r.Header.Get("x-forwarded-for")
	if clientIp == "" {
		re := regexp.MustCompile(`:\d+`)
		clientIp = re.ReplaceAllString(r.RemoteAddr, "") // strip port number
	}

	fmt.Printf("Received request from %s: %s\n", clientIp, id)
	url, err := shortener.UseShortUrl(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	if shortener.HasCallbackUrl() { // don't callback if the url is not configured
		data := &callback.Callback{
			EventType: "click",
			MsgId:     id,
			Link:      url,
			ClickTime: time.Now().Unix(),
			UserAgent: r.Header.Get("User-Agent"),
			ClientIp:  clientIp,
		}
		go callback.Clicked(shortener.CallbackUrl, data)
	}
}

func stats(w http.ResponseWriter, r *http.Request) {
	stats := shortener.Stats()
	body, _ := json.Marshal(stats)
	w.Write(body)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s configuration.yaml", os.Args[0])
	}
	file := os.Args[1]
	shortener = urlshortener.New(file)

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		log.Fatalln("AUTH_TOKEN missing in environment")
	}

	r := mux.NewRouter()
	r.HandleFunc("/shorten", shorten).Methods("POST").HeadersRegexp("X-Auth-Token", authToken)
	r.HandleFunc("/stats", stats).Methods("GET")
	r.HandleFunc("/{id}", redirect).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
