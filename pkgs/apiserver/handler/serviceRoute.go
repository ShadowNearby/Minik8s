package handler

import "github.com/gin-gonic/gin"

type ServiceRoutes struct {
	ServiceName string
	Routes      []Route
}

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

func (r *Route) register(engine *gin.Engine) {
	switch r.Method {
	case "get":
		engine.GET(r.Path, r.Handler)
	case "post":
		engine.POST(r.Path, r.Handler)
	case "put":
		engine.PUT(r.Path, r.Handler)
	case "delete":
		engine.DELETE(r.Path, r.Handler)
	case "GET":
		engine.GET(r.Path, r.Handler)
	case "POST":
		engine.POST(r.Path, r.Handler)
	case "PUT":
		engine.PUT(r.Path, r.Handler)
	case "DELETE":
		engine.DELETE(r.Path, r.Handler)
	default:
		panic("invalid HTTP method")
	}
}

func (service *ServiceRoutes) registerRoutes(engine *gin.Engine) {
	for i := range service.Routes {
		service.Routes[i].register(engine)
	}
}
