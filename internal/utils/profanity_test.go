package utils

import (
	"testing"
)

func TestProfanity(t *testing.T) {
	type testEntry struct {
		input string
		wanted string
	}
	tests := []testEntry{
		{input: "kerfuffle", wanted: "****"},
		{input: "kerfuffle sharbert fornax", wanted: "**** **** ****"},
		{input: "Hola kerfuffle", wanted: "Hola ****"},
		{input: "Hola kerfufflekerfuffle", wanted: "Hola kerfufflekerfuffle"},
		{input: "I really need a kerfuffle to go to bed sooner, Fornax !", wanted: "I really need a **** to go to bed sooner, **** !"}}
	for _, test := range tests {
		result:= ReplaceProfanity(test.input)
		if result != test.wanted {
			t.Fatalf("Unexpected result, wanted %s got: %s\n", test.wanted, result)
		}
	}
}
