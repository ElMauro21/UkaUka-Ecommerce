package cart

import (
	"database/sql"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetCartID(c *gin.Context, db *sql.DB) (int, error) {
	session := sessions.Default(c)
	email := session.Get("user")

	var cartID int
	if email != nil {
		userID,err  := auth.GetUserId(c, db)
		if err != nil {
			return 0,err
		}
		err = db.QueryRow("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		return cartID, err
	} else {
		sessionID := session.Get("cart_session_id")
		err := db.QueryRow("SELECT id FROM carts WHERE session_id = ?", sessionID).Scan(&cartID)
		return cartID, err
	}
}