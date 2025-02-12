package api

import (
    "encoding/json"
    "net/http"
 //   "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    rt.baseLogger.Printf("Getting conversations for user %d", user.Id)  // Add logging

    conversations, err := rt.db.GetConversations(user.Id)
    if err != nil {
        rt.baseLogger.Printf("Error getting conversations: %v", err)  // Add error logging
        http.Error(w, "Failed to get conversations", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(conversations); err != nil {
    http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    return
}
}