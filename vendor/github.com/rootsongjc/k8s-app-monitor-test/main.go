package main

import (
	"os"

	service "github.com/rootsongjc/k8s-app-monitor-test/service"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	server := service.NewServer()
	server.Run(":" + port)
}
