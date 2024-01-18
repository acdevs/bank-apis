package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

/* JWT Authentication */

func WithJWTAuth(next http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		token, err := ValidateJWT(tokenString)
		if err != nil {
			log.Println("Permission Denied", err)
			Error := &Error{Message: "Permission Denied"}
			WriteJSON(w, http.StatusUnauthorized, Error)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			log.Println("Permission Denied")
			Error := &Error{Message: "Permission Denied"}
			WriteJSON(w, http.StatusUnauthorized, Error)
			return
		}

		// panic(reflect.TypeOf(claims["number"])) ... actually good to know

		passedID := mux.Vars(r)["id"]
		account, err := store.GetAccountByID(passedID)
		if account.Number != int64(claims["number"].(float64)) || err != nil {
			log.Println("Permission Denied")
			Error := &Error{Message: "Permission Denied"}
			WriteJSON(w, http.StatusUnauthorized, Error)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := []byte("WeAReNeVEreVergEtTinGBaCkTOgEtHer")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}


func GenerateJWT(account *Account) (string, error) {
	secret := []byte("WeAReNeVEreVergEtTinGBaCkTOgEtHer")

	claims := &jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"id": account.ID,
		"number": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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

	router.HandleFunc("/api/login", server.handleLogin).Methods("POST")
	router.HandleFunc("/api/accounts", server.handleAccount).Methods("GET", "POST", "PUT")
	router.HandleFunc("/api/accounts/{id}", WithJWTAuth(server.handleGetAccountByID, server.Store)).Methods("GET")
	router.HandleFunc("/api/accounts/{id}", server.handleDeleteAccount).Methods("DELETE")
	router.HandleFunc("/api/accounts/transfer", server.handleTransferAccount).Methods("POST")
	
	log.Println("Listening on", server.ListenAddress)
	http.ListenAndServe(server.ListenAddress, router)
}


func (server *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		server.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		server.handleCreateAccount(w, r)
	}
	if r.Method == "PUT" {
		server.handleTransferAccount(w, r)
	}
}


func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (server *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	LoginReq := new(LoginRequest)
	json.NewDecoder(r.Body).Decode(LoginReq)

	account, err := server.Store.GetAccountByID(LoginReq.ID)
	if err != nil || !CheckPasswordHash(LoginReq.Password, account.Password) {
		log.Println("Incorrect Account ID or Password", err)
		Error := &Error{Message: "Incorrect Account ID or Password"}
		WriteJSON(w, http.StatusForbidden, Error)
		return
	}
	tokenString, err := GenerateJWT(account)
	if err != nil {
		log.Println("Error generating JWT", err)
		Error := &Error{Message: "Error generating JWT"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}

	// fmt.Println("tokenString : ", tokenString)
	// w.Header().Set("Authorization", tokenString)

	LoginResp := &LoginResponse{Token: tokenString}
	WriteJSON(w, http.StatusOK, LoginResp)
}


func (server *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	account, err := server.Store.GetAccountByID(id)

	if err != nil {
		log.Println("Account " + id + " not found", err)
		Error := &Error{Message: "Account " + id + " not found"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}
	WriteJSON(w, http.StatusOK, account)
}


func (server *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := server.Store.GetAccounts()

	if err != nil {
		log.Println("Error getting account", err)
		Error := &Error{Message: "Error getting account"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}
	WriteJSON(w, http.StatusOK, accounts)
}


func (server *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	req := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Println("Error decoding request body", err)
		Error := &Error{Message: "Invalid request body"}
		WriteJSON(w, http.StatusBadRequest, Error)
		return
	}

	account := NewAccount(req.FirstName, req.LastName, req.Password)
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
	id := mux.Vars(r)["id"]

	if err := server.Store.DeleteAccount(id); err != nil {
		log.Println("Error deleting account", err)
		Error := &Error{Message: "Error deleting account"}
		WriteJSON(w, http.StatusInternalServerError, Error)
		return
	}

	DeleteAccountResp := &DeleteAccountResponse{ID: id}
	WriteJSON(w, http.StatusOK, DeleteAccountResp)
}


func (server *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) {
	TransferAccountReq := new(TransferAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(TransferAccountReq); err != nil {
		log.Println("Error decoding request body", err)
		Error := &Error{Message: "Invalid request body"}
		WriteJSON(w, http.StatusBadRequest, Error)
		return
	}

	WriteJSON(w, http.StatusOK, TransferAccountReq)
}
