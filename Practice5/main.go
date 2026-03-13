package main

import (
	"Practice5/db"
	"Practice5/handler"
	"Practice5/repository"
	"log"
	"net/http"
)

func main() {
	database := db.NewDB()
	defer database.Close()

	db.Migrate(database)

	userRepo := repository.NewUserRepository(database)

	h := handler.NewHandler(userRepo, database)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
