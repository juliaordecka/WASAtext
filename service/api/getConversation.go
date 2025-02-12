package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	token := getToken(r.Header.Get("Authorization"))
	user.Id = token

	rt.baseLogger.Printf("Getting conversation. User ID: %d", user.Id)

	convId, err := strconv.Atoi(ps.ByName("conversation_id"))
	if err != nil {
		rt.baseLogger.Printf("Invalid conversation ID: %v", err)
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	rt.baseLogger.Printf("Conversation ID: %d", convId)

	// Check if user is participant
	isMember, err := rt.db.IsUserInGroup(user.Id, convId)
	if err != nil {
		rt.baseLogger.Printf("Error checking membership: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		rt.baseLogger.Printf("User %d is not a member of conversation %d", user.Id, convId)
		http.Error(w, "Not authorized to view conversation", http.StatusForbidden)
		return
	}

	rt.baseLogger.Printf("Getting conversation details")
	conversation, err := rt.db.GetConversationDetails(convId, user.Id)
	if err != nil {
		rt.baseLogger.Printf("Error getting conversation details: %v", err)
		http.Error(w, "Failed to get conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		rt.baseLogger.Printf("Error encoding conversation: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
