package dto

import "crypto/sha256"

// Full Database Properties
// ID       int64     `db:"c_id" json:"id" binding:"required"`
// Email    string    `db:"c_email" json:"email" binding:"required"`
// Nickname string    `db:"c_nickname" json:"nickname" binding:"required"`
// Username string    `db:"c_username" json:"username" binding:"required"`
// Password []byte    `db:"c_password"`
// Bio      string    `db:"c_bio" json:"bio" binding:"required"`
// JoinTime time.Time `db:"c_join_time" json:"joinTime" binding:"required"`

// Properties only in request
// PasswordPlain string `json:"password" binding:"required"`

// UserRegisterReq is used when requesting to register a new user
type UserRegisterReq struct {
	Username      string `db:"c_username" json:"username" binding:"required"`
	Password      []byte `db:"c_password"`
	PasswordPlain string `json:"password" binding:"required"`
	Nickname      string `db:"c_nickname" json:"nickname" binding:"required"`
	Email         string `db:"c_email" json:"email" binding:"required"`
}

// UserRegisterRes is the response of register user request
type UserRegisterRes struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

// UserLoginReq is usd when requesting to login
type UserLoginReq struct {
	Username      string `db:"c_username" json:"username" binding:"required"`
	Password      []byte `db:"c_password"`
	PasswordPlain string `json:"password" binding:"required"`
	TokenDuration int    `json:"tokenDuration" binding:"required"`
}

// UserLoginRes is the response of login request
type UserLoginRes struct {
	ID int64 `db:"c_id" json:"id" binding:"required"`
}

func (u *UserRegisterReq) fillPasswordHash() {
	if u.Password != nil {
		return
	}
	var hash [32]byte = sha256.Sum256([]byte(u.PasswordPlain))
	var hashSlice []byte = make([]byte, 0, 32)
	copy(hashSlice, hash[:])
	u.Password = hashSlice
}

func (u *UserLoginReq) fillPasswordHash() {
	if u.Password != nil {
		return
	}
	var hash [32]byte = sha256.Sum256([]byte(u.PasswordPlain))
	var hashSlice []byte = make([]byte, 0, 32)
	copy(hashSlice, hash[:])
	u.Password = hashSlice
}
