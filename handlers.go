package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Mood structure for a mood
type Mood struct {
	MoodID int
	Name   string
}

func (ms *myServer) handleListMoods(w http.ResponseWriter, r *http.Request) {

	// Return moods as slice of moods
	var moods []Mood

	for _, mood := range ms.moods {
		moods = append(moods, mood)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(moods)
}

func (ms *myServer) handleGetUserMood(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	mood, ok := ms.userMoods[userID]
	if !ok {
		mood = Mood{}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if mood == (Mood{}) {
		json.NewEncoder(w).Encode("null")
		return
	}

	json.NewEncoder(w).Encode(ms.userMoods[userID])
}

func (ms *myServer) handlePutUserMood(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Grab the user whose status we're updating
	userID := chi.URLParam(r, "userId")
	// Grab the user from context, assume string
	contextUserID := ctx.Value(currentUserID).(string)

	// UserIDs don't match; unauthorized
	if userID != contextUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var mood Mood
	if err := json.NewDecoder(r.Body).Decode(&mood); err != nil {
		msg := "failed to decode body"
		ms.Logger.WithError(err).Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	// Make sure valid mood was provided
	if _, ok := ms.moods[mood.MoodID]; !ok {
		msg := "invalid mood provided"
		ms.Logger.WithError(errors.New(msg)).Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	ms.userMoods[userID] = mood

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("status for userID: %s was updated to: %s", userID, mood.Name)
	json.NewEncoder(w).Encode(msg)
}
