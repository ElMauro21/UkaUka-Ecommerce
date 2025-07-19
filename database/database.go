package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) *sql.DB{
// Open (or create) the file uka.db
db,err := sql.Open("sqlite3","./uka.db")
if err != nil{
	log.Fatalf("Error al abrir la base de datos: %v", err)
}

_, err = db.Exec("PRAGMA foreign_keys = ON;")
if err != nil {
	log.Fatalf("Error al habilitar claves foráneas: %v", err)
}

runMigration(db)
return db
}

func runMigration(db *sql.DB){

// Run migrations: create tables if they don't already exist
mustExec(db,`
CREATE TABLE IF NOT EXISTS users(
id	INTEGER PRIMARY KEY AUTOINCREMENT,
names 	TEXT NOT NULL,
surnames TEXT,
id_number TEXT,
phone TEXT, 
email TEXT NOT NULL UNIQUE,
state TEXT,
city TEXT,
neighborhood TEXT,
address TEXT,
password_hash TEXT NOT NULL,
is_admin INTEGER NOT NULL DEFAULT 0,
failed_attempts INTEGER NOT NULL DEFAULT 0,
locked_until DATETIME
);`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS products (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL, 
description TEXT NOT NULL,
weight INTEGER NOT NULL,
size INTEGER NOT NULL,
price FLOAT NOT NULL,
quantity INTEGER NOT NULL,
image_url TEXT,
image_url_2 TEXT
);`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS carts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
user_id INTEGER,
session_id TEXT,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(user_id) REFERENCES users(id)
);`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS cart_items(
id INTEGER PRIMARY KEY AUTOINCREMENT,
cart_id INTEGER NOT NULL,
product_id INTEGER NOT NULL,
quantity INTEGER NOT NULL DEFAULT 1,
FOREIGN KEY(cart_id) REFERENCES carts(id) ON DELETE CASCADE,
FOREIGN KEY(product_id) REFERENCES products(id)
);`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS password_resets(
id INTEGER PRIMARY KEY AUTOINCREMENT,
email TEXT NOT NULL,
token TEXT NOT NULL,
expires_at DATETIME NOT NULL
);`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS transactions (
id INTEGER PRIMARY KEY AUTOINCREMENT,
user_id INTEGER,
reference_code TEXT UNIQUE NOT NULL,
total_amount FLOAT NOT NULL,
status TEXT NOT NULL DEFAULT 'pending',
shipped INTEGER NOT NULL DEFAULT 0,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(user_id) REFERENCES users(id)
);
`)
mustExec(db,`
CREATE TABLE IF NOT EXISTS shipping_info (
id INTEGER PRIMARY KEY AUTOINCREMENT,
transaction_id INTEGER NOT NULL,
full_name TEXT NOT NULL,
id_number TEXT NOT NULL,
phone TEXT NOT NULL,
email TEXT NOT NULL,
state TEXT NOT NULL,
city TEXT NOT NULL,
neighborhood TEXT NOT NULL,
address TEXT NOT NULL,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(transaction_id) REFERENCES transactions(id)
);`)
mustExec(db, `
CREATE TABLE IF NOT EXISTS transaction_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	transaction_id INTEGER NOT NULL,
	product_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
	FOREIGN KEY(product_id) REFERENCES products(id)
);
`)
}

// helper: panic if migration fails
func mustExec(db *sql.DB, stmt string){
	if _,err := db.Exec(stmt); err != nil{
		log.Fatalf("Error al crear tablas: %v\n-- statement:\n%s",err,stmt)
	}
}