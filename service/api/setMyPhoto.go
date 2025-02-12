package api

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
//    "strconv"
    "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
    "github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
    token := getToken(r.Header.Get("Authorization"))
    user.Id = token

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

    err := rt.db.SetUserPhoto(user.Id, req.Photo)
    if err != nil {
        http.Error(w, "Failed to set photo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
