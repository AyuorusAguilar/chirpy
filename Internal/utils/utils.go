package utils

import (
	"strings"
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