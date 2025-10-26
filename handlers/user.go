package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/users"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandleOpenProfile(c *gin.Context, db *sql.DB){
	
	u:= users.LoadUserInfo(c,db)
	
	msg,msgType := flash.GetMessage(c)
    view.Render(c, http.StatusOK, "profile.html", gin.H{
        "title": "Perfil",
        "User":  u,
		"Message": msg,
		"MessageType": msgType,
    })
}

func HandleUpdateProfile(c *gin.Context, db *sql.DB) {
	session := sessions.Default(c)
	email := session.Get("user")
	if email == nil {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}

	name := c.PostForm("name")
	surname := c.PostForm("surname")
	idNumber := c.PostForm("id-number")
	phone := c.PostForm("phone")
	state := c.PostForm("state")
	city := c.PostForm("city")
	neighborhood := c.PostForm("neighborhood")
	address := c.PostForm("address")

	var current struct {
		Names       string
		Surnames    string
		IDNumber    string
		Phone       string
		State       string
		City        string
		Neighborhood string
		Address     string
	}

	err := db.QueryRow(`
        SELECT names, surnames, id_number, phone, state, city, neighborhood, address
        FROM users WHERE email = ?`, email).Scan(
		&current.Names, &current.Surnames, &current.IDNumber, &current.Phone,
		&current.State, &current.City, &current.Neighborhood, &current.Address)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error no se puede verificar el perfil.")
		return
	}

	if 	name == current.Names &&
		surname == current.Surnames &&
		idNumber == current.IDNumber &&
		phone == current.Phone &&
		state == current.State &&
		city == current.City &&
		neighborhood == current.Neighborhood &&
		address == current.Address {
		view.RenderFlash(c,http.StatusOK,"No se han efectuado cambios.","info")
		return
	}

	_, err = db.Exec(`
        UPDATE users
        SET names = ?, surnames = ?, id_number = ?, phone = ?, state = ?, city = ?, neighborhood = ?, address = ?
        WHERE email = ?
    `, name, surname, idNumber, phone, state, city, neighborhood, address, email)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error al actualizar perfil.")
		return
	}

	view.RenderFlash(c,http.StatusOK,"Perfil actualizado correctamente.","success")
}

func HandleChangePassword(c *gin.Context, db *sql.DB) {
	session := sessions.Default(c)
	email := session.Get("user")
	if email == nil {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}

	currentPass := c.PostForm("current-password")
	newPass1 := c.PostForm("reg-password1")
	newPass2 := c.PostForm("reg-password2")

	if newPass1 != newPass2 {
		view.RenderFlash(c,http.StatusOK,"Las nuevas contraseñas deben coincidir","error")
		return
	}

	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = ?", email).Scan(&storedHash)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo verificar tu contraseña actual.")
		return
	}

	if err := auth.ComparePasswords(storedHash, currentPass); err != nil {
		view.RenderFlash(c,http.StatusOK,"La contraseña actual es incorrecta","error")
		return
	}

	newHash, err := auth.HashPassword(newPass1)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo actualizar la contraseña.")
		return
	}

	_, err = db.Exec("UPDATE users SET password_hash = ? WHERE email = ?", newHash, email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al guardar la nueva contraseña.")
		return
	}
	
	flash.SetMessage(c,"Contraseña actualizada correctamente","success")
	c.Header("HX-Redirect","/user/profile")
	c.Status(http.StatusSeeOther)
}

func HandleDeleteAccount(c *gin.Context,db *sql.DB){
	session := sessions.Default(c)
	email := session.Get("user")
	if email == nil {
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return
	}

	currentPass := c.PostForm("current-password")
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = ?", email).Scan(&storedHash)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo verificar tu contraseña actual.")
		return
	}

	if err := auth.ComparePasswords(storedHash, currentPass); err != nil {
		view.RenderFlash(c,http.StatusOK,"La contraseña actual es incorrecta","error")
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error al eliminar usuario.")
		return
	}

	session.Clear()
	session.Save()
	flash.SetMessage(c,"Cuenta eliminada correctamente","success")
	c.Header("HX-Redirect","/auth/login")
	c.Status(http.StatusSeeOther)
}