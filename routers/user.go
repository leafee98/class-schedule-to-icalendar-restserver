package routers

import (
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
}

func register(c *gin.Context) {
	var req dto.UserRegisterReq
	c.BindJSON(&req)
	id, err := db.CreateUser(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		logrus.Error(err)
	} else {
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.UserRegisterRes{ID: id}))
	}
}

func login(c *gin.Context) {
	var req dto.UserLoginReq
	c.Bind(&req)

	id, err := db.VerifyUser(req.Username, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
	} else {
		token := middlewares.RegisterToken(id)
		c.SetCookie("token", token, 3600*24*7, "/", "", true, false)
		c.JSON(http.StatusOK, dto.NewResponseFine(dto.UserLoginRes{ID: id}))
	}
}
