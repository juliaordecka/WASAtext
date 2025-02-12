package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CreateGroupRequest struct {
	Name      string   `json:"name"`
	Usernames []string `json:"usernames"` // List of usernames to add
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get creator's ID from token
	var user User
	token := getToken(r.Header.Get("Authorization"))
	user.Id = token

	// Decode request
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate group name
	if req.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	// Get creator's username
	creatorUsername, err := rt.db.GetUsernameById(user.Id)
	if err != nil {
		http.Error(w, "Failed to get creator info", http.StatusInternalServerError)
		return
	}

	// Validate usernames
	for _, username := range req.Usernames {
		// Check if trying to add self
		if username == creatorUsername {
			http.Error(w, "Cannot add yourself to the group", http.StatusBadRequest)
			return
		}

		// Check if user exists
		_, err := rt.db.GetUserIdByUsername(username)
		if err != nil {
			http.Error(w, "User not found: "+username, http.StatusBadRequest)
			return
		}
	}

	// Only create group if all validations pass
	group, err := rt.db.CreateGroup(req.Name, user.Id)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add other users
	for _, username := range req.Usernames {
		err = rt.db.AddUserToGroup(username, group.ConversationId)
		if err != nil {
			// If we fail to add users, rollback the group creation
			if delErr := rt.db.DeleteGroup(group.ConversationId); delErr != nil {
				http.Error(w, "Failed to rollback group creation", http.StatusInternalServerError)
			}
			http.Error(w, "Failed to add users to group", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
