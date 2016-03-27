package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var debug = os.Getenv("DEBUG")
func logDebug(format string, args ...interface{}) {
	if debug != "" {
		if len(args) == 0 {
			log.Print(format)
		} else {
			log.Printf(format, args)
		}
	}
}

func fallback(w http.ResponseWriter, r *http.Request, reason string) {
	location := "http://redirect.name/"
	if reason != "" {
		location = fmt.Sprintf("%s#reason=%s", location, url.QueryEscape(reason))
	}
	http.Redirect(w, r, location, 302)
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

	logDebug("Checking URL: %s", r.URL.String())

	var catch_alls []string
	logDebug("Begin record search.")
	for _, record := range txt {
		logDebug("Inspecting: %s", record)
		config := Parse(record)
		if strings.TrimSpace(config.From) == "" {
			logDebug("Saving catch all.")
			catch_alls = append(catch_alls, record)
			_ = catch_alls
			continue
		}
		redirect := Translate(r.URL.String(), Parse(record))
		if redirect != nil {
			http.Redirect(w, r, redirect.Location, redirect.Status)
			logDebug("Found matching record.")
			return
		}
	}
	logDebug("End record search")

	logDebug("Begin catch alls")
	for _, record := range catch_alls {
		logDebug("Inspecting: %s", record)
		redirect := Translate(r.URL.String(), Parse(record))
		if redirect != nil {
			http.Redirect(w, r, redirect.Location, redirect.Status)
			logDebug("Found catch all.")
			return
		}
	}
	logDebug("No match found.")

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
