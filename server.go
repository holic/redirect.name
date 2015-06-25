package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func fallback(w http.ResponseWriter, r *http.Request, reason string) {
	location := "http://redirect.name/"
	if reason != "" {
		location = fmt.Sprintf("%s#reason=%s", location, url.QueryEscape(reason))
	}
	http.Redirect(w, r, location, 302)
}

func handler(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		fallback(w, r, fmt.Sprintf("Could not split host (%v)", err))
		return
	}

	hostname := fmt.Sprintf("_redirect.%s", host)
	txt, err := net.LookupTXT(hostname)
	if err != nil {
		fallback(w, r, fmt.Sprintf("Could not resolve hostname (%v)", err))
		return
	}

	for _, record := range txt {
		redirect := Translate(r.URL.String(), Parse(record))
		if redirect != nil {
			http.Redirect(w, r, redirect.Location, redirect.Status)
			return
		}
	}

	fallback(w, r, "No paths matched")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	log.Printf("Listening on http://127.0.0.1:%s", port)
	log.Fatal(srv.ListenAndServe())
}
