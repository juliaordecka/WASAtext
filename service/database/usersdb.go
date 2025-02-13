package database

import (
	"database/sql"
	"log"
	"time"
)

func (db *appdbimpl) SetUserPhoto(userId uint64, photoData string) error {
	_, err := db.c.Exec("UPDATE users SET ProfilePhoto = ? WHERE Id = ?",
		photoData, userId)
	return err
}

func (db *appdbimpl) SetGroupPhoto(groupId int, photoData string) error {
	_, err := db.c.Exec("UPDATE conversations SET GroupPhoto = ? WHERE ConversationId = ? AND GroupId = 1",
		photoData, groupId)
	return err
}

func (db *appdbimpl) GetConversations(userId uint64) ([]ConversationPreview, error) {
	query := `
        SELECT DISTINCT  -- Add DISTINCT to prevent duplicates
            c.ConversationId,
            CASE 
                WHEN c.GroupId = 1 THEN c.Name
                WHEN u.Username IS NOT NULL THEN u.Username
                ELSE 'Unknown'
            END as Name,
            c.GroupPhoto as Photo,
            m.SendTime as LastMessageTime,
            m.Text as LastMessageText,
            CASE WHEN m.Photo IS NOT NULL AND m.Photo != '' THEN 1 ELSE 0 END as IsPhoto,  -- Fix photo check
            CASE WHEN c.GroupId = 1 THEN 1 ELSE 0 END as IsGroup
        FROM conversations c
        INNER JOIN participants p ON c.ConversationId = p.ConversationId AND p.UserId = ?
        LEFT JOIN messages m ON c.LastMessageId = m.MessageId
        LEFT JOIN participants p2 ON c.ConversationId = p2.ConversationId AND p2.UserId != ?
        LEFT JOIN users u ON p2.UserId = u.Id
        GROUP BY c.ConversationId  -- Group by to avoid duplicates
        ORDER BY COALESCE(m.SendTime, '1970-01-01') DESC`

	rows, err := db.c.Query(query, userId, userId)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var conversations []ConversationPreview
	for rows.Next() {
		var conv ConversationPreview
		var photoNull sql.NullString
		var textNull sql.NullString
		var timeNull sql.NullTime

		err := rows.Scan(
			&conv.ConversationId,
			&conv.Name,
			&photoNull,
			&timeNull,
			&textNull,
			&conv.IsPhoto,
			&conv.IsGroup,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}

		if photoNull.Valid {
			conv.Photo = photoNull.String
		}
		if textNull.Valid {
			conv.LastMessageText = textNull.String
		}
		if timeNull.Valid {
			conv.LastMessageTime = timeNull.Time
		} else {
			conv.LastMessageTime = time.Now()
		}

		conversations = append(conversations, conv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return conversations, nil
}

func (db *appdbimpl) GetConversationDetails(convId int, userId uint64) (ConversationDetails, error) {
	log.Printf("Getting details for conversation %d", convId)

	var conv ConversationDetails
	var photoNull sql.NullString // For handling NULL photo

	// Get conversation info
	err := db.c.QueryRow(`
        SELECT 
            c.ConversationId,
            CASE 
                WHEN c.GroupId = 1 THEN c.Name
                WHEN u.Username IS NOT NULL THEN u.Username
                ELSE 'Unknown'
            END as Name,
            c.GroupPhoto as Photo,
            c.GroupId = 1 as IsGroup
        FROM conversations c
        LEFT JOIN participants p ON c.ConversationId = p.ConversationId AND p.UserId != ?
        LEFT JOIN users u ON p.UserId = u.Id
        WHERE c.ConversationId = ?`,
		userId, convId).Scan(
		&conv.ConversationId,
		&conv.Name,
		&photoNull,
		&conv.IsGroup,
	)
	if err != nil {
		log.Printf("Error getting conversation info: %v", err)
		return conv, err
	}

	// Handle NULL photo
	if photoNull.Valid {
		conv.Photo = photoNull.String
	}

	// Get messages with sender info and comments
	rows, err := db.c.Query(`
        SELECT 
            m.MessageId,
            m.Text,
            m.SendTime,
            m.Status,
            m.SenderId,
            m.Photo,
            u.Username as SenderUsername
        FROM messages m
        JOIN users u ON m.SenderId = u.Id
        WHERE m.ConversationId = ?
        ORDER BY m.SendTime DESC`, convId)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return conv, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg MessageWithComments
		var photoNull sql.NullString
		err := rows.Scan(
			&msg.MessageId,
			&msg.Text,
			&msg.SendTime,
			&msg.Status,
			&msg.SenderId,
			&photoNull,
			&msg.SenderUsername,
		)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			return conv, err
		}

		if photoNull.Valid {
			msg.Photo = photoNull.String
		}

		// Get comments for this message
		comments, err := db.getMessageComments(msg.MessageId)
		if err != nil {
			log.Printf("Error getting comments: %v", err)
			return conv, err
		}
		msg.Comments = comments

		conv.Messages = append(conv.Messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error in messages: %v", err)
		return conv, err
	}

	return conv, nil
}

func (db *appdbimpl) getMessageComments(messageId int) ([]Comment, error) {
	rows, err := db.c.Query(`
        SELECT 
            c.UserId,
            u.Username,
            c.Emoji
        FROM comments c
        JOIN users u ON c.UserId = u.Id
        WHERE c.MessageId = ?`, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.UserId, &comment.Username, &comment.Emoji)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error in comments: %v", err)
		return nil, err
	}

	return comments, nil
}

func (db *appdbimpl) SearchUsers(query string) ([]User, error) {
	log.Printf("SearchUsers called with query: %s", query)

	// Use LIKE with wildcards for partial matching
	searchQuery := `
        SELECT Id, Username 
        FROM users 
        WHERE Username LIKE ?
        ORDER BY Username`

	log.Printf("Executing query: %s", searchQuery)
	rows, err := db.c.Query(searchQuery, "%"+query+"%")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		log.Printf("Found user: %+v", user)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	log.Printf("Total users found: %d", len(users))
	return users, nil
}
