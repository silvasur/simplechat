package main

import (
	"flag"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	laddr      = flag.String("laddr", ":8080", "Listen on this address")
	tplpath    = flag.String("tplpath", "tpls", "Path to templates")
	staticpath = flag.String("staticpath", "static", "Path to static page elements")
)

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticpath))))
	r.HandleFunc("/chat/{chatroom}/", Chatpage)
	r.HandleFunc("/chat/{chatroom}/socket", AcceptWebSock)
	http.Handle("/", r)
	http.ListenAndServe(*laddr, nil)
}
