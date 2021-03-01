package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
)

// Engine is gin Engine
var Engine *gin.Engine

// Init create the gin engine
func Init() {
	Engine = gin.Default()
}

// Run start listen and handle http request, call this func at last
func Run() {
	Engine.Run(config.RestEndpoint)
}
