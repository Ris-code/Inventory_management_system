package main

type item_info struct {
	Item_id  []string
	Club     string
	Club_id  string
	Items    []string
	Quantity []int
}

type club_info struct {
	Club_id []string
	Club    []string
	Info    []string
	Link    []string
}

type Item struct {
	ItemID   []string `json:"itemID"`
	Quantity []int    `json:"Quantity"`
	Club_id  string   `json:"club_id"`
	Club     string   `json:"club"`
	Return   string   `json:"returnDate"`
	Name     string   `json:"name"`
	ID       string   `json:"id"`
	Username string   `json:"username"`
}

type user struct {
	Username string
	ID       string
}

type email_Item struct {
	Name     string
	Quantity int
	Left     int
}
type return_email_Item struct {
	Name       string
	Quantity   int
	Left       int
	ReturnDate string
	Status     string
}

// Remove the duplicate declaration of inventory
type BorrowedItem struct {
	Name       string `bson:"name"`
	Quantity   int    `bson:"quantity"`
	ReturnDate string `bson:"return_date"`
}

type Club_present struct {
	Club_id       string         `bson:"club_id"`
	Club          string         `bson:"club"`
	Borrow_status string         `bson:"borrow_status"`
	Items         []BorrowedItem `bson:"items"`
}

type student struct {
	Username    string         `bson:"username"`
	Name        string         `bson:"name"`
	InstituteID string         `bson:"institute_id"`
	Club_info   []Club_present `bson:"club_info"`
}

type DeleteItemsRequest struct {
	Item       []string `json:"item"`
	Quantity   []string `json:"quantity"`
	ReturnDate []string `json:"returnDate"`
	Username   string   `json:"username"`
	Club       string   `json:"club"`
	ID         string   `json:"id"`
}

type Status struct {
	Status string `json:"status"`
}

type club_borrow struct {
	Club string `json:"club"`
}