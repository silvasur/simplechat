package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/kch42/simplechat/chat"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var (
	laddr      = flag.String("laddr", ":8080", "Listen on this address")
	tplpath    = flag.String("tplpath", "tpls", "Path to templates")
	staticpath = flag.String("staticpath", "static", "Path to static page elements")
	perroom    = flag.Int("perroom", -1, "Maximum amount of users per room (negative for unlimited)")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	if *perroom < 0 {
		*perroom = math.MaxInt32
	} else if *perroom == 0 {
		log.Fatalln("flag perroom must not be 0")
	}

	PrepTemplates()
	chat.InitRooms(*perroom)

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticpath))))
	r.HandleFunc("/chat/{chatroom:.+}/socket", AcceptWebSock)
	r.HandleFunc("/chat/{chatroom:.+}/", Chatpage)
	r.HandleFunc("/chat/{chatroom:.+}", Chatpage)
	http.Handle("/", r)
	http.ListenAndServe(*laddr, nil)
}
