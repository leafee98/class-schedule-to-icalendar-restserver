package model

import "time"

// User is stored in database and transformed in AJAX request and response
type User struct {
	ID       int64     `db:"c_id" json:"id"`
	Email    string    `db:"c_email" json:"email"`
	Nickname string    `db:"c_nickname" json:"nickname"`
	Username string    `db:"c_username" json:"username"`
	Password [32]byte  `db:"c_password" json:"password"`
	Bio      string    `db:"c_bio" json:"bio"`
	JoinTime time.Time `db:"c_join_time" json:"joinTime"`
}
