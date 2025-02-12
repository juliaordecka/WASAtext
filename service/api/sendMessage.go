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

 // Handle conversation lookup by name
 if message.ConversationName != "" {
	convId, err := rt.db.GetConversationIdByName(message.ConversationName)
	if err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}
	message.ConversationId = convId
}

// Handle recipient lookup by username
if message.RecipientUsername != "" {
	recipientId, err := rt.db.GetRecipientIdByUsername(message.RecipientUsername)
	if err != nil {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	}
	message.RecipientId = recipientId
}

// Check if it's a group message or direct message
if message.ConversationId != 0 {
	// It's a message to an existing conversation (group or direct)
	isParticipant, err := rt.db.IsUserInGroup(user.Id, message.ConversationId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Not authorized to send message to this conversation", http.StatusForbidden)
		return
	}
} else if message.RecipientId != 0 {
	// It's a new direct message
	newConvId, err := rt.db.GetOrCreateDirectConversation(user.Id, message.RecipientId)
	if err != nil {
		http.Error(w, "Failed to handle conversation", http.StatusInternalServerError)
		return
	}
	message.ConversationId = newConvId
} else {
	http.Error(w, "Must specify either conversation name/id or recipient username/id", http.StatusBadRequest)
	return
}


	// Set message metadata
    message.SenderId = user.Id
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
