package handler

import "github.com/gin-gonic/gin"

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

func (r *Route) Register(engine *gin.Engine) {
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
