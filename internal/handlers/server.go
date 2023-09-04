package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()
	notes := router.PathPrefix("/notes").Subrouter()
	notes.Use(h.UserCheck)
	notes.HandleFunc("", h.Create).Methods(http.MethodPost)
	notes.HandleFunc("/{id}", h.Read).Methods(http.MethodGet)
	notes.HandleFunc("", h.ReadAll).Methods(http.MethodGet)
	notes.HandleFunc("", h.Update).Queries("id", "{id}").Methods(http.MethodPatch)
	notes.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
	registration := router.PathPrefix("/registration").Subrouter()
	registration.HandleFunc("", h.UserRegistration).Methods(http.MethodPost)
	return router
}
