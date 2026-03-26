package main

import (
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
}

func (hndl SettingsHandler) Action(c *gin.Context) {

	switch c.Request.Method {
	case "GET":
		fallthrough
	case "":
		conf := Config{
			Volume:     config.Volume,
			WifiAPName: config.WifiAPName,
		}
		c.JSON(200, conf)
	case "POST":
		config := Config{}
		if err := c.ShouldBindBodyWithJSON(&config); err != nil {
			c.Status(400)
		} else {
			SetConfig(config)
		}
	default:
		c.Status(404)
	}
}

func (hndl SettingsHandler) Endpoint() string {
	return "/settings"
}
