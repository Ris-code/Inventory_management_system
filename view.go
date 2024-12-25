package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"crypto/tls"
	"crypto/x509"

	"github.com/go-sql-driver/mysql"

	"bytes"
	"encoding/json"
	"os"

	"net/smtp"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"log"
	"time"
	"strings"
)

var db *sql.DB
var collection *mongo.Collection
var err error

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	initDB()
	initMongoDB()
}

func initDB() {
	// caCertPath := "ca.pem"
	caCert := os.Getenv("CA_CERT")
	caCert = strings.ReplaceAll(caCert, "\\n", "\n")

	// if err != nil {
	// 	log.Fatalf("Error reading CA certificate: %v", err)
	// }

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM([]byte(caCert)) {
		log.Fatal("Failed to append CA certificate")
	}

	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatalf("Error registering custom TLS config: %v", err)
	}

	dsn := os.Getenv("DB_URI")
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connected to the MySQL database!")
}

func initMongoDB() {
	connectionString := os.Getenv("CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	collection = client.Database("Student_inventory_list").Collection("Item_list")
	log.Println("Connected to MongoDB!")
}

func ensureDBInitialized(w http.ResponseWriter) bool {
	if db == nil {
		http.Error(w, "Database connection not initialized", http.StatusInternalServerError)
		return false
	}
	return true
}

func home_before_login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home_before_login.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func home_after_login(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/home_after_login.html")
	t.Execute(w, nil)

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
	// fmt.Println("Photo Link:", photolink_arr)

	temp := club_info{
		Club_id: id_arr,
		Club:    club_arr,
		Info:    info_arr,
		Link:    photolink_arr,
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
	username := updateItem.Username
	email := db.QueryRow("SELECT email FROM clubs WHERE club_id=?", club_id)

	filter := bson.M{"username": username, "club_info.club": club}

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

		updated_value := val - updateItem.Quantity[idx]

		temp := email_Item{
			Name:     setitem,
			Quantity: updateItem.Quantity[idx],
			Left:     updated_value,
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

		update := bson.M{
			"$set": bson.M{
				"club_info.$.borrow_status": "Yes",
			},
			"$push": bson.M{
				"club_info.$.items": bson.M{
					"name":        setitem,
					"quantity":    updateItem.Quantity[idx],
					"return_date": return_date,
				},
			},
		}

		updateResult, err := collection.UpdateOne(context.Background(), filter, update)

		if err != nil {
			// log.Fatal(err)
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Matched %v document and modified %v document.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	fmt.Println("email_item_arr:", email_item_arr)

	send_email(email_id, club_id, name, return_date, club, id, email_item_arr)
}

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
		Club       string
		Id         string
		Items      []email_Item
	}{
		Name:       name,
		ReturnDate: return_date,
		Club:       club,
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

func inventory(w http.ResponseWriter, r *http.Request) {
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
		Item_id:  join_arr_id,
		Club:     club,
		Club_id:  id,
		Items:    join_arr,
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
	if !ensureDBInitialized(w) {
		return
	}

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

	institute_id := db.QueryRow("SELECT Institute_id FROM student WHERE username=?", loginUsername)

	var id string

	if err := institute_id.Scan(&id); err != nil {
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

	temp := user{
		Username: name,
		ID:       id,
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

		// uppercase the id
		id = strings.ToUpper(id)

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

		all_club, _ := db.Query("SELECT club_id, club FROM clubs")

		var club_data []Club_present

		for all_club.Next() {
			var club_id string
			var club string

			err = all_club.Scan(&club_id, &club)
			if err != nil {
				panic(err.Error())
			}

			temp := Club_present{
				Club_id: club_id,
				Club:    club,
				Items:   []BorrowedItem{},
			}

			club_data = append(club_data, temp)
		}
		// Insert a new document into the collection.
		user := student{Username: username, Name: name, InstituteID: id, Club_info: club_data}

		fmt.Println("user:", user)

		insertResult, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)

		// Fetch the inserted document from the collection
		var insertedUser student
		filter := bson.M{"_id": insertResult.InsertedID}
		err = collection.FindOne(context.Background(), filter).Decode(&insertedUser)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("Error:", err)
			return
		}

		// Print the entire inserted document
		fmt.Println("Inserted User:", insertedUser)
	}
}

func thank_page(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/thank.html")
	t.Execute(w, nil)
}

func borrow_list(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	if r.Method == "GET" {
		result, _ := db.Query("SELECT club_id, club FROM clubs")

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

		t, _ := template.ParseFiles("templates/Borrow_list.html")
		t.Execute(w, template.JS(jsonData))
	} else if r.Method == "POST" {
		r.ParseForm()

		var club = r.FormValue("club")
		var username = r.FormValue("username")

		fmt.Println("club:", club)

		filter := bson.M{"username": username, "club_info.club": club, "club_info.borrow_status": "Yes"}

		var student student

		err := collection.FindOne(context.Background(), filter).Decode(&student)

		if err != nil {
			// log.Fatal(err)
			var status Status
			status.Status = "unsuccessfull"

			jsonData, err := json.Marshal(status)

			w.Write([]byte(jsonData))
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Student:", student)

		temp := student.Club_info

		fmt.Println("Temp:", temp)

		var items []BorrowedItem

		for _, item := range temp {
			if item.Club == club {
				items = item.Items
			}
		}

		fmt.Println("Items:", items)

		// Convert the struct to JSON
		jsonData, err := json.Marshal(items)

		if err != nil {
			// Handle the error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(jsonData))
	}
}

func delete_items(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	// r.ParseForm()

	decoder := json.NewDecoder(r.Body)
	var deleteItemsRequest DeleteItemsRequest
	err := decoder.Decode(&deleteItemsRequest)
	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
		return
	}

	item := deleteItemsRequest.Item
	quantity := deleteItemsRequest.Quantity
	returnDate := deleteItemsRequest.ReturnDate
	username := deleteItemsRequest.Username
	club := deleteItemsRequest.Club
	id := deleteItemsRequest.ID

	fmt.Println("item:", item)
	fmt.Println("quantity:", quantity)
	fmt.Println("returnDate:", returnDate)

	name_id := db.QueryRow("SELECT name FROM student WHERE username=?", username)

	var name string

	if err := name_id.Scan(&name); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status

		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result := db.QueryRow("SELECT email FROM clubs WHERE club=?", club)

	var email_id string

	if err := result.Scan(&email_id); err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	delete_item_arr := []return_email_Item{}

	for idx, itemID := range item {
		var delete_item = string(itemID)
		var delete_quantity_str = quantity[idx]
		var delete_returnDate = string(returnDate[idx])
		delete_quantity, _ := strconv.Atoi(delete_quantity_str)

		fmt.Println("Item:", delete_item)
		fmt.Println("Quantity:", delete_quantity)
		fmt.Println("Return Date:", delete_returnDate)
		fmt.Println("Username:", username)
		fmt.Println("Club:", club)

		filter := bson.M{"username": username, "club_info.club": club, "club_info.items.name": delete_item}

		update := bson.M{
			"$pull": bson.M{
				"club_info.$.items": bson.M{
					"name":        delete_item,
					"quantity":    delete_quantity,
					"return_date": delete_returnDate,
				},
			},
		}

		// Execute the $pull operation
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("Error:", err)
			return

		}

		filter = bson.M{"username": username}

		// Then, check the length of the items array
		var result bson.M
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("Error:", err)
			return
		}

		// Check if "club_info" field exists
		clubInfo, ok := result["club_info"].(bson.A)
		if !ok {
			// log.Fatal("Club information not found.")
			fmt.Println("Club information not found.")
			return
		}

		// Find the specific club
		var targetClub bson.M
		for _, c := range clubInfo {
			if c.(bson.M)["club"].(string) == club {
				targetClub = c.(bson.M)
				break
			}
		}

		// Check if the "items" field exists and is an array (bson.A)
		item_arr, _ := targetClub["items"].(bson.A)

		fmt.Println("Item Array:", item_arr)
		fmt.Println("Length:", len(item_arr))

		var targetClubIndex int
		for i, c := range clubInfo {
			if c.(bson.M)["club"].(string) == club {
				targetClubIndex = i
				break
			}
		}

		if len(item_arr) == 0 {
			// If it's 0, set borrow_status to "No"
			update = bson.M{
				"$set": bson.M{
					fmt.Sprintf("club_info.%d.borrow_status", targetClubIndex): "No",
				},
			}
			_, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				// log.Fatal(err)
				fmt.Println("Error:", err)
			}
		}

		// Update the quantity of the item in the database
		value := db.QueryRow("SELECT quantity FROM items WHERE item=?", delete_item)

		var val int

		if err := value.Scan(&val); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		updated_value := val + delete_quantity

		currentTime := time.Now()
		// Time := currentTime.Format("2006-01-02")

		compareDate, _ := time.Parse("2006-01-02", delete_returnDate)

		var status string
		if currentTime.After(compareDate) {
			status = "Late"
		} else {
			status = "Within Time"
		}

		temp := return_email_Item{
			Name:       delete_item,
			Quantity:   delete_quantity,
			Left:       updated_value,
			ReturnDate: delete_returnDate,
			Status:     status,
		}

		delete_item_arr = append(delete_item_arr, temp)

		_, err = db.Exec("UPDATE items SET quantity=? WHERE item=?", updated_value, delete_item)

		if err != nil {
			// If there is an issue with the database, return a 500 error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Println("Item ID:", item)
		fmt.Println("Quantity:", quantity)
		fmt.Println("Return Date:", returnDate)
		fmt.Println("Username:", username)
		fmt.Println("Club:", club)

	}

	send_return_email(name, club, id, email_id, delete_item_arr)
}

func send_return_email(username string, club string, id string, email_id string, items []return_email_Item) {
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

	t, _ := template.ParseFiles("templates/return_email_template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Item returned by student \n%s\n\n", mimeHeaders)))

	// Prepare data for template
	data := struct {
		Name  string
		Club  string
		Id    string
		Items []return_email_Item
	}{
		Name:  username,
		Club:  club,
		Id:    id,
		Items: items,
	}

	t.Execute(&body, data)

	err := smtp.SendMail(host+":"+port, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

// coordinator functions

func coordinator_login(w http.ResponseWriter, r *http.Request) {
	if !ensureDBInitialized(w) {
		return
	}

	fmt.Println("method:", r.Method) // log request method

	if r.Method == "GET" {
		// Query database for clubs
		rows, err := db.Query("SELECT club_id, club FROM clubs")
		if err != nil {
			log.Printf("Error querying clubs: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var clubs []string
		for rows.Next() {
			var id, name string
			if err := rows.Scan(&id, &name); err != nil {
				log.Printf("Error scanning row: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			clubs = append(clubs, name)
		}

		// Handle potential errors from rows iteration
		if err = rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Items []string
		}{
			Items: clubs,
		}
		fmt.Println("Clubs:", clubs)

		// Convert to JSON and serve the HTML
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/club_signin.html")
		if err != nil {
			log.Printf("Error loading template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, template.JS(jsonData))
	} else if r.Method == "POST" {
		r.ParseForm()
		club := r.FormValue("club")
		id := r.FormValue("id")

		fmt.Println("club:", club)
		fmt.Println("id:", id)

		var existingClub string
		err := db.QueryRow("SELECT club FROM clubs WHERE unique_id=?", id).Scan(&existingClub)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid Unique ID", http.StatusUnauthorized)
				return
			}
			log.Printf("Error querying club: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if existingClub != club {
			http.Error(w, "Club name does not match Unique ID", http.StatusUnauthorized)
			return
		}

		fmt.Println("Coordinator login successful")
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

	item_check := db.QueryRow("SELECT item FROM items WHERE item=?", item)

	var existingItem string
	if err := item_check.Scan(&existingItem); err == nil {
		// Username already exists
		fmt.Println("Item already exists")

		data := Status{Status: "item_exists"}
		jsonData, _ := json.Marshal(data)
		w.Write([]byte(jsonData))

		return
	}
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

	id_result := db.QueryRow("SELECT item_id FROM items ORDER BY item_id DESC LIMIT 1")

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
	data := Status{Status: "success"}
	fmt.Println("Record inserted successfully")
	// w.Write([]byte("success"))
	jsonData, _ := json.Marshal(data)

	w.Write([]byte(jsonData))

}

func update_info(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	r.ParseForm()
	// logic part of sign up
	var info = r.FormValue("desp")
	var pic = r.FormValue("pic")
	var email = r.FormValue("email")
	var name = r.FormValue("name")

	fmt.Println("info:", info)
	fmt.Println("pic:", pic)
	fmt.Println("email:", email)
	fmt.Println("name:", name)

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

	if pic != "" {
		// Insert the new user into the database
		insert, err := db.Prepare("UPDATE clubs SET Img_link=? WHERE club_id=?")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(pic, club_id)
		if err != nil {
			panic(err.Error())
		}
	}

	if info != "" {
		// Insert the new user into the database
		insert, err := db.Prepare("UPDATE clubs SET Info=? WHERE club_id=?")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(info, club_id)
		if err != nil {
			panic(err.Error())
		}

	}

	if email != "" {
		insert, err := db.Prepare("UPDATE clubs SET email=? WHERE club_id=?")
		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		// Execute the prepared statement with form values
		_, err = insert.Exec(info, club_id)
		if err != nil {
			panic(err.Error())
		}

	}

	data := Status{Status: "success"}

	jsonData, _ := json.Marshal(data)
	// fmt.Println("Record inserted successfully")
	w.Write([]byte(jsonData))
}

func club_borrow_list(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) // get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/club_borrower_list.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {

		var requestData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Extract the club name from the JSON data
		club := requestData["club"]

		fmt.Println("club:", club)

		// select all the users and their items they borrowed from this club from mongoDB

		filter := bson.M{"club_info.club": club, "club_info.borrow_status": "Yes"}

		cur, err := collection.Find(context.Background(), filter)
		fmt.Println("Cursor:", cur)

		if err != nil {

			// log.Fatal(err)
			fmt.Println("Error:", err)
			return

		}

		defer cur.Close(context.Background())

		var students []student

		for cur.Next(context.Background()) {

			var s student

			err := cur.Decode(&s)

			if err != nil {

				// log.Fatal(err)
				fmt.Println("Error:", err)
				return

			}

			// Filter the student data to include only the "Sangam" club
			var filteredClubInfo []Club_present
			for _, c := range s.Club_info {
				if c.Club == club && c.Borrow_status == "Yes" {
					filteredClubInfo = append(filteredClubInfo, c)
				}
			}
			fmt.Println("Filtered Club Info:", len(filteredClubInfo))
			s.Club_info = filteredClubInfo

			if len(filteredClubInfo) > 0 {
				students = append(students, s)
			}

		}

		if err := cur.Err(); err != nil {

			// log.Fatal(err)
			fmt.Println("Error:", err)
			return

		}

		cur.Close(context.Background())

		fmt.Println("Student:", students)

		// Convert the struct to JSON

		jsonData, _ := json.Marshal(students)

		w.Write([]byte(jsonData))
	}
}

func inventorylist(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) // get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/inventory_list.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {

		var requestData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Extract the club name from the JSON data
		club := requestData["club"]

		fmt.Println("club:", club)

		// select all the users and their items they borrowed from this club from mongoDB
		id := db.QueryRow("SELECT club_id FROM clubs WHERE club=?", club)

		var club_id string

		if err := id.Scan(&club_id); err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		result, err := db.Query("SELECT items.item, items.quantity FROM items INNER JOIN clubs ON clubs.club_id=items.club_id WHERE clubs.club_id=?", club_id)

		if err != nil {
			// Handle the error (log it or return an error response)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var items []item_list

		for result.Next() {
			var item string
			var quantity int

			_ = result.Scan(&item, &quantity)
			temp := item_list{
				Item:     item,
				Quantity: quantity,
			}

			items = append(items, temp)
		}

		// Convert the struct to JSON
		fmt.Println("Items:", items)
		jsonData, _ := json.Marshal(items)
		// fmt.Println("JSON:", jsonData)

		w.Write([]byte(jsonData))
	}
}

func edit_inventory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // get request method

	// r.ParseForm()
	// // logic part of sign up
	// var item = r.FormValue("item")
	// var quantity = r.FormValue("quantity")
	decoder := json.NewDecoder(r.Body)
	var item item_list
	err := decoder.Decode(&item)
	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
		return
	}

	item_name := item.Item
	quantity := item.Quantity

	fmt.Println("item:", item_name)
	fmt.Println("quantity:", quantity)

	// Insert the new user into the database
	insert, err := db.Prepare("UPDATE items SET quantity=? WHERE item=?")
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	// Execute the prepared statement with form values
	_, err = insert.Exec(quantity, item_name)
	if err != nil {
		panic(err.Error())
	}
	data := Status{Status: "success"}
	fmt.Println("Record inserted successfully")
	// w.Write([]byte("success"))
	jsonData, _ := json.Marshal(data)

	w.Write([]byte(jsonData))
}

func delete_inventory(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the club name from the JSON data
	item := requestData["item"]

	fmt.Println("item:", item)

	// delete item from items

	_, err := db.Exec("DELETE FROM items WHERE item=?", item)

	if err != nil {
		// If there is an issue with the database, return a 500 error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := Status{Status: "success"}

	jsonData, _ := json.Marshal(data)

	w.Write([]byte(jsonData))
}
