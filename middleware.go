package main

import (
	"fmt"
	"net/http"

	"github.com/zherenx/rssagg/internal/auth"
	"github.com/zherenx/rssagg/internal/database"
)

/*
Note:
I came across quite a few articles/tutorials/posts on the
handler middleware/adapters/wrappers topic, which most of them
name it xxxHandler while the underlying thing is obviously a
HandlerFunc, I understand that HandlerFunc is kind of a syntactic
sugar thing to Handler in go, but why couldn't or shouldn't we
call them xxxHandlerFunc, I'm confused..
*/
type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, *database.User)

func (apiCfg *apiConfig) AuthMiddleware(h AuthenticatedHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not get API Key: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Could not get user with API Key: %v", err))
			return
		}

		h(w, r, &user)
	}
}
