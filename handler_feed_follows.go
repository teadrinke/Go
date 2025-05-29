package main

import (
	"encoding/json"
	"fmt"
	
	
	
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/teadrinke/Go/internal/database"
	
)



func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err!= nil {
		
		respondWithError(w,400, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err!=nil{
		
		respondWithError(w,400, fmt.Sprintf("Couldnot create feed follow: %v", err))
		return
	}
	respondWithJson(w, 201, feedFollow)
	// respondWithJson(w, 200, user)
}

func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollow, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)

	if err!=nil{
		respondWithError(w,400, fmt.Sprintf("Couldnot get feed follow: %v", err))
		return
	}

	feedFollowList := []database.FeedFollow{}
	for _, feed := range feedFollow {
		feedFollowList = append(feedFollowList, feed)
	}
	respondWithJson(w, 201, feedFollowList)
	// respondWithJson(w, 200, user)
}
