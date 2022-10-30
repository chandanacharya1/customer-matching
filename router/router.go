package router

import (
	"github.com/chandanacharya1/customer-matching/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/login", middleware.GetJWT).Methods("POST", "OPTIONS")
	router.Handle("/request", middleware.ValidateJWT(middleware.CustomerRequest)).Methods("POST", "OPTIONS")
	router.Handle("/partner", middleware.ValidateJWT(middleware.ListPartners)).Methods("GET", "OPTIONS")
	router.Handle("/partner/{partnerid}", middleware.ValidateJWT(middleware.GetPartner)).Methods("GET")
	return router
}
