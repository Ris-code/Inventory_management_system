package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/gin-gonic/gin"
)

func main() {
	// db, err := sql.Open("mysql", "root:rishav@2003@tcp(127.0.0.1:3306)/PERSON")
	// if err != nil {
	//     panic(err.Error())
	// }

	// defer db.Close()
	// fmt.Println("Success!")

	// insert, err := db.Query("INSERT INTO customers VALUES ( '6', 'Risu', '7890567890')")
	// if err != nil {
	//     panic(err.Error())
	// }

	// // _, err := db.Exec("CREATE TABLE IF NOT EXISTS mytable (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, some_text TEXT NOT NULL)")
	// // if err != nil {
	// //     panic(err)
	// // }

	// defer insert.Close()
	// fmt.Println("Success!")/

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/signup/", signup)
	http.HandleFunc("/signin/", login)
	initDB()
	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
