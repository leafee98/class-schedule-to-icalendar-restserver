package db

import (
	"crypto/sha256"
	"database/sql"
	"errors"
)

// CreateUser create a new user entry in database with basic info
func CreateUser(username string, password string, email string, nickname string) (int64, error) {
	rows, err := DB.Query("select c_id from t_user where c_username = ? or c_email = ?", username, email)
	if err != nil {
		return 0, err
	} else if rows.Next() {
		return 0, errors.New("duplicated username or email")
	}

	var hashPassArr [32]byte = passwordHash(password)
	res, err := DB.Exec("insert into t_user (c_username, c_password, c_email, c_nickname) values (?, ?, ?, ?)",
		username, hashPassArr[:], email, nickname)
	return res.LastInsertId()
}

// VerifyUser will verify user's username and password,
// return user's ID if success, other will return a error
func VerifyUser(username, password string) (int64, error) {
	var dbID int64
	var dbPassword [32]byte
	row := DB.QueryRow("select c_id, c_password from t_user where c_username = ?", username)

	var tmp []byte
	switch err := row.Scan(&dbID, &tmp); err {
	case sql.ErrNoRows:
		return 0, errors.New("user not exists")
	case nil:
		// fine
	default:
		return 0, err
	}

	copy(dbPassword[:], tmp)

	var passCipher [32]byte = passwordHash(password)
	if passCipher == dbPassword {
		return dbID, nil
	}
	return 0, errors.New("password incorrect")
}

// ChangeEmail will change the user's email,
// will NOT check if the email address is valid,
// error if length does not meet the limit
func ChangeEmail(userID int, newEmail string) error {
	return nil
}

// ChangeNickname will change user's nickname,
// error if userID not exist or the newNickname is longer than limit.
func ChangeNickname(userID int, newNickname string) error {
	return nil
}

// ModifyBio modify user's bio, userId required
// error if userID not exist or the newBio is longer than limit.
func ModifyBio(userID int, newBio string) error {
	return nil
}

// ChangePassword change the password of a user,
// the newPassword should be plaintext and will be encrypted by SHA256 in this function
func ChangePassword(UserID int, newPassword string) error {
	return nil
}

// passwordHash do sha256 as hash to avoid using plain text
func passwordHash(p string) [32]byte {
	return sha256.Sum256([]byte(p))
}
