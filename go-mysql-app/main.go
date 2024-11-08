package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

// User struct represents a user in MySQL
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var db *sql.DB

// Initialize the MySQL database connection
func initDB() {
    var err error
    dsn := "yourUsername:yourPassword@tcp(127.0.0.1:3306)/yourDatabaseName"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error connecting to database: %v\n", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Error pinging database: %v\n", err)
    }

    log.Println("Connected to MySQL database")
}

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, _ := result.LastInsertId()
    user.ID = int(id)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func main() {
    initDB()
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")

    log.Println("Server is running on port 8000")
    log.Fatal(http.ListenAndServe(":8000", router))
}
