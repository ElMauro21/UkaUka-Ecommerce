package view

import (
	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, code int, name string, data gin.H){
	templateData, _ := c.Get("templateData")
	baseData, _ := templateData.(gin.H)

	for k, v := range data {
		baseData[k] = v
	}

	c.HTML(code,name,baseData)
}

func RenderFlash(c *gin.Context, status int, message, messageType string){
	flash.SetMessage(c, message, messageType)
	msg, msgType := flash.GetMessage(c)

	Render(c, status, "flash.html", gin.H{
		"Message": msg,
		"MessageType": msgType,
	})
}