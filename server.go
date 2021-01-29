package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func fallback(w http.ResponseWriter, r *http.Request, reason string) {
	location := os.Getenv("FALLBACK_URL")
	if location == "" {
		location = "http://redirect.name/"
	}
	if reason != "" {
		location = fmt.Sprintf("%s#reason=%s", location, url.QueryEscape(reason))
	}
	http.Redirect(w, r, location, 302)
}

func getRedirect(txt []string, url string) (*Redirect, error) {
	var catchAlls []*Config
	for _, record := range txt {
		config := Parse(record)
		if config.From == "" {
			catchAlls = append(catchAlls, config)
			continue
		}
		redirect := Translate(url, config)
		if redirect != nil {
			return redirect, nil
		}
	}

	var config *Config
	for _, config = range catchAlls {
		redirect := Translate(url, config)
		if redirect != nil {
			return redirect, nil
		}
	}

	return nil, errors.New("No paths matched")
}

func handler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.Host, ":")
	host := parts[0]

	hostname := fmt.Sprintf("_redirect.%s", host)
	txt, err := net.LookupTXT(hostname)
	if err != nil {
		fallback(w, r, fmt.Sprintf("Could not resolve hostname (%v)", err))
		return
	}

	redirect, err := getRedirect(txt, r.URL.String())
	if err != nil {
		fallback(w, r, err.Error())
	} else {
		http.Redirect(w, r, redirect.Location, redirect.Status)
	}
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
