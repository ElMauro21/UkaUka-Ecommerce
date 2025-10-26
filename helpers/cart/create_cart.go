package cart

import (
	"database/sql"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCart(c *gin.Context, db *sql.DB) error{

	session := sessions.Default(c)
	email := session.Get("user")
	
	if email != nil{
		userID,err := auth.GetUserId(c,db)
		if err != nil {
			return err
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM carts WHERE user_id = ?", userID).Scan(&count)
		if err != nil {
			return err
		}
		
		if count > 0{
			return nil
		}

		_,err = db.Exec(`INSERT INTO carts
		(user_id) VALUES (?)`,
		userID,
		)
		if err != nil{
			return err
		}
	} else {
		sessionID := session.Get("cart_session_id")
		if sessionID == nil {
			sessionID = uuid.New().String()
			session.Set("cart_session_id",sessionID)
			if err := session.Save(); err != nil {
				return err
			}
		}

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM carts WHERE session_id = ?", sessionID).Scan(&count)
		if err != nil {
			return err
		}
		
		if count > 0{
			return nil
		}

		_,err = db.Exec(`INSERT INTO carts
		(session_id) VALUES (?)`,
		sessionID,
		)
		if err != nil{
			return err 
		}
	}
	return nil
}


