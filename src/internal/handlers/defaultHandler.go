package handlers

import (
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			http.Error(w, "Not implemented", http.StatusInternalServerError)
			// TODO: implement simple informational handler
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}
