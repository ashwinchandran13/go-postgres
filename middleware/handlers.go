package middleware

import (
	"database/sql"
	"encoding/json" //package to encode & decode json into struct & vice versa
	"fmt"
	"go-postgres/models" //Job schema model package
	"log"
	"net/http" //to acess request and response object of api
	"os"       //to read environment variables to ensure security
	"strconv"

	//to get params from route
	"github.com/gorilla/mux"
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
func GetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//get the jobid from request params, key is "id"
	params := mux.Vars(r)

	//calling getJob function with job id to retrieve that job details
	job, err := getJob(string(params["name"]))

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	//send the response
	json.NewEncoder(w).Encode(job)
}

// GetAllJob will return all the users
func GetAllJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	jobs, err := getAllJobs()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(jobs)
}

// DeleteJob to delete entries
func DeleteJob(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//get job id from the request params
	params := mux.Vars(r)

	//convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to covert string id into int. %v", err)
	}

	//call the deleteJob, convert int to int64
	deletedRows := deleteJob(int64(id))

	//format the message string
	msg := fmt.Sprintf("job deleted successfully. Rows affected %v", deletedRows)

	// format response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------Handler Functions check 2------------
// insert user in DB

func insertJob(job models.Job) int64 {

	// create postgres DB connection
	db := createConnection()

	// close db connection
	defer db.Close()

	// create the insert sql queery
	// returning jobid will return the id of the inserted job
	sqlStatement := `INSERT INTO jobs (jobname, location, openings) VALUES ($1, $2, $3) RETURNING jobid`

	// inserted id will store in this id
	var id int64

	// execute the sql statement
	// scan function will save insert id in the id
	err := db.QueryRow(sqlStatement, job.Name, job.Location, job.Openings).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single recored %v", id)

	// return the inserted id
	return id
}

func getJob(name string) (models.Job, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a job of models.Job type
	var job models.Job

	// create the select sql query
	sqlStatement := `SELECT * FROM jobs WHERE jobname=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, name)

	// unmarshal the row object to job
	err := row.Scan(&job.ID, &job.Name, &job.Openings, &job.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Printf("No entry with job name as: %v", name)
		return job, nil
	case nil:
		return job, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty job on error
	return job, err
}

// get one user from the DB by its userid
func getAllJobs() ([]models.Job, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var jobs []models.Job

	// create the select sql query
	sqlStatement := `SELECT * FROM jobs`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var job models.Job

		// unmarshal the row object to user
		err = rows.Scan(&job.ID, &job.Name, &job.Openings, &job.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		jobs = append(jobs, job)

	}

	// return empty user on error
	return jobs, err
}

// delete user in the DB
func deleteJob(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM jobs WHERE jobid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total entries deleted :%v", rowsAffected)

	return rowsAffected
}
