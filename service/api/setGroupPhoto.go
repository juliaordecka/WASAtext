package api

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

type PhotoRequest struct {
    Photo string `json:"photo"` // base64 encoded photo
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    groupId, err := strconv.Atoi(ps.ByName("group_id"))
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    // Check if user is in group
    isMember, err := rt.db.IsUserInGroup(user.Id, groupId)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if !isMember {
        http.Error(w, "Not authorized to modify group", http.StatusForbidden)
        return
    }

    var req PhotoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate base64 photo
    if _, err := base64.StdEncoding.DecodeString(req.Photo); err != nil {
        http.Error(w, "Invalid photo format", http.StatusBadRequest)
        return
    }

    err = rt.db.SetGroupPhoto(groupId, req.Photo)
    if err != nil {
        http.Error(w, "Failed to set photo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

