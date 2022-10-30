package main

import (
	"fmt"
	"github.com/chandanacharya1/customer-matching/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
