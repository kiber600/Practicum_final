package main

import (
	"fmt"
	"log"

	"Practicum_final/pkg/server" //"github.com/kiber600/Practicum_final/pkg/server"
)

func main() {

	err := server.StartServer()
	if err != nil {
		fmt.Printf("Server don't start")
		log.Println(err)
	}

}
