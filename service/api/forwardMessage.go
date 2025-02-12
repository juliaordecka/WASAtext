package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

type ForwardMessageRequest struct {
    ConversationId int `json:"conversationId"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Get user ID from token
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    // Get message ID from URL
    messageId, err := strconv.Atoi(ps.ByName("message_id"))
    if err != nil {
        http.Error(w, "Invalid message ID", http.StatusBadRequest)
        return
    }

    // Decode request
    var req ForwardMessageRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Check if user is part of target conversation
    isInConv, err := rt.db.IsUserInGroup(user.Id, req.ConversationId)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if !isInConv {
        http.Error(w, "Not authorized to forward to this conversation", http.StatusForbidden)
        return
    }

    // Forward message
    forwardedMsg, err := rt.db.ForwardMessage(messageId, user.Id, req.ConversationId)
    if err != nil {
        http.Error(w, "Failed to forward message", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(forwardedMsg); err != nil {
    http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    return
}
}

