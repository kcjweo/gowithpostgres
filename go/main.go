package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// Book Struct (Model)
type Book struct {
	ID          string `json:"id"`
	Category_id int    `json:"category_id"`
	Author_id   int    `json:"author_id"`
	Title       string `json:"title"`
	UPDATED_AT  string `json:"updated_at"`
	CREATED_AT  string `json:"created_at"`
}

// !! TODO !!
//
// type Author struct {
//  ID          string `json:"id"`
// 	Firstname string `json:"firstname"`
// 	Lastname  string `json:"lastname"`
// }
// type Category struct {
//  ID          string `json:"id"`
// 	Category_name string `json:"category_name"`
// }

// Init books var as a slice Book struct
var books []Book

func createConnection() (*sql.DB, error) {
	Db, err := sql.Open("postgres", "postgres://app_user:yourP@ss2022@postgres:5432/app_db?sslmode=disable")

	return Db, err
}

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book = nil // reset
	Db, err := createConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	// select all books from DB
	rows, err := Db.Query("SELECT id, title, category_id, author_id, updated_at, created_at FROM books;")
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var t Book
		rows.Scan(&t.ID, &t.Title, &t.Category_id, &t.Author_id, &t.UPDATED_AT, &t.CREATED_AT)
		books = append(books, t)
	}

	// return with json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Single Books
func getBook(w http.ResponseWriter, r *http.Request) {
	var books []Book = nil // reset
	Db, err := createConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	params := mux.Vars(r) // Get params

	// select all books from DB
	get_one_sql := fmt.Sprintf("SELECT id, title, category_id, author_id, updated_at, created_at FROM books WHERE id = %s;", params["id"])
	rows, err := Db.Query(get_one_sql)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var t Book
		rows.Scan(&t.ID, &t.Title, &t.Category_id, &t.Author_id, &t.UPDATED_AT, &t.CREATED_AT)
		books = append(books, t)
	}

	// return with json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Create a New Books
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	Db, err := createConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	insert_sql := fmt.Sprintf("INSERT INTO BOOKS (category_id, author_id, title, created_at) VALUES(%d, %d, '%s', CURRENT_TIMESTAMP);", book.Category_id, book.Author_id, book.Title)
	_, err = Db.Query(insert_sql)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return with json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func updateIds(column string, value int, id string) error {
	Db, err := createConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}

	update_sql := fmt.Sprintf("UPDATE BOOKS SET %s = %d WHERE id = %s;", column, value, id)
	_, err = Db.Query(update_sql)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// Update a New Books
func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	params := mux.Vars(r)

	if book.Category_id != 0 {
		err := updateIds("category_id", book.Category_id, params["id"])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if book.Author_id != 0 {
		err := updateIds("author_id", book.Author_id, params["id"])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if book.Title != "" {
		Db, err := createConnection()
		if err != nil {
			fmt.Println(err)
			return
		}
		update_sql := fmt.Sprintf("UPDATE BOOKS SET title = '%s' WHERE id = %s;", book.Title, params["id"])
		_, err = Db.Query(update_sql)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// return with json
	w.Header().Set("Content-Type", "application/json")
	return_message := fmt.Sprintf("UPDATED id: %s", params["id"])
	json.NewEncoder(w).Encode(return_message)
}

// Delete a Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	Db, err := createConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	params := mux.Vars(r)
	delete_sql := fmt.Sprintf("DELETE FROM BOOKS WHERE id = %s;", params["id"])
	_, err = Db.Query(delete_sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	return_message := fmt.Sprintf("DELETED id: %s", params["id"])

	// return with json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(return_message)
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Router Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
