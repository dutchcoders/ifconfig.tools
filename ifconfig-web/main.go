package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var config struct {
}

func init() {
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()

	r.PathPrefix("/scripts/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/styles/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/images/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/fonts/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/ico/").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/robots.txt").Handler(http.FileServer(http.Dir("./static/")))
	r.HandleFunc("/health.html", healthHandler).Methods("GET")
	r.HandleFunc("/{action}", actionHandler).Methods("GET")
	r.HandleFunc("/", viewHandler).Methods("GET")

	port := flag.String("port", "8080", "port number, default: 8080")
	flag.Parse()

	log.Printf("Ifconfig server started. :%v", *port)
	log.Printf("---------------------------")

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: handlers.PanicHandler(RedirectHandler(LoveHandler(handlers.LogHandler(r, handlers.NewLogOptions(log.Printf, "_default_")))), nil),
	}

	log.Panic(s.ListenAndServe())
	log.Printf("Server stopped.")

}
