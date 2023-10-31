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
		Name *string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// I changed the Name field in parameters struct to *string type to enable
	// nil check/handle (i.e. "name" field not exist case), and differentiate
	// nil case and empty string case
	if params.Name == nil {
		respondWithError(w, http.StatusBadRequest, "Name is not found")
		return
	} else if *params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}

	// Note: this method is generated by sqlc
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      *params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
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

	// TODO:
	// I think there is a potential improvement here, I think this function
	// shouldn't return the API Key, and we should also consider hiding the
	// API Key (e.g. storing only the hash value in the database)
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
