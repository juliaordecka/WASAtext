package api

import (
	"encoding/json"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from Authorization header
	rt.baseLogger.Println("SendMessage endpoint called")
	var user User
	token := getToken(r.Header.Get("Authorization"))
	user.Id = token

	// Decode the request body
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		rt.baseLogger.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate message
	if message.Text == "" {
		http.Error(w, "Cannot send an empty message", http.StatusBadRequest)
		return
	}

	// Validate recipient
	if message.RecipientId == 0 {
		http.Error(w, "Recipient ID is required", http.StatusBadRequest)
		return
	}

	// Get or create conversation with recipient
	conversationId, err := rt.db.GetOrCreateDirectConversation(user.Id, message.RecipientId)
	if err != nil {
		http.Error(w, "Failed to handle conversation", http.StatusInternalServerError)
		return
	}

	// Set message metadata
	message.SenderId = user.Id
	message.ConversationId = conversationId
	message.SendTime = time.Now()
	message.Status = "Sent"

	// Store message in database
	dbMsg := message.ToDatabase()
	dbMsg, err = rt.db.CreateMessage(dbMsg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Convert back from database
	message.FromDatabase(dbMsg)

	// Update conversation's last message
	err = rt.db.UpdateLastMessage(message.MessageId, message.ConversationId)
	if err != nil {
		http.Error(w, "Failed to update conversation", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
