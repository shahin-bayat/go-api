package api

import (
	"net/http"

	"github.com/shahin-bayat/go-api/internal/store"
	"github.com/shahin-bayat/go-api/internal/util"
)

type APIServer struct {
	listenAddr string
	store      store.Store
}

func NewAPIServer(listenAddr string, store store.Store) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error here
			util.WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
