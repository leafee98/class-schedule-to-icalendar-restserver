package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/db"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/dto"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/middlewares"
)

func init() {
	RegisterRouter("/config-create", "post", configCreate)
}

// ConfigCreate will create a new Config
func configCreate(c *gin.Context) {
	var req dto.ConfigCreateReq
	err := c.Bind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewResponseBad(err.Error()))
		return
	}

	ownerIDInterface, exist := c.Get(middlewares.Key.UserID)
	if exist == false {
		c.AbortWithStatusJSON(http.StatusForbidden,
			dto.NewResponseBad("unauthorized action is forbidden"))
		return
	}
	var ownerID int64 = ownerIDInterface.(int64)

	configID, err := db.ConfigCreate(req.Name, req.Content, req.Format, ownerID, req.Remark)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.NewResponseBad(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewResponseFine(dto.ConfigCreateRes{ID: configID}))
}
