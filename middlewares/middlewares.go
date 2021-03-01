package middlewares

import "github.com/gin-gonic/gin"

var ms []gin.HandlerFunc = make([]gin.HandlerFunc, 0, 4)

// Init is used to initialzed middlewares with gin.Engine
func Init(r *gin.Engine) {
	for _, m := range ms {
		r.Use(m)
	}
}

func registerMiddleware(m gin.HandlerFunc) {
	ms = append(ms, m)
}
