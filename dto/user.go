package dto

// UserRegisterReq is used when requesting to register a new user
type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// UserRegisterRes is the response of register user request
type UserRegisterRes struct {
	ID int64 `json:"id"`
}

// UserLoginReq is usd when requesting to login
type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginRes is the response of login request
type UserLoginRes struct {
	ID int64 `json:"id"`
}
