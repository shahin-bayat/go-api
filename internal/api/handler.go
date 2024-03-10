package api

import (
	"encoding/json"
	"net/http"

	"github.com/shahin-bayat/go-api/internal/model"
	"github.com/shahin-bayat/go-api/internal/util"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var loginRequest = new(model.LoginRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(loginRequest); err != nil {
		return err
	}

	account, err := s.store.GetAccountByIban(loginRequest.Iban)
	if err != nil {
		return err
	}

	if !account.ValidatePassword(loginRequest.Password) {
		return util.WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "not "})
	}

	token, err := createJwt(account)
	if err != nil {
		return err
	}
	w.Header().Set("x-jwt-token", token)
	return util.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return util.WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := util.GetId(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(model.CreateAccountRequest)
	// or
	// createAccountRequest := CreateAccountRequest{}
	// Decode accepts a pointer to a value
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(req); err != nil {
		return err
	}

	account, err := model.NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusCreated, account)

}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := util.GetId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return util.WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	transferRequest := new(model.TransferRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(transferRequest); err != nil {
		return err
	}
	return util.WriteJSON(w, http.StatusOK, transferRequest)
}
