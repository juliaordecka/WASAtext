package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    convId, err := strconv.Atoi(ps.ByName("conversation_id"))
    if err != nil {
        http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
        return
    }

    // Check if user is participant
    isMember, err := rt.db.IsUserInGroup(user.Id, convId)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if !isMember {
        http.Error(w, "Not authorized to view conversation", http.StatusForbidden)
        return
    }

    conversation, err := rt.db.GetConversationDetails(convId, user.Id)
    if err != nil {
        http.Error(w, "Failed to get conversation", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(conversation); err != nil {
    http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    return
}
}

