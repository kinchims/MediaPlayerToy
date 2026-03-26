package api

import "github.com/gin-gonic/gin"

type Handler interface {
	Endpoint() string
	Action(c *gin.Context)
}
