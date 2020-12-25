package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	log.SetPrefix("hello-gin: ")
	gin.DisableConsoleColor()

	log.Println("here we go again...")


	server := NewServer()
	server.Routes()
	server.Run()
	server.WaitSignal()

	log.Printf("Server exit")
	os.Exit(0)
}
