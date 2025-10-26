package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InjectTemplateData() gin.HandlerFunc{
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := session.Get("user")
		loggedIn := session.Get("loggedIn")
		isAdmin := session.Get("isAdmin")

		templateData := gin.H{
			"user": user,
			"loggedIn": loggedIn,
			"isAdmin": isAdmin,
		}

		c.Set("templateData",templateData)
		c.Next()

	}
}