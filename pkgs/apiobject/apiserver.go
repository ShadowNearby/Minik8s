package core

import "github.com/gin-gonic/gin"

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

// InfoType either Data or Error exists
type InfoType struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}
