package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/leafee98/class-schedule-to-icalendar-restserver/config"
)

// Router type alias function to handle some request
type Router struct {
	path   string
	method string
	f      func(*gin.Context)
}

var routers = make([]Router, 0)

// RegisterRouter function store the router to register to Gin temporarily
func RegisterRouter(path string, method string, f func(*gin.Context)) {
	routers = append(routers, Router{path: path, method: method, f: f})
}

// Init register all stored router to Gin
// This request a *gin.Engine initialized in server package
func Init(engine *gin.Engine) error {
	routerGroup := engine.Group(config.HTTPBasepath)
	for _, r := range routers {
		switch r.method {
		case "get":
			routerGroup.GET(r.path, r.f)
		case "post":
			routerGroup.POST(r.path, r.f)
		}
	}
	return nil
}
