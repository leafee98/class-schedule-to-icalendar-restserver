package middlewares

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/sirupsen/logrus"
)

// KeyStruct is the words used as keys setted in gin.Context
type KeyStruct struct {
	UserID string
}

// Key stored the key values setted in gin.Context
var Key KeyStruct

func init() {
	Key = KeyStruct{
		UserID: "userID",
	}

	registerMiddleware(verifyUser)
}

// add Key.UserID's valuein gin.Context base on request's token
func verifyUser(c *gin.Context) {
	v, err := c.Cookie("token")
	if err == nil {
		var userID int64
		row := db.DB.QueryRow("select c_user_id from t_login_token where c_token = ?", v)
		err := row.Scan(&userID)

		if err == nil {
			c.Set(Key.UserID, userID)
			logrus.Infof("userID=%v verified", userID)
			return
		} else if err != sql.ErrNoRows {
			logrus.Error(err.Error())
		}
	}
	logrus.Info("unauthorized visitor")
}

// RegisterToken will add an random token and userId in userTokenMap, and return the token,
// so the middleware could add userId info from token to gin.Context per HTTP request later.
// the token will be uuid string compatiable with rfc4122 removed dashes.
//
// duartion in hours
//
// Register new token will remove the older token.
func RegisterToken(userID int64, duration int) string {
	u, _ := uuid.NewRandom()
	var uuid string = u.String()

	// remove dashes in token
	var builder strings.Builder
	builder.Grow(len(uuid))
	for _, c := range uuid {
		if c != '-' {
			builder.WriteRune(c)
		}
	}
	token := builder.String()

	// remove the older token
	_, err := db.DB.Exec("delete from t_login_token where c_user_id = ?", userID)
	if err != nil {
		logrus.Error(err.Error())
	}

	// insert new token into database
	_, err = db.DB.Exec("insert into t_login_token (c_user_id, c_token, c_expire_time) values"+
		" (?, ?, now() + interval ? hour)",
		userID, token, duration)
	if err != nil {
		logrus.Error(err.Error())
	}

	return token
}

// ExpireToken will remove the token in userTokenMap
func ExpireToken(token string) {
	db.DB.Exec("deleted from t_login_token where c_token = ?", token)
}
