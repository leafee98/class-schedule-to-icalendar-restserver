package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// KeyStruct is the struct used for keys setted in gin.Context
type KeyStruct struct {
	UserID string
}

// Key stored the key values setted in gin.Context
var Key KeyStruct

func init() {
	Key = KeyStruct{
		UserID: "userID",
	}
	userTokenMap = make(map[string]int64)

	registerMiddleware(verifyUser)
}

var userTokenMap map[string]int64

func verifyUser(c *gin.Context) {
	v, err := c.Cookie("token")
	if err == nil {
		userID, exists := userTokenMap[v]
		if exists {
			c.Set(Key.UserID, userID)
			logrus.Infof("userID=%v verified", userID)
			return
		}
	}
	logrus.Infof("unauthorized visitor")
	c.Set(Key.UserID, nil)
}

// RegisterToken will add and random token and userId in userTokenMap, and return the token,
// so the middleware could add userId info to gin.Context per HTTP request later
func RegisterToken(userID int64) string {
	u, _ := uuid.NewRandom()
	userTokenMap[u.String()] = userID
	return u.String()
}

// ExpireToken will remove the token in userTokenMap
func ExpireToken(token string) {
	delete(userTokenMap, token)
}
