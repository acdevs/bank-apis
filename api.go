package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
	
type APIServer struct {
	ListenAddress string
	Store Storage
}

func NewAPIServer(ListenAddress string, store Storage) *APIServer {
	return &APIServer{
		ListenAddress: ListenAddress,
		Store: store,
	}
}


func (server *APIServer) Run(){
	router := mux.NewRouter()

	router.HandleFunc("/api/accounts", server.handleAccount)
	router.HandleFunc("/api/accounts/{id}", server.handleGetAccount)

	log.Println("Listening on", server.ListenAddress)
	http.ListenAndServe(server.ListenAddress, router)
}


func (server *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		server.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		server.handleDeleteAccount(w, r)
	}
	if r.Method == "PUT" {
		server.handleTransferAccount(w, r)
	}
}


func (server *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	account, err := server.Store.GetAccountByID(id)

	if err != nil {
		log.Println("Error getting account", err)
		Error := &Error{Message: "Error getting account"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}
	WriteJSON(w, http.StatusOK, account)
}


func (server *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	CreateAccountReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(CreateAccountReq); err != nil {
		log.Println("Error decoding request body", err)
		Error := &Error{Message: "Invalid request body"}
		WriteJSON(w, http.StatusBadRequest, Error)
		return
	}

	account := NewAccount(CreateAccountReq.FirstName, CreateAccountReq.LastName)
	if err := server.Store.CreateAccount(account); err != nil {
		log.Println("Error creating account", err)
		Error := &Error{Message: "Error creating account"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}

	CreateAccountResp := &CreateAccountResponse{ID: account.ID}
	WriteJSON(w, http.StatusOK, CreateAccountResp)
}


func (server *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	return
}


func (server *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) {
	return
}
