package handler

import (
	"encoding/json"
	"mime"
	"net/http"
)

var allowedContent = map[string]bool{
	"application/json":   true,
	"application/x-gzip": true,
}

func validateContentType(w http.ResponseWriter, r *http.Request) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !allowedContent[mediaType] {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	dec := json.NewDecoder(r.Body)
	//dec.DisallowUnknownFields() // запрещать поля, которых нет в структуре
	if err := dec.Decode(&dst); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return true
	}
	return false
}
