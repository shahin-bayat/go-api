package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	// go get github.com/gorilla/mux
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin)).Methods("POST")
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", withJwtAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store)).Methods("GET")
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer)).Methods("POST")

	log.Println("Starting server on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var loginRequest = new(LoginRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(loginRequest); err != nil {
		return err
	}

	account, err := s.store.GetAccountByIban(loginRequest.Iban)
	if err != nil {
		return err
	}

	if !account.validatePassword(loginRequest.Password) {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not authorized"})
	}

	token, err := createJwt(account)
	if err != nil {
		return err
	}
	w.Header().Set("x-jwt-token", token)
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	// or
	// createAccountRequest := CreateAccountRequest{}
	// Decode accepts a pointer to a value
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(req); err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)

}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	transferRequest := new(TransferRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(transferRequest); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, transferRequest)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func withJwtAuth(handleFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJwt(tokenString)
		if err != nil || !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		userId, err := getId(r)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		account, err := s.GetAccountById(userId)
		if err != nil {
			// improvement: create a custom error func
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			// you can still log the error to some logging service (elastic) to knw=ow what went wrong
			return
		}

		// cast the claims to jwt.MapClaims
		claims := token.Claims.(jwt.MapClaims)

		if account.IBAN != claims["iban"] {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
			return
		}

		handleFunc(w, r)
	}
}

func validateJwt(tokenString string) (*jwt.Token, error) {
	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#section-readme
	// you should use export JWT_SECRET=your_secret in your terminal
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

func createJwt(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"iban":      account.IBAN,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getId(r *http.Request) (int, error) {
	var idStr = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idStr)
	}
	return id, nil
}
