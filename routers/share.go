package routers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
)

func getUserIDOrAbort(c *gin.Context) (int64, error) {
	idInterface, exist := c.Get(middlewares.Key.UserID)
	if exist == false {
		c.AbortWithStatusJSON(http.StatusForbidden,
			dto.NewResponseBad("unauthorized action is forbidden"))
		return 0, errors.New("unanthorized")
	}
	id := idInterface.(int64)
	return id, nil
}

func checkConfigTypeRange(t int8) bool {
	return t <= dto.LimitConfigTypeMax && t >= dto.LimitConfigFormatMin
}

func checkConfigFormatRange(r int8) bool {
	return r <= dto.LimitConfigFormatMax && r >= dto.LimitConfigFormatMin
}

func bindOrResponseFailed(c *gin.Context, req interface{}) error {
	err := c.Bind(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad("invalid request parameters"))
	}
	return err
}
