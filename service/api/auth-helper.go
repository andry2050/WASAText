package api

import (
	"net/http"
	"strings"
)

// extractBearer estrae l'ID utente dall'header Authorization della richiesta HTTP.
// Se l'header manca o non è formattato correttamente, restituisce una stringa vuota.
func extractBearer(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}