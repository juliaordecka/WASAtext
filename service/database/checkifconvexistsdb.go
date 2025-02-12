package database

// import (
//	"database/sql"
// )

func (db *appdbimpl) CheckIfConversationExists(conversationId int) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM conversations WHERE ConversationId = ?)", conversationId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
