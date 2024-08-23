package handlers

import (
	"c07_practica/pkg/model"
	models "c07_practica/pkg/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UserRouterHandlers(router *mux.Router, db *sql.DB, secret string) {
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/users/register", createUser(db)).Methods("POST")
	router.HandleFunc("/users/login", loginUser(db, secret)).Methods("POST")
	router.HandleFunc("/users/{id}", getUserByID(db)).Methods("GET")
	router.Handle("/users/update", model.AuthMiddleware(updateUser(db), secret)).Methods("PUT")
}

func loginUser(db *sql.DB, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user_id, err := model.Login(db, user.Email, user.Password)

		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		jwt := model.CreateJWT(secret, user_id)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwt)
	}
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := models.GetUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := models.CreateUser(db, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User Created Successfully")
	}
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.UpdatedUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user_id, ok := model.ObtenerIdJWT(r.Context())

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userId := int(user_id)
		user.ID = &userId

		if err := models.UpdateUser(db, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User Updated Successfully")
	}
}

func getUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
