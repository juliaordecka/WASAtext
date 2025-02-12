package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	rt.baseLogger.Println("searchUsers endpoint called")

	// Get query parameter
	query := r.URL.Query().Get("username")
	rt.baseLogger.Printf("Search query: %s", query)

	if query == "" {
		rt.baseLogger.Println("Empty search query")
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	rt.baseLogger.Printf("Calling database SearchUsers with query: %s", query)
	users, err := rt.db.SearchUsers(query)
	if err != nil {
		rt.baseLogger.Printf("Database error in SearchUsers: %v", err)
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Printf("Found %d users", len(users))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		rt.baseLogger.Printf("Error encoding users: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
