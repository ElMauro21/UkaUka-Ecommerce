package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/products"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-gonic/gin"
)

func HandleOpenShop(c *gin.Context, db *sql.DB){

	msg,msgType := flash.GetMessage(c)

	products := products.LoadProducts(db)

	view.Render(c,http.StatusOK,"shop.html",gin.H{
		"title": "Tienda",
		"Message": msg,
		"MessageType": msgType,
		"products": products,
	})
}
