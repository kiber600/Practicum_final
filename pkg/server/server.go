package server

import (
	"Practicum_final/pkg/api"
	"Practicum_final/pkg/db"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const webDir = "./web"

func StartServer() error {
	err := db.CheckDbFile()

	if err != nil {
		panic(err)
	}
	api.Init()
	defer db.CloseDB()

	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	port := fmt.Sprintf("%s", os.Getenv("TODO_PORT"))
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server started on port: %s", port)
	err = http.ListenAndServe(":"+port, nil)
	return err

}
