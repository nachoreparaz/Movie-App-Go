package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Importa el driver de MySQL
)

func Connect(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func CreateTabla(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users(
		    id INT AUTO_INCREMENT, 
		    name VARCHAR(100) NOT NULL,
		    email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
		    PRIMARY KEY(id)
		)
       `

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Table users created or already exists")
	return nil
}
