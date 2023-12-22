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
	"github.com/gorilla/mux"
	"encoding/json"
)

var db *sql.DB
var err error

type item_info struct{
	Item_id []string
	Club string
	Items []string
	Quantity []int
}

type Item struct {
	ItemID   []string `json:"itemID"`
	Quantity []int    `json:"Quantity"`
}

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

func club_option(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/club1.html")
	t.Execute(w, nil)
}

func cart(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/cart.html")
	t.Execute(w, nil)
}

func update(w http.ResponseWriter, r *http.Request) {
	var updateItem Item

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateItem); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	for idx, itemID := range updateItem.ItemID {
		fmt.Println("Item ID:", itemID)
		fmt.Println("Quantity:", updateItem.Quantity[idx])

		value := db.QueryRow("SELECT quantity FROM items WHERE item_id=?", itemID)
		// Update the quantity of the item in the database
		var val int
		if err := value.Scan(&val); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		
		updated_value:=val-updateItem.Quantity[idx]

		_, err := db.Exec("UPDATE items SET quantity=? WHERE item_id=?", updated_value, itemID)

		if err != nil {
			// If there is an issue with the database, return a 500 error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("Item ID:", updateItem.ItemID)
	fmt.Println("Quantity:", updateItem.Quantity)

}

func inventory(w http.ResponseWriter, r *http.Request){
	fmt.Println("method:", r.Method)
	t, _ := template.ParseFiles("templates/inventory.html")

	vars := mux.Vars(r)
    id := vars["ID"]
	fmt.Println("ID:", id)

	result := db.QueryRow("SELECT club FROM clubs WHERE club_id=?", id)

	var club string
	
	if err := result.Scan(&club); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Club:", club)

	join_output, err := db.Query("SELECT  items.item_id, items.item, items.quantity, items.club_id FROM items INNER JOIN clubs ON clubs.club_id=items.club_id WHERE items.club_id=?", id)

	// fmt.Println(join_output)
	if err != nil {
        // Handle the error (log it or return an error response)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer join_output.Close()

	join_arr := []string{}
	join_arr_id := []string{}
	join_arr_quantity := []int{}

	for join_output.Next() {
		var item_id string
		var item string
		var quantity int
		var club_id string

		err = join_output.Scan(&item_id, &item, &quantity, &club_id)
		if err != nil {
			panic(err.Error())
		}

		join_arr = append(join_arr, item)
		join_arr_quantity = append(join_arr_quantity, quantity)
		join_arr_id = append(join_arr_id, item_id)
	}

	fmt.Println("Items:", join_arr)
	fmt.Println("Quantity:", join_arr_quantity)


	temp := item_info{
		Item_id: join_arr_id,
		Club: club,
		Items: join_arr,
		Quantity: join_arr_quantity,
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(temp)
	if err != nil {
		// Handle the error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template with the JSON data
	t.Execute(w, template.JS(jsonData))
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

