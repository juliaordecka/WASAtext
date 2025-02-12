package api

import (
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
        http.Error(w, "You are not a member of this group", http.StatusForbidden)
        return
    }

    // Leave group
    err = rt.db.LeaveGroup(user.Id, groupId)
    if err != nil {
        http.Error(w, "Failed to leave group", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}