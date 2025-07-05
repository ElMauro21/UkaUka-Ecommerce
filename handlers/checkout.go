package handlers

import (
	"net/http"

	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-gonic/gin"
)

func HandleOpenCheckout(c *gin.Context){
	
	msg,msgType := flash.GetMessage(c)

	view.Render(c,http.StatusOK,"checkout.html",gin.H{
		"title": "Checkout",
		"Message": msg,
		"MessageType": msgType,
	})
}