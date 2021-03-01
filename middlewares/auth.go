package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var userTokenMap map[string]int64 = make(map[string]int64)

func verifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, err := c.Cookie("token")
		if err == nil {
			c.Set("userId", userTokenMap[v])
		} else {
			c.Set("userId", nil)
		}
	}
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
