package admin

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdminUser(db *sql.DB){

	password := os.Getenv("ADMIN_PASSWORD")
  	if password == ""{
    	log.Fatal("ADMIN_PASSWORD is not set")
  	}

	name := "Admin"
	surname := "Admin"
	idNumber := "00000000"
	phone := "000000000"
	email := "admin@example.com"
	state := "00000000"
	city := "00000000"
	neigb := "00000000"
	addr := "00000000"

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if count > 0{
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	_,err := db.Exec(`INSERT INTO users 
	(names,surnames,id_number,phone, email,state,city,neighborhood,address, password_hash, is_admin) 
	VALUES (?,?,?,?,?,?,?,?,?,?,1)
	`,name,surname,idNumber,phone,email,state,city,neigb,addr, hashedPassword)
	if err != nil {
		log.Printf("Error al crear usuario admin: %v", err)
	}
}