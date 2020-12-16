package middleware

import (
	"database/sql"
	"encoding/json" //package to encode & decode json into struct & vice versa
	"fmt"
	"go-postgres/models" //Job schema model package
	"log"
	"net/http" //to acess request and response object of api
	"os"       //to read environment variables to ensure security

	//to get params from route
	"github.com/joho/godotenv" //package to read .env
	//postgres goland driver
)

//response struct
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

//create connection to postgres db
func createConnection() *sql.DB {
	//load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//open connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	//check connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("DB Connection Success")
	//return connection
	return db
}

// PostJob creates a job posting in db
func PostJob(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//create an empty job of type models.Job
	var job models.Job

	//decode json request to job
	err := json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	//call insert job function and pass the job
	insertID := insertJob(job)

	//format response object
	res := response{
		ID:      insertID,
		Message: "Job added successfully",
	}

	json.NewEncoder(w).Encode(res)
}

//GetJob will return job by its id
func GetUser()
