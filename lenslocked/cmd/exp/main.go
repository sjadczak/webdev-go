package main

import (
	"fmt"

	"github.com/sjadczak/webdev-go/lenslocked/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("lenslocked> connected")

	us := models.UserService{
		DB: db,
	}

	user, err := us.Create("theo@jadczak.com", "bork123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("lenslocked> created user: %v\n", *user)

	// Create a table...
	//fmt.Println("lenslocked> creating tables...")
	//_, err = db.Exec(`
	//CREATE TABLE IF NOT EXISTS users (
	//id SERIAL PRIMARY KEY,
	//name TEXT,
	//email TEXT UNIQUE NOT NULL
	//);

	//CREATE TABLE IF NOT EXISTS orders (
	//id SERIAL PRIMARY KEY,
	//user_id INT,
	//amount INT,
	//description TEXT,
	//FOREIGN KEY(user_id) REFERENCES users(id)
	//);
	//`)
	//if err != nil {
	//panic(err)
	//}
	//fmt.Println("lenslocked> tables created.")

	// Insert some data
	//name := "Jessica Jadczak"
	//email := "jessica@test.com"
	//row := db.QueryRow(`
	//INSERT INTO users (name, email)
	//VALUES ($1, $2)
	//RETURNING id;
	//`, name, email)
	//var id int
	//err = row.Scan(&id)
	//if err != nil {
	//panic(err)
	//}
	//fmt.Printf("lenslocked> user created with id %d.\n")

	//id := 1
	//row := db.QueryRow(`
	//SELECT name, email
	//FROM users
	//WHERE id=$1;
	//`, id)
	//var name, email string
	//err = row.Scan(&name, &email)
	//if err != nil {
	//panic(err)
	//}
	//fmt.Printf("lenslocked> User Info: name=%s email=%s\n", name, email)

	//userID := 3
	//for i := 1; i <= 6; i++ {
	//amount := i * 20
	//desc := fmt.Sprintf("Fake order #%d", i+5)
	//_, err := db.Exec(`
	//INSERT INTO orders(user_id, amount, description)
	//VALUES ($1, $2, $3);
	//`, userID, amount, desc)
	//if err != nil {
	//panic(err)
	//}
	//}
	//fmt.Println("lenslocked> created faker orders.")

	//type Order struct {
	//ID          int
	//UserID      int
	//Amount      int
	//Description string
	//}
	//var orders []Order

	//userID := 1
	//rows, err := db.Query(`
	//SELECT id, amount, description
	//FROM orders
	//WHERE user_id=$1
	//`, userID)
	//if err != nil {
	//panic(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//var order Order
	//order.UserID = userID
	//err = rows.Scan(&order.ID, &order.Amount, &order.Description)
	//if err != nil {
	//panic(err)
	//}
	//orders = append(orders, order)
	//}
	//if err = rows.Err(); err != nil {
	//panic(err)
	//}
	//fmt.Printf("lenslocked> found %d orders for user %d.\n", len(orders), userID)
	//for _, order := range orders {
	//fmt.Printf("#> Order %d: Quantity=%d Description=%s\n", order.ID, order.Amount, order.Description)
	//}
}
