package flash

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetMessage(c *gin.Context)(any,any){
	session := sessions.Default(c)
	msg := session.Get("message")
	msgType := session.Get("messageType")

	session.Delete("message")
	session.Delete("messageType")
	session.Save()

	return msg,msgType
}

func SetMessage(c *gin.Context,msg,msgtype string){
	session := sessions.Default(c)
	session.Set("message",msg)
	session.Set("messageType",msgtype)
	session.Save()
}