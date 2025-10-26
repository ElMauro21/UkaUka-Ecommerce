package handlers

import (
	"net/http"

	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-gonic/gin"
)

func HandleOpenHome(c *gin.Context){
	view.Render(c,http.StatusOK,"index.html",gin.H{
		"title": "UkaUka fb",
	})
}