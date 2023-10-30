package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zherenx/rssagg/internal/auth"
	"github.com/zherenx/rssagg/internal/database"
)

/*
Note: In go, the function signiture of an http handler can't change, but we
do want to pass into this function an addition piece of data, the apiConfig,
so what we do is to make this function a method
*/
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// this method is generated by sqlc
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not get API Key: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Could not get user with API Key: %v", err))
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
