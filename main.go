package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/silvasur/simplechat/chat"
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
	usetls     = flag.Bool("tls", false, "Should TLS be used?")
	certfile   = flag.String("tlscert", "", "TLS certificate file")
	keyfile    = flag.String("tlskey", "", "TLS key file")
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

	var listenServe func() error
	if *usetls {
		listenServe = func() error {
			return http.ListenAndServeTLS(*laddr, *certfile, *keyfile, nil)
		}
	} else {
		listenServe = func() error {
			return http.ListenAndServe(*laddr, nil)
		}
	}

	if err := listenServe(); err != nil {
		log.Fatal(err)
	}
}
