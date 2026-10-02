package hashing

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var params = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

func HashPass(s string) (string, error) {
	return argon2id.CreateHash(s, params)
}

func CheckPass(pass, hash string) (bool, error) {
	isValid, err := argon2id.ComparePasswordAndHash(pass, hash)
	return isValid, err
}

func MakeJWT(id uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   id.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	a := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &a,
		func(t *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.UUID{}, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(subject)
}

func GetBearerToken(h http.Header) (string, error) {
	s := h.Get("Authorization")
	trimed, hadIt := strings.CutPrefix(s, "Bearer ")
	if !hadIt {
		return "", fmt.Errorf("Invalid Token")
	}
	return trimed, nil
}

func MakeRefreshToken() string {
	cont := make([]byte, 32)
	rand.Read(cont)
	return hex.EncodeToString(cont)
}

func GetAPIKey(h http.Header) (string, error)  {
	s := h.Get("Authorization")
	trimed, hadIt := strings.CutPrefix(s, "ApiKey ")
	if !hadIt {
		return "", fmt.Errorf("Invalid Key")
	}
	return trimed, nil
}