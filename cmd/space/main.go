package main

import (
	"log"
	/* "space/internal/api" */
	"space/internal/app"
)

func main() {
	log.Println("Application start!")

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	application.StartServer()

	log.Println("Application terminated")
}
