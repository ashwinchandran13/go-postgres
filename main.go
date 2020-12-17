package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	//check 1
	fmt.Println("Starting server on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
