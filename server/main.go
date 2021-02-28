package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":9000", router)
}

// Message is used to send a title and body as error message to client
type Message struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Token   string `json:"token"`
	Timeout int64  `json:"timeout"`
}
