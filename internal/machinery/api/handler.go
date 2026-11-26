package api

import "github.com/gin-gonic/gin"

// Handler A handler represents a group of endpoints that should be registered with an API
type Handler interface {
	Prefix() string
	Register(*gin.Engine)
}
