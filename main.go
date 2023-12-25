package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	// "github.com/gin-gonic/gin"
)

func main() {

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", home_before_login)

	// student
	r.HandleFunc("/home/", home_after_login)
	r.HandleFunc("/signup/", signup)
	r.HandleFunc("/signin/", login)
	r.HandleFunc("/club/", club_option)
	r.HandleFunc("/inventory/{ID}/", inventory)
	r.HandleFunc("/cart/", cart)
	r.HandleFunc("/update/", update)
	r.HandleFunc("/thank/", thank_page)

	// club 
	r.HandleFunc("/clublogin/", coordinator_login)
	r.HandleFunc("/clubhome/", club_home)
	r.HandleFunc("/additem/", add_inventory)
	r.HandleFunc("/updateinfo/", update_info)


	 // Enable CORS
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	http.Handle("/", handlers.CORS(headers, methods, origins)(r))
	// http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
