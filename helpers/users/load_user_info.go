package users

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type user struct {
	Names        string
	Surnames     string
	IDNumber     string
	Phone        string
	Mail         string
	State        string
	City         string
	Neighborhood string
	Address      string
}

func LoadUserInfo(c *gin.Context, db *sql.DB) user {
	session := sessions.Default(c)
	email := session.Get("user")
	if email == nil {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return user{}
	}
	
	emailStr, ok := email.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Invalid session email.")
		return user{}
	}

	var u user
	u.Mail = emailStr

	err := db.QueryRow(`
		SELECT names, surnames, id_number, phone, state, city, neighborhood, address 
		FROM users WHERE email = ?`, emailStr).
		Scan(&u.Names, &u.Surnames, &u.IDNumber, &u.Phone,
			&u.State, &u.City, &u.Neighborhood, &u.Address)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error al cargar perfil.")
		return user{}
	}

	return u
}