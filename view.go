package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"bytes"
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"regexp"
	"strconv"
)

var db *sql.DB
var err error

type item_info struct{
	Item_id []string
	Club string
	Club_id string
	Items []string
	Quantity []int
}

type club_info struct{
	Club_id []string
	Club []string
	Info []string
	Link []string
}

type Item struct {
	ItemID   []string `json:"itemID"`
	Quantity []int    `json:"Quantity"`
	Club_id string `json:"club_id"`
	Club string `json:"club"`
	Return string `json:"returnDate"`
	Name string `json:"name"`
	ID string `json:"id"`
}

type user struct {
	Username string 
}

type email_Item struct {
	Name     string
	Quantity int
	Left int
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

func home_before_login(w http.ResponseWriter, r *http.Request) {	
	t, _ := template.ParseFiles("templates/home_before_login.html")
	t.Execute(w, nil)
}

func home_after_login(w http.ResponseWriter, r *http.Request) {	
	t, _ := template.ParseFiles("templates/home_after_login.html")
	t.Execute(w, nil);

}



func club_option(w http.ResponseWriter, r *http.Request) {

	join_output, err := db.Query("SELECT club_id, club, Info, Img_link FROM clubs")

	// fmt.Println(join_output)
	if err != nil {
        // Handle the error (log it or return an error response)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer join_output.Close()

	id_arr := []string{}
	club_arr := []string{}
	info_arr := []string{}
	photolink_arr := []string{}

	for join_output.Next() {
		var id string
		var club string
		var info string
		var link string

		err = join_output.Scan(&id, &club, &info, &link)
		if err != nil {
			panic(err.Error())
		}

		id_arr = append(id_arr, id)
		club_arr = append(club_arr, club)
		info_arr = append(info_arr, info)
		photolink_arr = append(photolink_arr, link)
	}

	fmt.Println("ID:", id_arr)
	fmt.Println("Club:", club_arr)
	fmt.Println("Info:", info_arr)
	fmt.Println("Photo Link:", photolink_arr)


	temp := club_info{
		Club_id: id_arr,
		Club: club_arr,
		Info: info_arr,
		Link: photolink_arr,
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(temp)
	if err != nil {
		// Handle the error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template with the JSON data
	t, _ := template.ParseFiles("templates/club1.html")
	t.Execute(w, template.JS(jsonData))
	
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

	club_id := updateItem.Club_id
	club := updateItem.Club
	name := updateItem.Name
	return_date := updateItem.Return
	id := updateItem.ID
	email := db.QueryRow("SELECT email FROM clubs WHERE club_id=?", club_id)

	var email_id string
	if err := email.Scan(&email_id); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	email_item_arr := []email_Item{}

	for idx, itemID := range updateItem.ItemID {
		fmt.Println("Item ID:", itemID)
		fmt.Println("Quantity:", updateItem.Quantity[idx])

		get_item := db.QueryRow("SELECT item FROM items WHERE item_id=?", itemID)

		var setitem string

		if err := get_item.Scan(&setitem); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

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


		temp := email_Item{
			Name:   setitem,
			Quantity: updateItem.Quantity[idx],
			Left: updated_value,
		}
		email_item_arr = append(email_item_arr, temp)

		_, err := db.Exec("UPDATE items SET quantity=? WHERE item_id=?", updated_value, itemID)

		if err != nil {
			// If there is an issue with the database, return a 500 error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Println("Item ID:", updateItem.ItemID)
		fmt.Println("Quantity:", updateItem.Quantity)
	}
	
	
	fmt.Println("email_item_arr:", email_item_arr)

    send_email(email_id, club_id, name, return_date, club, id, email_item_arr)
}

// func readHTMLTemplate(path string) (string, error) {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(file), nil
// }

func send_email(email_id string, club_id string, name string, return_date string, club string, id string, items []email_Item) {
		// Send an email to the club

	// sender data
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	fmt.Println("Email:", from)
	// fmt.Println("Password:", password)

	// receiver email address
	to := []string{
		email_id, // could be more than one receiver email addresses separated by comma  
	}

	fmt.Println("Email ID:", to)
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"

	
	auth := smtp.PlainAuth("", from, password, host)

	t, _ := template.ParseFiles("templates/email_template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Item borrowed by student \n%s\n\n", mimeHeaders)))

	// Prepare data for template
	data := struct {
		Name       string
		ReturnDate string
		Club	   string
		Id         string
		Items      []email_Item
	}{
		Name:       name,
		ReturnDate: return_date,
		Club: club,
		Id:         id,
		Items:      items,
	}

	t.Execute(&body, data)

	err := smtp.SendMail(host+":"+port, auth, from, to, body.Bytes())
	if err != nil {
	  fmt.Println(err)
	  return
	}
	fmt.Println("Email Sent!")
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
		Club_id: id,
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
	fmt.Println("method:", r.Method) // get request method

	if r.Method == "GET" {
		// Serve the login page
		http.ServeFile(w, r, "templates/signin.html")
		return
	}

	r.ParseForm()
	// logic part of log in
	loginUsername := r.FormValue("username")
	loginPassword := r.FormValue("password")

	fmt.Println("username:", loginUsername)
	fmt.Println("password:", loginPassword)

	// Get the existing entry present in the database for the given username
	result := db.QueryRow("SELECT password FROM student WHERE username=?", loginUsername)
	result1 := db.QueryRow("SELECT name FROM student WHERE username=?", loginUsername)

	// Declare a variable to store the retrieved hashed password
	var hashedPassword string
	var name string

	if err := result1.Scan(&name); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.Write([]byte("unsuccessfull"))
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Name:", name)

	// Scan the result into the hashedPassword variable
	if err := result.Scan(&hashedPassword); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.Write([]byte("unsuccessfull"))
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
    

	// Compare the stored hashed password with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginPassword)); err != nil {
		// If the two passwords don't match, return a 401 status
		fmt.Println("unsuccessfull")
		w.Write([]byte("unsuccessfull"))
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("success")
	// If we reach this point, that means the user's password was correct, and they are authorized
	// Send a success response to the client
	temp := user{
		Username: name,
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(temp)
	if err != nil {
		// Handle the error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Name:", name)
	w.Write([]byte(jsonData))	
}


func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of sign up
		var username = r.FormValue("username")
		var password = r.FormValue("password")
		var confirm_password = r.FormValue("confirm-password")
		var name = r.FormValue("name")
		var id = r.FormValue("id")

		fmt.Println("username:", username)
		fmt.Println("password:", password)
		fmt.Println("confirm-password:", confirm_password)
		fmt.Println("name:", name)
		fmt.Println("id:", id)
		// Check if username already exists
		result := db.QueryRow("SELECT username FROM student WHERE username=?", username)

		var existingUsername string
		err := result.Scan(&existingUsername)

		if err == nil {
			// Username already exists
			fmt.Println("Username already exists")
			w.Write([]byte("username"))
			
			return
		} else if err != sql.ErrNoRows {
			// If the error is of any other type, send a 500 status
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if password != confirm_password {
			// Passwords do not match
			w.Write([]byte("password"))
			http.Redirect(w, r, "/signup/", http.StatusSeeOther)
			return
		}

		// Hash the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

		

		// Insert the new user into the database
		insert, err := db.Prepare("INSERT INTO student (username, name, password, Institute_id) VALUES (?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(username, name, hashedPassword, id)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Record inserted successfully")
		w.Write([]byte("success"))
	}
}

func thank_page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/thank.html")
	t.Execute(w, nil)
}



// coordinator functions

func coordinator_login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	if r.Method == "GET" {
		result,_ := db.Query("SELECT club_id, club FROM clubs")

		var club []string

		for result.Next() {
			var id string
			var name string

			err = result.Scan(&id, &name)
			if err != nil {
				panic(err.Error())
			}

			// temp := coordinator{
			// 	club_name: name,
			// 	club_id: id,
			// }

			club = append(club, name)
		}
		data := struct {
			Items []string
		}{
			Items: club,
		}
		fmt.Println("Club:", club)

		jsonData, err := json.Marshal(data)

		if err != nil {
			// Handle the error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		t, _ := template.ParseFiles("templates/club_signin.html")
		t.Execute(w, template.JS(jsonData))
	} else if r.Method == "POST" {
		r.ParseForm()
		// logic part of sign up
		var club = r.FormValue("club")
		var id = r.FormValue("id")

		fmt.Println("club:", club)
		fmt.Println("id:", id)
		// Check if username already exists
		result := db.QueryRow("SELECT club FROM clubs WHERE unique_id=?", id)

		var existingClub string
		if err := result.Scan(&existingClub) ; err != nil {
			fmt.Println("Unique ID does not exist")
			w.Write([]byte("unique_id"))

			return
		}

		// if err != nil {
		// 	// Username already exists
		// 	fmt.Println("Unique ID does not exist")
		// 	w.Write([]byte("unique_id"))
			
		// 	return
		// } else if err != sql.ErrNoRows {
		// 	// If the error is of any other type, send a 500 status
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return
		// }

		if existingClub != club {
			// Passwords do not match
			w.Write([]byte("not_club"))
			// http.Redirect(w, r, "/signup/", http.StatusSeeOther)
			return
		}

		fmt.Println("Entered successfully")
		w.Write([]byte("success"))
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}

func club_home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/club_home.html")
	t.Execute(w, nil)
}

func extractNumericPart(itemID string) (int, error) {
	// Use regular expression to extract numeric part
	re := regexp.MustCompile(`(\d+)`)
	match := re.FindStringSubmatch(itemID)

	if len(match) < 2 {
		return 0, fmt.Errorf("numeric part not found in item ID: %s", itemID)
	}

	numericPart, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("failed to convert numeric part to integer: %v", err)
	}

	return numericPart, nil
}

func add_inventory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

		r.ParseForm()
		// logic part of sign up
		var item = r.FormValue("item")
		var quantity = r.FormValue("quantity")
		var name = r.FormValue("name")

		fmt.Println("item:", item)
		fmt.Println("quantity:", quantity)

		// Get the existing entry present in the database for the given username
		result := db.QueryRow("SELECT club_id FROM clubs WHERE club=?", name)

		// Declare a variable to store the retrieved hashed password
		var club_id string

		if err := result.Scan(&club_id); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				w.Write([]byte("unsuccessfull"))
				// http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Println("club_id:", club_id)


		id_result :=db.QueryRow("SELECT item_id FROM items ORDER BY item_id DESC LIMIT 1")

		var item_id string

		if err := id_result.Scan(&item_id); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				w.Write([]byte("unsuccessfull"))
				// http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Println("item_id:", item_id)

		numericPart, err := extractNumericPart(item_id)

		if err != nil {
			fmt.Println(err)
			return
		}

		numericPart = numericPart + 1

		stringnumericPart := strconv.Itoa(numericPart)

		// fmt.Println("stringnumericPart:", len(stringnumericPart))
		if len(stringnumericPart) == 1 {
			item_id = "IT0" + strconv.Itoa(numericPart)
		} else {
			item_id = "IT" + strconv.Itoa(numericPart)
		}

		// Insert the new user into the database
		insert, err := db.Prepare("INSERT INTO items (item_id, item, quantity, club_id) VALUES (?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(item_id, item, quantity, club_id)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Record inserted successfully")
		w.Write([]byte("success"))
	}
