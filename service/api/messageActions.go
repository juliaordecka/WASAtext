package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
	
)

type CommentRequest struct {
    Emoji string `json:"emoji"`
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    messageId, err := strconv.Atoi(ps.ByName("message_id"))
    if err != nil {
        http.Error(w, "Invalid message ID", http.StatusBadRequest)
        return
    }

    // Check if user owns the message
    isOwner, err := rt.db.IsMessageOwner(messageId, user.Id)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if !isOwner {
        http.Error(w, "Not authorized to delete this message", http.StatusForbidden)
        return
    }

    err = rt.db.DeleteMessage(messageId, user.Id)
    if err != nil {
        http.Error(w, "Failed to delete message", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    messageId, err := strconv.Atoi(ps.ByName("message_id"))
    if err != nil {
        http.Error(w, "Invalid message ID", http.StatusBadRequest)
        return
    }

    // Check if user is trying to comment their own message
    isOwner, err := rt.db.IsMessageOwner(messageId, user.Id)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    if isOwner {
        http.Error(w, "Cannot comment on your own message", http.StatusBadRequest)
        return
    }

    var req CommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = rt.db.CommentMessage(messageId, user.Id, req.Emoji)
    if err != nil {
        http.Error(w, "Failed to comment message", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

    messageId, err := strconv.Atoi(ps.ByName("message_id"))
    if err != nil {
        http.Error(w, "Invalid message ID", http.StatusBadRequest)
        return
    }

    err = rt.db.UncommentMessage(messageId, user.Id)
    if err != nil {
        http.Error(w, "Failed to remove comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

