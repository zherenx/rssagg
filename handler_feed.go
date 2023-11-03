package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zherenx/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedWithUser(w http.ResponseWriter, r *http.Request, user *database.User) {
	type parameters struct {
		Name *string `json:"name"`
		Url  *string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// TODO: consider refactor this into a helper function
	if params.Name == nil {
		respondWithError(w, http.StatusBadRequest, "Name field is required but not found")
		return
	} else if *params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name shouldn't be empty")
		return
	}
	if params.Url == nil {
		respondWithError(w, http.StatusBadRequest, "Url field is required but not found")
		return
	} else if *params.Url == "" {
		// TODO: potential improvment, valid url check
		respondWithError(w, http.StatusBadRequest, "Url shouldn't be empty")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      *params.Name,
		Url:       *params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		// TODO: should i also log user id here?
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not get feeds: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
