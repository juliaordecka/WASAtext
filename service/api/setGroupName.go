package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

type SetGroupNameRequest struct {
    Name string `json:"name"`
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Get user ID from token
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    // Get group ID from URL
    groupId, err := strconv.Atoi(ps.ByName("group_id"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    // Check if user is in group
    isInGroup, err := rt.db.IsUserInGroup(user.Id, groupId)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if !isInGroup {
        http.Error(w, "Not authorized to change group name", http.StatusForbidden)
        return
    }

    // Decode request
    var req SetGroupNameRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate new name
    if req.Name == "" {
        http.Error(w, "Group name cannot be empty", http.StatusBadRequest)
        return
    }

    // Set new name
    err = rt.db.SetGroupName(groupId, req.Name)
    if err != nil {
        http.Error(w, "Failed to update group name", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

