package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type PlayFileHandler struct {
}

func (hndl PlayFileHandler) Action(c *gin.Context) {

	switch c.Request.Method {
	case "POST":
		id := c.Param("id")
		SetState(NewPlayState(fmt.Sprintf("%s/%s", PlayableFiles, id)))
	default:
		c.Status(404)
	}
}

func (hndl PlayFileHandler) Endpoint() string {
	return "/files/:id/play"
}
