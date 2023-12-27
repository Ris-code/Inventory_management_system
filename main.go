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
	r.HandleFunc("/borrowlist/", borrow_list)
	r.HandleFunc("/deleteItems/", delete_items)

	// club 
	r.HandleFunc("/clublogin/", coordinator_login)
	r.HandleFunc("/clubhome/", club_home)
	r.HandleFunc("/additem/", add_inventory)
	r.HandleFunc("/updateinfo/", update_info)
	r.HandleFunc("/clubborrowlist/", club_borrow_list)
	r.HandleFunc("/inventorylist/", inventorylist)
	r.HandleFunc("/inventorylist/edit/", edit_inventory)
	r.HandleFunc("/inventorylist/delete/", delete_inventory)
	

	// Enable CORS for all routes
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
	)

	// Apply CORS middleware to the router
	http.Handle("/", corsHandler(r))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
