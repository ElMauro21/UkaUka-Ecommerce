package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElMauro21/UkaUkafb/helpers/auth"
	"github.com/ElMauro21/UkaUkafb/helpers/flash"
	"github.com/ElMauro21/UkaUkafb/helpers/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandleOpenLogin(c *gin.Context){
	
	msg,msgType := flash.GetMessage(c)

	view.Render(c,http.StatusOK,"login.html",gin.H{
		"title": "Login",
		"Message": msg,
		"MessageType": msgType,
	})
}

func HandleLogin(c *gin.Context, db *sql.DB){
	email := c.PostForm("log-email")
	password := c.PostForm("log-password")

	var storedHash string
	var isAdmin int
	var failedAttempts int
	var lockedUntil sql.NullTime

	err := db.QueryRow(`SELECT password_hash, is_admin, failed_attempts, locked_until 
	FROM users 
	WHERE email = ?`, email).Scan(&storedHash, &isAdmin, &failedAttempts, &lockedUntil)
	if err != nil {
		view.RenderFlash(c,http.StatusOK,"Contraseña o usuario incorrectos","error")
		return
	}

	if lockedUntil.Valid && time.Now().Before(lockedUntil.Time){
		view.RenderFlash(c, http.StatusOK, "Cuenta bloqueada. Inténtalo más tarde.", "error")
		return
	}

	err = auth.ComparePasswords(storedHash,password)
	if err != nil {
		failedAttempts ++
		if failedAttempts >= 5 {
			lockDuration := 15 * time.Minute
			_, _ = db.Exec("UPDATE users SET failed_attempts = ?, locked_until = ? WHERE email = ?",
				failedAttempts, time.Now().Add(lockDuration), email)
			view.RenderFlash(c, http.StatusOK, "Demasiados intentos fallidos. Cuenta bloqueada por 15 minutos.", "error")
		} else {
			_, _ = db.Exec("UPDATE users SET failed_attempts = ? WHERE email = ?", failedAttempts, email)
			view.RenderFlash(c, http.StatusOK, "Contraseña o usuario incorrectos", "error")
		}
		return
	}

	_, _ = db.Exec("UPDATE users SET failed_attempts = 0, locked_until = NULL WHERE email = ?", email)

	session := sessions.Default(c)
	session.Set("user", email)
	session.Set("isAdmin", isAdmin == 1)
	session.Set("loggedIn", true)
	session.Save()
	c.Header("HX-Redirect","/")
	c.Status(http.StatusOK)
}

func HandleLogout(c *gin.Context){
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/auth/login")
}

func HandleRegister(c *gin.Context, db *sql.DB){
	name := c.PostForm("name")
	surname := c.PostForm("surname")
	idNumber := c.PostForm("id-number")
	phone := c.PostForm("phone")
	mail := c.PostForm("reg-mail")
	state := c.PostForm("state")
	city := c.PostForm("city")
	neigb := c.PostForm("neighborhood")
	addr := c.PostForm("address")
	pass1 := c.PostForm("reg-password1")
	pass2 := c.PostForm("reg-password2")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", mail).Scan(&count)
	if count > 0{
		view.RenderFlash(c,http.StatusOK,"No se ha podido crear usuario","error")
		return
	}

	if pass1 != pass2{
		view.RenderFlash(c,http.StatusOK,"Las contraseñas deben coincidir","error")
		return
	}

	hashedPassword, err := auth.HashPassword(pass1)
	if err != nil{
		c.String(http.StatusInternalServerError, "No se ha podido crear usuario.")
		return
	}

	_, err = db.Exec(`INSERT INTO users 
	(names, surnames, id_number, phone, email, state, city, neighborhood, address, password_hash) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, 
	name, surname, idNumber, phone, mail, state, city, neigb, addr, hashedPassword)
	if err != nil{
		c.String(http.StatusInternalServerError, "No se ha podido crear usuario.")
		return
	}
	
	flash.SetMessage(c,"Usuario creado correctamente","success")
	c.Header("HX-Redirect","/auth/login")
	c.Status(http.StatusSeeOther)
}

func HandleCreateRecoveryLink(c *gin.Context, db *sql.DB){
	email:= c.PostForm("recover-email")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil{
		c.String(http.StatusInternalServerError, "Error en la base de datos.")
		return 
	}
	
	if count == 0{
		flash.SetMessage(c,"Si este correo existe, ya ha sido enviado un correo con las instrucciones","success")
		c.Redirect(http.StatusSeeOther, "/auth/login")
		return 
	}

	token, err := auth.GenerateRandomToken(32)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo crear token de reseteo.")
		return
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	_, err = db.Exec(`INSERT INTO password_resets 
	(email, token, expires_at) VALUES (?, ?, ?)`,
	 email,token,expiresAt)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo guardar token de reseteo")
		return
	}

	resetLink := fmt.Sprintf("http://localhost:8080/auth/recover/reset?token=%s", token)
	
	c.Set("reset_link",resetLink)
	c.Set("reset_email",email)

	go func(cCopy *gin.Context){
		if err := auth.SendRecoveryEmail(cCopy); err != nil {
			log.Println("No se pudo enviar correo de recuperación:", err)
		}
	}(c.Copy())

	flash.SetMessage(c,"Si este correo existe, ya ha sido enviado un correo con las instrucciones","success")
	c.Redirect(http.StatusSeeOther, "/auth/login")
}

func HandleShowResetForm(c *gin.Context){
	token := c.Query("token")

	c.HTML(http.StatusOK,"reset-password.html",gin.H{
		"token": token,
		"title": "UkaUka fb",
	})
}

func HandleResetPassword(c *gin.Context, db *sql.DB){
	token := c.PostForm("token")
	newPassword := c.PostForm("recover-pass")

	if token == "" || newPassword == ""{
		view.RenderFlash(c,http.StatusOK,"Se requiere token y y contraseña","error")
		return
	}

	var email string
	var expiresAt time.Time
	err := db.QueryRow("SELECT email, expires_at FROM password_resets WHERE token = ?", token).Scan(&email, &expiresAt)
	if err != nil {
		view.RenderFlash(c,http.StatusOK,"El token ha expirado!","error")
		return
	}

	if time.Now().After(expiresAt) {
		view.RenderFlash(c,http.StatusOK,"El token ha expirado!","error")
		return
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo actualizar contraseña.")
		return
	}

	_, err = db.Exec("UPDATE users SET password_hash = ? WHERE email = ?", hash, email)
	if err != nil {
		c.String(http.StatusInternalServerError, "No se pudo actualizar contraseña..")
		return
	}

	_, _ = db.Exec("DELETE FROM password_resets WHERE token = ?", token)

	flash.SetMessage(c,"La contraseña ha sido actualizada correctamente","success")
	c.Redirect(http.StatusSeeOther, "/auth/login")
}
