package interfaces

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdfpt05/go-event-sourcing/internal/middleware"
)

func NewRouter(ah *AccountHandler) *mux.Router {
	r := mux.NewRouter()
	
	api := r.PathPrefix("/api/v1").Subrouter()
	
	api.HandleFunc("/account/create", middleware.Logging(middleware.Recover(ah.CreateAccount))).Methods("POST")
	api.HandleFunc("/account/deposit", middleware.Logging(middleware.Recover(ah.Deposit))).Methods("POST")
	api.HandleFunc("/account/withdraw", middleware.Logging(middleware.Recover(ah.Withdraw))).Methods("POST")
	api.HandleFunc("/account/balance", middleware.Logging(middleware.Recover(ah.GetBalance))).Methods("GET")

	return r
}