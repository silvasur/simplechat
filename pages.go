package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"path"
)

var (
	TplHome, TplChat *template.Template
)

type ChatpageData struct {
	Websock, Roomname string
}

func PrepTemplates() {
	TplHome = template.Must(template.ParseFiles(path.Join(*tplpath, "home.html")))
	TplChat = template.Must(template.ParseFiles(path.Join(*tplpath, "chat.html")))
}

func Home(rw http.ResponseWriter, req *http.Request) {
	TplHome.Execute(rw, nil) // TODO: Should we log the error?
}

func Chatpage(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	TplChat.Execute(rw, ChatpageData{"ws://" + req.Host + req.URL.Path + "socket", vars["chatroom"]})
}
