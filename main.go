package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/gin-gonic/gin"
)

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/signup/", signup)
	http.HandleFunc("/signin/", login)
	http.HandleFunc("/club/", club_option)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
