package handlers

import "net/http"

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	case http.MethodPost:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	case http.MethodDelete:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}
