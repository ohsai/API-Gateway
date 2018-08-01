package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	// "Signin" and "Signup" are handler that we will implement
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	// initialize our database connection
	initDB()
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDB() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", "user=ubuntu password=9173162abc dbname=mydb sslmode=disable")
	if err != nil {
		log.Println(err)
	}
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}
type DBCredentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
	Role     string `json:"role", db:"role"`
}
type ResponseCred struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &DBCredentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	// Next, insert the username, along with the hashed password into the database
	if _, err = db.Query("insert into users (username, password, role) values ($1, $2, $3)", creds.Username, string(hashedPassword), creds.Role); err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func Signin(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	// Get the existing entry present in the database for the given username
	//result := db.QueryRow("select password from users where username=$1", creds.Username)
	result := db.QueryRow("select * from users where username=$1", creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &DBCredentials{}
	// Store the obtained password in `storedCreds`
	//err = result.Scan(&storedCreds.Password)
	var role_temp sql.NullString
	err = result.Scan(&storedCreds.Username, &storedCreds.Password, &role_temp)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err.Error())
			return
		}
		// If the error is of any other type, send a 500 status
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err.Error())
		return
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
	if role_temp.Valid {
		storedCreds.Role = role_temp.String
	}
	responsecred := ResponseCred{
		Username: storedCreds.Username,
		Role:     storedCreds.Role,
	}
	b, err := json.Marshal(responsecred)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(b)
	if err != nil {
		log.Println(err.Error())
		return
	}

}
