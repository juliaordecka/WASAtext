package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) CreateConversation(userId uint64, conversationId int) (Conversation, error) {
	// First create the conversation
	res, err := db.c.Exec("INSERT INTO conversations (GroupId, LastMessageId) VALUES (?, 0)", conversationId)
	if err != nil {
		return Conversation{}, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return Conversation{}, err
	}

	// Then add the user as a participant
	_, err = db.c.Exec("INSERT INTO participants (ConversationId, UserId) VALUES (?, ?)", lastInsertID, userId)
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{
		ConversationId: int(lastInsertID),
		GroupId:        conversationId,
		LastMessageId:  0, // Initializing with no messages
	}, nil
}

func (db *appdbimpl) UpdateLastMessage(messageId int, conversationId int) error {
	_, err := db.c.Exec("UPDATE conversations SET LastMessageId = ? WHERE ConversationId = ?", messageId, conversationId)
	return err
}

func (db *appdbimpl) GetOrCreateDirectConversation(userId, recipientId uint64) (int, error) {
	// First try to find existing conversation
	var conversationId int
	err := db.c.QueryRow(`
        SELECT c.ConversationId 
        FROM conversations c
        JOIN participants p1 ON c.ConversationId = p1.ConversationId
        JOIN participants p2 ON c.ConversationId = p2.ConversationId
        WHERE p1.UserId = ? AND p2.UserId = ? AND c.GroupId = 0
        AND NOT EXISTS (
            SELECT 1 FROM participants p3 
            WHERE p3.ConversationId = c.ConversationId 
            AND p3.UserId NOT IN (?, ?)
        )
    `, userId, recipientId, userId, recipientId).Scan(&conversationId)

	if err == nil {
		return conversationId, nil
	}

	// If not found, create new conversation
	tx, err := db.c.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("Transaction rollback failed: %v\n", err)
		}
	}()

	// Create conversation
	result, err := tx.Exec("INSERT INTO conversations (GroupId, LastMessageId) VALUES (0, 0)")
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	conversationId = int(lastId)

	// Add participants
	_, err = tx.Exec("INSERT INTO participants (ConversationId, UserId) VALUES (?, ?)", conversationId, userId)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO participants (ConversationId, UserId) VALUES (?, ?)", conversationId, recipientId)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return conversationId, nil
}
