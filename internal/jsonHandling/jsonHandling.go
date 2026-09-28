package jsonHandling

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, code int, headers map[string]string, payload any) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		w.Header().Set(key, value)
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshalling payload")
	}
	w.Write(body)
	return nil
}

func WriteJSONError(w http.ResponseWriter, code int, msg string) error {
	return WriteJSONResponse(w, code, map[string]string{}, map[string]string{"error": msg})
}