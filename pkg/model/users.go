package model

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatedUser struct {
	ID       *int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatalf("Error al hashear la contraseña: %v", err)
	}
	result, err := db.Exec(query, user.Name, user.Email, hashedPassword)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

func UpdateUser(db *sql.DB, user *UpdatedUser) error {
	query := "UPDATE users SET name = COALESCE(?, name), email = COALESCE(?, email), password = COALESCE(?, password) WHERE id = ?"
	hashedPassword, bc_error := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if bc_error != nil {
		return bc_error
	}
	fmt.Print(user)
	_, err := db.Exec(query, user.Name, user.Email, hashedPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	var user User
	query := "SELECT id, name, email FROM users WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Login(db *sql.DB, email string, password string) (int, error) {
	var user User
	query := "SELECT id, password FROM users WHERE email = ?"

	// Ejecutar la consulta
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Usuario no encontrado, manejar de forma apropiada
			log.Printf("Usuario no encontrado: %v", err)
		} else {
			// Error en la consulta
			log.Printf("Error al ejecutar la consulta: %v", err)
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		log.Printf("Email y/o Password Incorrectas")
		return 0, err
	}

	return user.ID, nil
}
