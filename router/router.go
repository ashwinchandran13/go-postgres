package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

//Router is exported and used in main.go
func Router *mux.Router() {

	router := mux.NewRouter()

	router.HandleFunc("/api/job/{id}", middleware.GetJob).Methods("GET", "OPTIONS") 
	router.HandleFunc("/api/job/{name}", middleware.FilterJob).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newjob", middleware.PostJob).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletejob/{id}", middleware.DeleteJob).Methids("DELETE", "OPTIONS")
	
	return router
} 
