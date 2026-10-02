package api

import (
	"fmt"
	"net/http"
	"sort"
	"time"
	"github.com/AyuorusAguilar/chirpy/internal/database"
	"github.com/AyuorusAguilar/chirpy/internal/hashing"
	"github.com/AyuorusAguilar/chirpy/internal/jsonHandling"
	"github.com/AyuorusAguilar/chirpy/internal/state"
	"github.com/AyuorusAguilar/chirpy/internal/utils"
	"github.com/google/uuid"
)

func Ok(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func PostChirps(w http.ResponseWriter, r *http.Request, s *state.State) {
	/* Token validation */
	token, err := hashing.GetBearerToken(r.Header)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid Header: %v", err))
		return
	}
	userId, err := hashing.ValidateJWT(token, s.Secret)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid token: %v", err))
		return
	}

	/* Request checking */
	type expected struct {
		Body string `json:"body"`
	}
	var req expected
	err = utils.UnmarshalAt(r.Body, &req, w)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	/* Chirp validation */
	if len(req.Body) > 140 {
		jsonHandling.WriteJSONError(w, 401, "Chirp is too long!")
		return
	}
	cleanString := utils.ReplaceProfanity(req.Body)

	/* Chirp creation */
	createdChirp, err := s.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanString,
		UserID: userId,
	})
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("An error ocurred during database insertion, make sure to use a valid user!:\n%v", err))
		return
	}

	/* Response writing */
	jsonHandling.WriteJSONResponse(w, 201, map[string]string{}, map[string]string{
		"id":         createdChirp.ID.String(),
		"created_at": createdChirp.CreatedAt.String(),
		"updated_at": createdChirp.UpdatedAt.String(),
		"body":       createdChirp.Body,
		"user_id":    createdChirp.UserID.String()})
}
func DeleteChirp(w http.ResponseWriter, r *http.Request, s *state.State) {
	/* Token validation */
	token, err := hashing.GetBearerToken(r.Header)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid Header: %v", err))
		return
	}
	userId, err := hashing.ValidateJWT(token, s.Secret)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid token: %v", err))
		return
	}
	/* Chirp validation */
	reqChirp, err := uuid.Parse(r.PathValue("ChirpId"))
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid chirp uuid:%v", err))
		return
	}
	ActualChirp, err := s.Db.GetChirpById(r.Context(), reqChirp)
	if err != nil {
		jsonHandling.WriteJSONError(w, 404, fmt.Sprintf("Chirp not found:%v", err))
		return
	}
	if ActualChirp.UserID != userId {
		jsonHandling.WriteJSONError(w, 403, "Unauthorized request:")
		return
	}
	/* Deleting chirp */
	err = s.Db.DeleteChirp(r.Context(), ActualChirp.ID)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Error deleting the chirp!:%v", err))
		return
	}
	/* Response writing */
	w.WriteHeader(204)
}
func GetChirps(w http.ResponseWriter, r *http.Request, s *state.State) {
	var chirps []database.Chirp
	var err error
	/* Filtering */
	author := r.URL.Query().Get("author_id")
	if author != "" {
		id, err := uuid.Parse(author)
		if err != nil {
			jsonHandling.WriteJSONError(w, 400, fmt.Sprintf("Invalid Id:%v", err))
			return
		}
		chirps, err = s.Db.GetByChirpsByAuthor(r.Context(), id)
	} else {
		chirps, err = s.Db.GetChirps(r.Context())
		if err != nil {
			jsonHandling.WriteJSONError(w, 1, fmt.Sprintf("An error ocurred fetching data from the database, make sure there are entries to begin with!.\nError: %v", err))
			return
		}
	}
	/* Sorting */
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
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
	var requestedString string = r.PathValue("ChirpId")
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

func PostUser(w http.ResponseWriter, r *http.Request, s *state.State) {
	/* Request checking */
	type expected struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req expected
	err := utils.UnmarshalAt(r.Body, &req, w)
	if err != nil {
		return
	}

	/* Password hashing */
	hash, err := hashing.HashPass(req.Password)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, fmt.Sprintf("Invalid7unhashable password:\n%v", err))
	}

	createdUser, err := s.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		Hash:  hash,
	})
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, "Error creating user on the database!")
	}

	jsonHandling.WriteJSONResponse(w, 201, map[string]string{}, map[string]any{
		"id":         createdUser.ID.String(),
		"created_at": createdUser.CreatedAt.String(),
		"updated_at": createdUser.UpdatedAt.String(),
		"email":      createdUser.Email,
		"is_chirpy_red":createdUser.IsChirpyRed})
}
func PutUser(w http.ResponseWriter, r *http.Request, s *state.State) {
	/* Token validation */
	authToken, err := hashing.GetBearerToken(r.Header)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid Header:%v", err))
	}
	UserId, err := hashing.ValidateJWT(authToken, s.Secret)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid Auth Token (JWT):%v", err))
	}
	/* Request checking */
	type expected struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req expected
	err = utils.UnmarshalAt(r.Body, &req, w)
	if err != nil {
		return
	}
	/* Password hashing */
	hash, err := hashing.HashPass(req.Password)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid/unhashable password:\n%v", err))
	}
	/* Updating User */
	createdUser, err := s.Db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: UserId,
		Email: req.Email,
		Hash:  hash,
	})
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, "Error updating user on the database!")
	}

	jsonHandling.WriteJSONResponse(w, 200, map[string]string{}, map[string]any{
		"id":         createdUser.ID.String(),
		"created_at": createdUser.CreatedAt.String(),
		"updated_at": createdUser.UpdatedAt.String(),
		"email":      createdUser.Email,
		"is_chirpy_red":createdUser.IsChirpyRed})
}
func PostLogin(w http.ResponseWriter, r *http.Request, s *state.State) {
	/* Request checking */
	type expected struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	var req expected
	err := utils.UnmarshalAt(r.Body, &req, w)
	if err != nil {
		return
	}

	/* Check if exists */
	user, err := s.Db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, "Incorrect email or password")
		return
	}

	/* Validate password */
	if isValid, err := hashing.CheckPass(req.Password, user.Hash); !isValid || err != nil {
		jsonHandling.WriteJSONError(w, 401, "Incorrect email or password")
		return
	}
	/* TokenMaking */
	token, err := hashing.MakeJWT(user.ID, s.Secret)
	refreshToken := hashing.MakeRefreshToken()
	rToken, err := s.Db.CreateToken(r.Context(), database.CreateTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
	})
	/* Write response */
	jsonHandling.WriteJSONResponse(w, 200, map[string]string{}, map[string]any{
		"id":			user.ID.String(),
		"created_at":	user.CreatedAt.String(),
		"updated_at":	user.UpdatedAt.String(),
		"email":		user.Email,
		"token":		token,
		"is_chirpy_red":user.IsChirpyRed,
		"refresh_token":rToken.Token})
}
func RefreshToken(w http.ResponseWriter, r *http.Request, s *state.State){
	refreshToken, err := hashing.GetBearerToken(r.Header)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid header:%v", err))
	}
	ActualRefreshToken, err := s.Db.GetToken(r.Context(), refreshToken)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid token:%v", err))
	}
	if ActualRefreshToken.RevokedAt.Valid {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Expired/Revoked token:%v", err))
	}

	jwt, err := hashing.MakeJWT(ActualRefreshToken.UserID, s.Secret)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("An error ocurred creating the JWT:%v", err))
	}
	jsonHandling.WriteJSONResponse(w, 200, map[string]string{}, map[string]string{
		"token": jwt})
}
func RevokeToken(w http.ResponseWriter, r *http.Request, s *state.State){
	refreshToken, err := hashing.GetBearerToken(r.Header)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid header:%v", err))
	}
	
	err = s.Db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Error revoking token: %v", err))
	}
	w.WriteHeader(204)
}

func UpgradeUser(w http.ResponseWriter, r *http.Request, s *state.State){
	/* Request checking */
	type expected struct {
		Event	string `json:"event"`
		Data 	struct{
			User_id string `json:"user_id"`
		} `json:"data"`
	}
	var req expected
	err := utils.UnmarshalAt(r.Body, &req, w)
	if err != nil {
		return
	}
	/* Header checking */
	key, err := hashing.GetAPIKey(r.Header)
	if err != nil || key != s.PKey {
		jsonHandling.WriteJSONError(w, 401, fmt.Sprintf("Invalid Key:%v", err))
	}
	/* Event validation */
	if req.Event != "user.upgraded" {
		jsonHandling.WriteJSONError(w, 204, "Invalid event")
	}

	/* Id parsing */
	id, err := uuid.Parse(req.Data.User_id)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, fmt.Sprintf("Invalid ID: %v", err))
	}

	/* Upgrading */
	_, err = s.Db.UpgradeUser(r.Context(), id)
	if err != nil {
		jsonHandling.WriteJSONError(w, 400, fmt.Sprintf("An error ocurred upgrading the user:%v", err))
	}
	w.WriteHeader(204)
}


func MdwareAccessState(s *state.State, f func(w http.ResponseWriter, r *http.Request, s *state.State)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, s)
	}
}
