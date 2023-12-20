package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func initDB() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("mysql", "root:rishav@2003@tcp(127.0.0.1:3306)/club")
	if err != nil {
		panic(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/signin.html")
		t.Execute(w, nil)
	} else {
		// db, err = sql.Open("mysql", "root:rishav@2003@tcp(127.0.0.1:3306)/club")
		r.ParseForm()
		// logic part of log in
		var login_username = r.FormValue("login-username")
		var login_password = r.FormValue("login-password")

		fmt.Println("username:", login_username)
		fmt.Println("password:", login_password)

		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Get the existing entry present in the database for the given username
		result := db.QueryRow("select password from student where username=$1", login_username)
		if err != nil {
			// If there is an issue with the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Compare the stored hashed password, with the hashed version of the password that was received
		hashedPassword := []byte{}
		if err := result.Scan(&hashedPassword); err != nil {
			// If there is an issue with retrieving the hashed password from the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(login_password)); err != nil {
			// If the two passwords don't match, return a 401 status
			w.WriteHeader(http.StatusUnauthorized)
		}
		
		fmt.Println("Record fetched successfully")
		// If we reach this point, that means the users password was correct, and that they are authorized
		// The default 200 status is sent
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of sign up
		var username = r.FormValue("username")
		var password = r.FormValue("password")
		var name = r.FormValue("name")

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

		fmt.Println("username:", username)
		fmt.Println("password:", hashedPassword)
		fmt.Println("name:", name)

		insert, err := db.Prepare("INSERT INTO student (username, name, password) VALUES (?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(username, name, hashedPassword)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Record inserted successfully")
	}
}

