package main

import (
	"bufio"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/4rgon4ut/lightblocks-assignment/pkg/ioutils"
	"github.com/4rgon4ut/lightblocks-assignment/pkg/queue"
)

var (
    queueURL  string
    queueName string
    inputFile string
)

func init() {
    flag.StringVar(&queueURL, "queue-url", "amqp://guest:guest@localhost:5672/", "URL of the RabbitMQ server")
    flag.StringVar(&queueName, "queue-name", "commands", "Name of the queue to use")
    flag.StringVar(&inputFile, "input-file", "", "File to read commands from")
    flag.Parse()

    if inputFile == "" {
        log.Fatalf("Input file is required")
    }
}

func main() {
    q, err := queue.NewRabbitMQ(queueURL, queueName)
    if err != nil {
        log.Fatalf("Failed to connect to queue: %v", err)
    }
    defer q.Close()


    file, err := os.Open(inputFile)
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)


    cmdCh, err := ioutils.ReadCommandsCh(scanner)

    // Create a channel to listen for OS signals
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    for {
        select {
            case cmd, ok := <-cmdCh:
                if !ok {
                    // Channel is closed, all commands have been processed
                    return
                }
                if err := q.Send(cmd); err != nil {
                    log.Errorf("Failed to send command: %v", err)
                }
            case <-sigCh:
                log.Println("Received shutdown signal. Exiting...")
                return
        }
    }
}

