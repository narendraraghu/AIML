package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Add this line to import the PostgreSQL driver
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "bunny"
)

func main() {
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Create the "users" table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Error creating 'users' table:", err)
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction:", err)
	}

	// Insert a record into the "users" table
	var userID int
	err = tx.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id;", "John").Scan(&userID)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error inserting record into 'users' table:", err)
	}

	// Insert a record into the "user_addresses" table with an invalid user_id (9999)
	_, err = tx.Exec("INSERT INTO user_addresses (user_id, address) VALUES ($1, $2);", 9999, "123 Main St")
	if err != nil {
		tx.Rollback()
		log.Fatal("Error inserting record into 'user_addresses' table:", err)
	}

	// Simulate an error during the second insert (invalid column name 'agee')
	_, err = tx.Exec("INSERT INTO users (name, agee) VALUES ($1, $2);", "Jane", 28)
	if err != nil {
		tx.Rollback()
		log.Println("Error inserting record into 'users' table:", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	fmt.Println("Records inserted successfully!")
}
