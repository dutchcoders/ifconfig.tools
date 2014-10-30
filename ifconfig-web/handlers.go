/*
The MIT License (MIT)

Copyright (c) 2014 DutchCoders [https://github.com/dutchcoders/]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	html_template "html/template"
	_ "math/rand"
	"net"
	"net/http"
	text_template "text/template"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Approaching Neutral Zone, all systems normal and functioning.")
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	action := vars["action"]

	switch action {
	case "hostname":
		addr := getIpAddress(r)

		names, err := net.LookupAddr(addr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, name := range names {
			fmt.Fprintf(w, "%s\n", name)
		}
	case "ip":
		fmt.Fprintf(w, "%s\n", getIpAddress(r))
	case "useragent":
		fmt.Fprintf(w, "%s\n", r.UserAgent())
	default:
		fmt.Fprintf(w, "Unknown action %s.\n", action)

	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	addr := getIpAddress(r)

	data := struct {
		IpAddress string
		Hostname  string
		UserAgent string
	}{
		addr,
		"",
		r.UserAgent(),
	}

	names, err := net.LookupAddr(addr)
	if err == nil {
		data.Hostname = names[0]
	}

	// vars := mux.Vars(r)
	if acceptsHtml(r.Header) {
		tmpl, err := html_template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		tmpl, err := text_template.ParseFiles("static/index.txt")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

func RedirectHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ipAddrFromRemoteAddr(r.Host) != "ifconfig.tools" && ipAddrFromRemoteAddr(r.Host) != "127.0.0.1" && r.URL.Path != "/health.html" {
			http.Redirect(w, r, "https://ifconfig.tools"+r.RequestURI, http.StatusMovedPermanently)
			return
		}

		if acceptsHtml(r.Header) && ipAddrFromRemoteAddr(r.Host) == "ifconfig.tools" && r.Header.Get("X-Forwarded-Proto") != "https" && r.Method == "GET" {
			http.Redirect(w, r, "https://ifconfig.tools"+r.RequestURI, http.StatusMovedPermanently)
			return
		}

		h.ServeHTTP(w, r)
	}
}

// Create a log handler for every request it receives.
func LoveHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-made-with", "<3 by DutchCoders")
		w.Header().Set("x-served-by", "Proudly served by DutchCoders")
		w.Header().Set("Server", "ifconfig.tools HTTP Server 1.0")
		h.ServeHTTP(w, r)
	}
}
