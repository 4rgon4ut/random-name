package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/4rgon4ut/lightblocks-assignment/pkg/server"
)



func main() {
    // Define command-line flags
    queueURL := flag.String("queue-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
    queueName := flag.String("queue-name", "commands", "Queue name")
    outputFile := flag.String("output", "output.txt", "Output file path")
    flag.Parse()

    // Create and start the server
    srv, err := server.New(*queueURL, *queueName, *outputFile)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }

    log.Println("Server started. Press Ctrl+C to stop.")

    // Set up signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Wait for termination signal
    <-sigChan

    log.Println("Shutting down server...")
    srv.Close()
    log.Println("Server stopped.")
}