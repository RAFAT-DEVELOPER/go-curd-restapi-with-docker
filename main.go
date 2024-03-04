package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *pgxpool.Pool

func main() {
	initDB()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func initDB() {
	var err error
	db, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	createTableIfNotExists()
}

func createTableIfNotExists() {
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT
		)
	`)
	if err != nil {
		log.Fatal("Unable to create table:", err)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		handleError(w, err)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			handleError(w, err)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	err := db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		handleError(w, err)
		return
	}

	err := db.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		handleError(w, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec(context.Background(), "UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
	if err != nil {
		handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		handleError(w, err)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		handleError(w, fmt.Errorf("no user with ID %s found", id))
		return
	}

	json.NewEncoder(w).Encode("User deleted")
}

func handleError(w http.ResponseWriter, err error) {
	log.Println("Error:", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
