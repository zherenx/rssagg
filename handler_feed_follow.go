package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/zherenx/rssagg/internal/database"
)

func (apiCfg *apiConfig) HandlerCreateFeedFollowForUser(w http.ResponseWriter, r *http.Request, user *database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	// TODO: potential improvement as the json parsing code are repetitive
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// TODO: FeedId nil check to provide more useful error message?

	feedFollow, err := apiCfg.DB.CreateFeedFollowForUser(r.Context(), database.CreateFeedFollowForUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) HandlerGetFeedFollowsOfUser(w http.ResponseWriter, r *http.Request, user *database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollowsOfUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not get feed follows: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) HandlerDeleteFeedFollowForUser(w http.ResponseWriter, r *http.Request, user *database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse feed follow id: %v", err))
		return
	}

	// TODO: should we check if the feedFollowId exist in the database, or check
	// afterward if the delete query actually delete something? if the entry does
	// not exist, should we still return 200?

	err = apiCfg.DB.DeleteFeedFollowForUser(r.Context(), database.DeleteFeedFollowForUserParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unfollow was not successful: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
