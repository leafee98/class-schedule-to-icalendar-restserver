package routers

import (
	"crypto/sha256"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
	"github.com/sirupsen/logrus"
)

// contains routers relative to user's account

// run automatically to register routers in this file
func init() {
	RegisterRouter("/register", "post", register)
	RegisterRouter("/login", "post", login)
	RegisterRouter("/logout", "post", logout)
	RegisterRouter("/logout", "get", logout)
}

func register(c *gin.Context) {
	var req dto.UserRegisterReq
	if bindOrResponseFailed(c, &req) != nil {
		return
	}

	var cnt int64
	err := db.DB.Get(&cnt, "select count(c_id) from t_user where c_username = ? or c_email = ?",
		req.Username, req.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
		logrus.Errorf("err while register: %v", err.Error())
		return
	}
	if cnt > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("duplicated username or email"))
		logrus.Warnf("err while register: %s", "duplicated username or email")
		return
	}

	// file the hash password
	var hashPassArr [32]byte = passwordHash(req.PasswordPlain)
	req.Password = hashPassArr[:]

	res, err := db.DB.Exec("insert into t_user (c_username, c_password, c_email, c_nickname) "+
		"values (?, ?, ?, ?)", req.Username, req.Password, req.Email, req.Email)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		logrus.Error(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		logrus.Error(err)
	} else {
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.UserRegisterRes{ID: id}))
	}
}

func login(c *gin.Context) {
	var req dto.UserLoginReq
	if bindOrResponseFailed(c, &req) != nil {
		return
	}

	var dbID int64
	var dbPassword [32]byte
	row := db.DB.QueryRow("select c_id, c_password from t_user where c_username = ?", req.Username)

	var tmp []byte
	switch err := row.Scan(&dbID, &tmp); err {
	case sql.ErrNoRows:
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("user not exists"))
		return
	case nil:
		// fine
	default:
		c.AbortWithStatusJSON(http.StatusBadGateway, dto.NewResponseBad(err.Error()))
		return
	}

	copy(dbPassword[:], tmp)

	var passCipher [32]byte = passwordHash(req.PasswordPlain)
	if passCipher == dbPassword {
		// logdin success, register token and set cookie
		token := middlewares.RegisterToken(dbID, req.TokenDuration)
		c.SetCookie("token", token, 3600*24*req.TokenDuration, "/", "", true, false)
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.UserLoginRes{ID: dbID}))
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("password incorrect"))
	}
}

func logout(c *gin.Context) {
	token, err := c.Cookie("token")
	c.SetCookie("token", "000", -1, "/", "", true, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("unauthorized logout is forbidden"))
	}
	middlewares.ExpireToken(token)
}

// passwordHash do sha256 as hash to avoid using plain text
func passwordHash(p string) [32]byte {
	return sha256.Sum256([]byte(p))
}
