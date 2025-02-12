package database

import (
	"database/sql"
	"fmt"
	"errors"
)

func (db *appdbimpl) CreateGroup(name string, creatorId uint64) (Conversation, error) {
	tx, err := db.c.Begin()
	if err != nil {
		return Conversation{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("Transaction rollback failed: %v\n", err)
		}
	}()

	// Create the conversation/group
	result, err := tx.Exec("INSERT INTO conversations (GroupId, LastMessageId, Name) VALUES (1, 0, ?)", name)
	if err != nil {
		return Conversation{}, err
	}

	conversationId, err := result.LastInsertId()
	if err != nil {
		return Conversation{}, err
	}

	// Add creator to participants
	_, err = tx.Exec("INSERT INTO participants (ConversationId, UserId) VALUES (?, ?)",
		conversationId, creatorId)
	if err != nil {
		return Conversation{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{
		ConversationId: int(conversationId),
		GroupId:        1, // 1 indicates it's a group
		LastMessageId:  0,
		Name: 			name,
	}, nil
}

func (db *appdbimpl) AddUserToGroup(username string, groupId int) error {
	// First get the user ID from username
	var userId uint64
	err := db.c.QueryRow("SELECT Id FROM users WHERE Username = ?", username).Scan(&userId)
	if err != nil {
		return err // Return error if user doesn't exist
	}

	// Check if user is already in group
	var exists bool
	err = db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM participants 
            WHERE UserId = ? AND ConversationId = ?
        )`, userId, groupId).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // User is already in group
	}

	// Add user to group
	_, err = db.c.Exec("INSERT INTO participants (ConversationId, UserId) VALUES (?, ?)",
		groupId, userId)
	return err
}

func (db *appdbimpl) GetUserIdByUsername(username string) (uint64, error) {
	var userId uint64
	err := db.c.QueryRow("SELECT Id FROM users WHERE Username = ?", username).Scan(&userId)
	return userId, err
}

func (db *appdbimpl) IsUserInGroup(userId uint64, groupId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM participants 
            WHERE UserId = ? AND ConversationId = ?
        )`, userId, groupId).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) DeleteGroup(groupId int) error {
	// Start transaction
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("Transaction rollback failed: %v\n", err)
		}
	}()

	// Delete from participants first (due to foreign key)
	_, err = tx.Exec("DELETE FROM participants WHERE ConversationId = ?", groupId)
	if err != nil {
		return err
	}

	// Delete the conversation
	_, err = tx.Exec("DELETE FROM conversations WHERE ConversationId = ?", groupId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *appdbimpl) LeaveGroup(userId uint64, groupId int) error {
    result, err := db.c.Exec("DELETE FROM participants WHERE UserId = ? AND ConversationId = ?", 
        userId, groupId)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return errors.New("user not found in group")
    }
    return nil
}

func (db *appdbimpl) SetGroupName(groupId int, newName string) error {
    result, err := db.c.Exec("UPDATE conversations SET Name = ? WHERE ConversationId = ? AND GroupId = 1", 
        newName, groupId)
    if err != nil {
        return err
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return errors.New("group not found")
    }
    return nil
}
