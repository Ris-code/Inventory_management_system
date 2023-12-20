package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"os"
	"github.com/joho/godotenv"
)

var db *sql.DB
var err error

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Initialize the database connection
	initDB()
}

func initDB() {
	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")

	// Connect to the MySQL database
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/club", dbUser, dbPass, dbHost))
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
		r.ParseForm()
		// logic part of log in
		var loginUsername = r.FormValue("login-username")
		var loginPassword = r.FormValue("login-password")

		fmt.Println("username:", loginUsername)
		fmt.Println("password:", loginPassword)

		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		// Get the existing entry present in the database for the given username
		result := db.QueryRow("SELECT password FROM student WHERE username=?", loginUsername)
		if err != nil {
			// If there is an issue with the database, return a 500 error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Declare a variable to store the retrieved hashed password
		var hashedPassword string

		// Scan the result into the hashedPassword variable
		if err := result.Scan(&hashedPassword); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Compare the stored hashed password with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
			// If the two passwords don't match, return a 401 status
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If we reach this point, that means the user's password was correct, and they are authorized
		// The default 200 status is sent
		fmt.Println("Login successful")
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

