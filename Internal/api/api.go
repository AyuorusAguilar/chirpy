package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AyuorusAguilar/chirpy/Internal/jsonHandling"
	"github.com/AyuorusAguilar/chirpy/Internal/state"
	"github.com/AyuorusAguilar/chirpy/Internal/utils"
	"github.com/AyuorusAguilar/chirpy/internal/database"
	"github.com/google/uuid"
)

func Ok(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func PostChirps(w http.ResponseWriter, r *http.Request, s *state.State) {
	type expected struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	var req expected

	byt, err := io.ReadAll(r.Body)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, "Couldn't read the request!")
		return
	}
	json.Unmarshal(byt, &req)

	if len(req.Body) > 140 {
		jsonHandling.WriteJSONError(w, 400, "Chirp is too long!")
		return
	}

	cleanString := utils.ReplaceProfanity(req.Body)

	createdChirp, err := s.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanString,
		UserID: req.User_id,
	})
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, fmt.Sprintf("An error ocurred during database insertion, make sure to use a valid user!:\n%v", err))
		return
	}

	jsonHandling.WriteJSONResponse(w, 201, map[string]string{}, map[string]string{
		"id":         createdChirp.ID.String(),
		"created_at": createdChirp.CreatedAt.String(),
		"updated_at": createdChirp.UpdatedAt.String(),
		"body":       createdChirp.Body,
		"user_id":    createdChirp.UserID.String()})
}

func GetChirps(w http.ResponseWriter, r *http.Request, s *state.State) {
	chirps, err := s.Db.GetChirps(r.Context())
	if err != nil {
		jsonHandling.WriteJSONError(w, 1, fmt.Sprintf("An error ocurred fetching data from the database, make sure there are entries to begin with!.\nError: %v", err))
		return
	}
	if len(chirps) < 1 {
		fmt.Printf("Error!!!")
		jsonHandling.WriteJSONError(w, 1, "The database is empty, send some chirps to begin with!")
		return
	}
	var chirpsMap = make([]map[string]any, len(chirps))
	for i, chirp := range chirps {
		chirpsMap[i] = make(map[string]any)
		chirpsMap[i]["id"] = chirp.ID.String()
		chirpsMap[i]["created_at"] = chirp.CreatedAt
		chirpsMap[i]["updated_at"] = chirp.UpdatedAt
		chirpsMap[i]["body"] = chirp.Body
		chirpsMap[i]["user_id"] = chirp.UserID.String()
	}

	jsonHandling.WriteJSONResponse(w, 200, map[string]string{}, chirpsMap)
}

func GetChirpsById(w http.ResponseWriter, r *http.Request, s *state.State) {
	var requestedString string = r.PathValue("GetChirpsById")
	requestedId, err := uuid.Parse(requestedString)
	if err != nil {
		jsonHandling.WriteJSONError(w, 404, fmt.Sprintf("Invalid requested Id:%v", requestedString))
		return
	}
	returnedChirp, err := s.Db.GetChirpById(r.Context(), requestedId)
	if err != nil {
		jsonHandling.WriteJSONError(w, 404, fmt.Sprintf(
			"An error ocurred fetching data from the database, make sure there are entries to begin with!.\nError: %v", err))
		return
	}

	jsonHandling.WriteJSONResponse(w, 200, map[string]string{}, map[string]string{
		"id":         returnedChirp.ID.String(),
		"created_at": returnedChirp.CreatedAt.String(),
		"updated_at": returnedChirp.UpdatedAt.String(),
		"body":       returnedChirp.Body,
		"user_id":    returnedChirp.UserID.String()})
}

func CreateUser(w http.ResponseWriter, r *http.Request, s *state.State) {
	type expected struct {
		Email string `json:"email"`
	}
	var req expected

	byt, err := io.ReadAll(r.Body)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, "Couldn't read the request!")
		return
	}
	json.Unmarshal(byt, &req)

	fmt.Println(req.Email)

	createdUser, err := s.Db.CreateUser(r.Context(), req.Email)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, "Error creating user on the database!")
	}

	jsonHandling.WriteJSONResponse(w, 201, map[string]string{}, map[string]string{
		"id":         createdUser.ID.String(),
		"created_at": createdUser.CreatedAt.String(),
		"updated_at": createdUser.UpdatedAt.String(),
		"email":      createdUser.Email})
}

func MdwareAccessState(s *state.State, f func(w http.ResponseWriter, r *http.Request, s *state.State)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s)
	}
}
