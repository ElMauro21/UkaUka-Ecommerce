package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/users"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-gonic/gin"
)

func HandleOpenCheckout(c *gin.Context, db *sql.DB){
	
	msg,msgType := flash.GetMessage(c)

	u := users.LoadUserInfo(c,db)

	total := c.PostForm("total")
	description := c.PostForm("description")

	view.Render(c,http.StatusOK,"checkout.html",gin.H{
		"title": "Checkout",
		"User":  u,
		"Message": msg,
		"MessageType": msgType,
		"Total":       total,
		"Description": description,
	})
}