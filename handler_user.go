package main

import (
	"encoding/json"
	"fmt"
	
	
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/teadrinke/Go/internal/database"
	
)



func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err!= nil {
		
		respondWithError(w,400, fmt.Sprintf("error parsing json: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err!=nil{
		
		respondWithError(w,400, fmt.Sprintf("Couldnot create user: %v", err))
		return
	}
	respondWithJson(w, 201, user)
	// respondWithJson(w, 200, user)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User){
	// Extract user ID from the URL path
	respondWithJson(w, 200, user)
}