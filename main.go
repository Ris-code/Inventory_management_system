package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	// "github.com/gin-gonic/gin"
)

func main() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/signup/", signup)
	r.HandleFunc("/signin/", login)
	r.HandleFunc("/club/", club_option)
	r.HandleFunc("/inventory/{ID}/", inventory)

	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
