package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (ms *myServer) handleListUserStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ms.userStatuses)
}

func (ms *myServer) handleGetUserStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{userID: ms.userStatuses[userID]})
}

func (ms *myServer) handlePutUserStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Grab the user whose status we're updating
	userID := chi.URLParam(r, "userId")
	// Grab the user from context, assume string
	contextUser := ctx.Value("currentUser").(string)

	// UserIDs don't match; unauthorized
	if userID != contextUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload map[string]string

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		msg := "failed to decode body"
		ms.Logger.WithError(err).Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	newStatus, ok := payload["status"]
	if !ok {
		msg := "no status provided"
		ms.Logger.WithError(errors.New(msg)).Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	ms.userStatuses[userID] = newStatus

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("status for userID: %s was updated to: %s", userID, newStatus)
	json.NewEncoder(w).Encode(msg)
}
