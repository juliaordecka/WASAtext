package database

import (
	"database/sql"
)

func (db *appdbimpl) CreateUser(u User) (User, error) {
	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", u.Username)
	if err != nil {
		var user User
		if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.Username).Scan(&user.Id, &user.Username); err != nil {
			if err == sql.ErrNoRows {
				return user, ErrUserDoesNotExist
			}
		}
		return user, nil
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = uint64(lastInsertID)
	return u, nil
}

func (db *appdbimpl) SetUsername(u User, username string) (User, error) {
	res, err := db.c.Exec(`UPDATE users SET Username=? WHERE Id=? AND Username=?`, u.Username, u.Id, username)
	if err != nil {
		return u, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return u, err
	} else if affected == 0 {
		return u, err
	}
	return u, nil
}

func (db *appdbimpl) GetUsernameById(userId uint64) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT Username FROM users WHERE Id = ?", userId).Scan(&username)
	return username, err
}

// User operations
func (db *appdbimpl) GetRecipientIdByUsername(username string) (uint64, error) {
	var userId uint64
	err := db.c.QueryRow("SELECT Id FROM users WHERE Username = ?", username).Scan(&userId)
	return userId, err
}

func (db *appdbimpl) GetConversationIdByName(name string) (int, error) {
	var convId int
	err := db.c.QueryRow("SELECT ConversationId FROM conversations WHERE Name = ?", name).Scan(&convId)
	return convId, err
}
