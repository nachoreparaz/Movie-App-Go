package main

import (
	"c07_practica/internal/config"
	"c07_practica/internal/database"
	"c07_practica/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//	Lectura del archivo config.yaml - configuraicones
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//	Conectarme a la DB
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	if err := database.CreateTabla(db); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	//	instancio el router de gorilla mux
	router := mux.NewRouter()
	handlers.UserRouterHandlers(router, db, cfg.SecretJWT)
	handlers.MovieRouterHandlers(router, cfg.MoviesbaseURL, cfg.API_KEY, cfg.SecretJWT, db)

	//	Levantamos el servidor
	log.Printf("Server starting on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
