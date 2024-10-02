package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// DecodeRequest decodes the JSON request body into the provided struct
func decodeRequest[T any](w http.ResponseWriter, r *http.Request, request T) error {
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return err
	}
	return nil
}

// respondWithJSON sends a JSON response with the provided status code and body
func respondWithJSON[T any](w http.ResponseWriter, responseCode int, body T) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(&body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing HTTP response: %s", err)
		w.WriteHeader(500)
		return
	}
}

// respondWithError sends a JSON error response with the provided status code and error message
func respondWithError(w http.ResponseWriter, responseCode int, errorMsg string) {
	responseStruct := struct {
		Error string `json:"error"`
	}{
		Error: errorMsg,
	}
	data, err := json.Marshal(&responseStruct)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing HTTP response: %s", err)
		w.WriteHeader(500)
		return
	}
}
