package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"math/rand"
	"net/http"
	"path"
	"strings"
)

var (
	TplHome, TplChat *template.Template
)

type ChatpageData struct {
	Websock, Roomname string
}

type HomeData struct {
	RandomChat string
}

var randroomAlphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLNOPQRSTUVWXYZ")

const randroomLen = 10

func randomRoom() string {
	name := make([]rune, randroomLen)
	for i := 0; i < randroomLen; i++ {
		name[i] = randroomAlphabet[rand.Intn(len(randroomAlphabet))]
	}
	return string(name)
}

func PrepTemplates() {
	TplHome = template.Must(template.ParseFiles(path.Join(*tplpath, "home.html")))
	TplChat = template.Must(template.ParseFiles(path.Join(*tplpath, "chat.html")))
}

func Home(rw http.ResponseWriter, req *http.Request) {
	TplHome.Execute(rw, HomeData{randomRoom()})
}

func Chatpage(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	wsproto := "ws://"
	if *usetls {
		wsproto = "wss://"
	}
	TplChat.Execute(rw, ChatpageData{wsproto + req.Host + strings.Replace(req.URL.Path+"/socket", "//", "/", -1), vars["chatroom"]})
}
