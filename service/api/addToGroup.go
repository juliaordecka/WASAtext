package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type AddToGroupRequest struct {
	Username string `json:"username"`
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get requester's ID from token
	var user User
	token := getToken(r.Header.Get("Authorization"))
	user.Id = token

	// Get group ID from URL
	groupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Check if requester is in group
	isInGroup, err := rt.db.IsUserInGroup(user.Id, groupId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isInGroup {
		http.Error(w, "Not authorized to add members", http.StatusForbidden)
		return
	}

	// Decode request
	var req AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add user to group directly using username
	err = rt.db.AddUserToGroup(req.Username, groupId)
	if err != nil {
		http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
