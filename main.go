package main

import (
	"log"
	"os"
	"os/signal"
	"server/internal/app"
	"syscall"
)

func Run() {
	app := app.New()
	app.MustRun()
	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := app.Stop(); err != nil {
		log.Println("Server stopped gracefully")
	}

}

func main() {
	Run()
}
