// package main
package handler

import (
	// "log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

// func main() {

// 	r := mux.NewRouter()

// 	fs := http.FileServer(http.Dir("static"))
// 	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
// 	r.HandleFunc("/", home_before_login)

// 	// student
// 	r.HandleFunc("/home/", home_after_login)
// 	r.HandleFunc("/signup/", signup)
// 	r.HandleFunc("/signin/", login)
// 	r.HandleFunc("/club/", club_option)
// 	r.HandleFunc("/inventory/{ID}/", inventory)
// 	r.HandleFunc("/cart/", cart)
// 	r.HandleFunc("/update/", update)
// 	r.HandleFunc("/thank/", thank_page)
// 	r.HandleFunc("/borrowlist/", borrow_list)
// 	r.HandleFunc("/deleteItems/", delete_items)

// 	// club 
// 	r.HandleFunc("/clublogin/", coordinator_login)
// 	r.HandleFunc("/clubhome/", club_home)
// 	r.HandleFunc("/additem/", add_inventory)
// 	r.HandleFunc("/updateinfo/", update_info)
// 	r.HandleFunc("/clubborrowlist/", club_borrow_list)
// 	r.HandleFunc("/inventorylist/", inventorylist)
// 	r.HandleFunc("/inventorylist/edit/", edit_inventory)
// 	r.HandleFunc("/inventorylist/delete/", delete_inventory)
	

// 	// Enable CORS for all routes
// 	corsHandler := handlers.CORS(
// 		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
// 		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
// 		handlers.AllowedOrigins([]string{"*"}),
// 	)

// 	// Apply CORS middleware to the router
// 	http.Handle("/", corsHandler(r))

// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}

// }

func Handler(w http.ResponseWriter, req *http.Request) {
	router := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Define routes
	router.HandleFunc("/", home_before_login)
	router.HandleFunc("/home/", home_after_login)
	router.HandleFunc("/signup/", signup)
	router.HandleFunc("/signin/", login)
	router.HandleFunc("/club/", club_option)
	router.HandleFunc("/inventory/{ID}/", inventory)
	router.HandleFunc("/cart/", cart)
	router.HandleFunc("/update/", update)
	router.HandleFunc("/thank/", thank_page)
	router.HandleFunc("/borrowlist/", borrow_list)
	router.HandleFunc("/deleteItems/", delete_items)

	// Club routes
	router.HandleFunc("/clublogin/", coordinator_login)
	router.HandleFunc("/clubhome/", club_home)
	router.HandleFunc("/additem/", add_inventory)
	router.HandleFunc("/updateinfo/", update_info)
	router.HandleFunc("/clubborrowlist/", club_borrow_list)
	router.HandleFunc("/inventorylist/", inventorylist)
	router.HandleFunc("/inventorylist/edit/", edit_inventory)
	router.HandleFunc("/inventorylist/delete/", delete_inventory)

	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
	)

	// Apply CORS middleware to the router
	corsHandler(router).ServeHTTP(w, req)
}