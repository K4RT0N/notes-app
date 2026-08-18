package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"notesapp/api"
	"notesapp/dao"
	"notesapp/service"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	connMsg := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", connMsg)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tm := dao.NewTransactionManager(db)
	authDataRepository := dao.NewAuthDataRepository(db)
	sessionRepository := dao.NewSessionRepository(db)
	userRepository := dao.NewUserRepository(db)
	noteRepository := dao.NewNoteRepository(db)

	sessionService := service.NewSessionService(sessionRepository, authDataRepository)
	noteService := service.NewNoteService(noteRepository, sessionRepository)
	userService := service.NewUserService(userRepository, sessionRepository)
	registrationService := service.NewRegistrationService(tm, authDataRepository)

	noteHandler := api.NewNoteHandler(noteService)
	userHandler := api.NewUserHandler(userService)
	authHandler := api.NewAuthHandler(sessionService)
	registrationHandler := api.NewRegistrationHandler(registrationService)

	router := mux.NewRouter()

	router.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Get).Methods("GET")
	router.HandleFunc("/api/users/{userId:[0-9]+}/notes", noteHandler.GetUserNotes).Methods("GET")
	router.HandleFunc("/api/users/", userHandler.GetAll).Methods("GET")
	router.HandleFunc("/api/users/me", userHandler.GetMe).Methods("GET")
	router.HandleFunc("/api/users/me/notes", noteHandler.GetMyNotes).Methods("GET")

	router.HandleFunc("/api/notes/", noteHandler.CreateNote).Methods("POST")
	router.HandleFunc("/api/notes/{noteId:[0-9]+}", noteHandler.DeleteById).Methods("DELETE")
	router.HandleFunc("/api/notes/{noteId:[0-9]+}", noteHandler.GetById).Methods("GET")

	router.HandleFunc("/api/auth/registration", registrationHandler.Registrate).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Authorize).Methods("POST")
	router.HandleFunc("/api/auth/logout", authHandler.Logout).Methods("POST")

	log.Println("Server is listening...")
	http.ListenAndServe(":8080", router)
}
