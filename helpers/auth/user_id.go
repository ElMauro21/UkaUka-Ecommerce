package auth

import (
	"database/sql"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context, db *sql.DB) (int,error) {
	session := sessions.Default(c)
	email := session.Get("user")

	if email == nil {
		return 0,errors.New("usuario no está en sesión")
	}

	var userID int 
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
		if err != nil {
			return 0,err
		}
		return userID,nil
}