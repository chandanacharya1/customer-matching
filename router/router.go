package router

import (
	"fmt"
	"github.com/chandanacharya1/customer-matching/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	/*router.HandleFunc("/api/jwt", middleware.GetJWT).Methods("POST", "OPTIONS")*/
	router.HandleFunc("/api/login", middleware.Login).Methods("POST", "OPTIONS")
	router.Handle("/api/request", middleware.Validate(middleware.CustomerRequest)).Methods("POST", "OPTIONS")
	router.Handle("/api/partner", middleware.Validate(middleware.ListPartners)).Methods("GET", "OPTIONS")
	router.Handle("/api/partner/{partner_id}", middleware.Validate(middleware.GetPartner)).Methods("GET")
	fmt.Println("Starting API server on port 8080")
	http.ListenAndServe(":8080", router)
	return router
}
