package hashing

import (
	"testing"
)

func TestHashing(t *testing.T)  {
	passwords := []string{
		"IsmuContrañero",
		"Passamela_boteia2",
		"contrasenasegura123"}

	for _, pass := range passwords {
		hash, err := HashPass(pass)
		if err != nil {
			t.Errorf("Unhashable pass? Error:\n\t%v", err)
		}
		if isValid, err := CheckPass(pass, hash); !isValid || err != nil {
			t.Errorf("???... Eh... UnUnhashable pass??? Error:\n\t%v", err)
		}
	}
}

func TestToken(t *testing.T) {
	tok := MakeRefreshToken()
	t.Logf("Hey!: %s\n", tok)
} 