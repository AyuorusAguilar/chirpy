package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/AyuorusAguilar/chirpy/internal/jsonHandling"
)
var profanity = []string{"kerfuffle", "sharbert", "fornax"}
func ReplaceProfanity(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		for _, profWord := range profanity {
			if strings.ToLower(word) == profWord {
				words[i] = "****"
			}
		}
	}
	result := strings.Join(words, " ")
	return result
}

func UnmarshalAt(content io.Reader, container any, w http.ResponseWriter) error {
	byt, err := io.ReadAll(content)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, "Couldn't read the request!")
		return err
	}
	json.Unmarshal(byt, container)
	return nil
}