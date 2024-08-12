package server

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/4rgon4ut/lightblocks-assignment/pkg/ordmap"
	"github.com/4rgon4ut/lightblocks-assignment/pkg/queue"
	"github.com/4rgon4ut/lightblocks-assignment/pkg/types"
)

type Server struct {
	queue    queue.Queue
	omp      *ordmap.OrderedMap
	resultCh chan string
	file     *os.File
	done     chan struct{}
}

func New(queueURL, queueName, outputFile string) (*Server, error) {
	q, err := queue.NewRabbitMQ(queueURL, queueName)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Server{
		queue:    q,
		omp:      ordmap.New(),

		resultCh: make(chan string, 10000), // Buffered channel to prevent blocking
											// TODO: make configurable
		file:     f,
		done:     make(chan struct{}),
	}, nil
}

func (s *Server) Start() error {
	receiveCh, err := s.queue.ReceiveCh()
	if err != nil {
		return fmt.Errorf("failed to create receive channel: %w", err)
	}
	go s.processQueue(receiveCh)
	// separate goroutine to write results
	go s.writeResults()
	return nil
}

func (s *Server) processQueue(receiveCh <-chan []byte) {
	for {
		select {
		case <-s.done:
			return
		case msg, ok := <-receiveCh:
			if !ok {
				return
			}
			var cmd types.Command
			if err := json.Unmarshal(msg, &cmd); err != nil {
				log.Errorf("Failed to unmarshal command: %v", err)
				continue
			}
			go s.Execute(&cmd)
		}
	}
}

func (s *Server) Execute(cmd *types.Command) {
	if !cmd.IsValid() {
		log.Errorf("invalid command: %v", cmd)
		return
	}

	switch cmd.Type {
	case types.GetItem:
		value, ok := s.omp.Get(cmd.Key)
		if ok {
			s.resultCh <- fmt.Sprintf("%s: %s", cmd.Key, value)
		}
	case types.GetAllItems:
		values := s.omp.GetAll()
		for k, v := range values {
			s.resultCh <- fmt.Sprintf("%s: %s", k, v)
		}
	case types.AddItem:
		s.omp.Set(cmd.Key, cmd.Value)
	case types.DeleteItem:
		s.omp.Delete(cmd.Key)
	}
}

func (s *Server) writeResults() {
	for {
		select {
		case result := <-s.resultCh:
			_, err := fmt.Fprintln(s.file, result)
			if err != nil {
				log.Errorf("Error writing to file: %v", err)
			}
			// TODO: prob we can accumulate to the file buffer and sync every N time
			// in parallel to not wait for disc i/o every time
			if err := s.file.Sync(); err != nil {
				log.Errorf("Error syncing file: %v", err)
			}
		case <-s.done:
			return
		}
	}
}

func (s *Server) Close() {
	close(s.done)
	s.queue.Close()
	close(s.resultCh)
	s.file.Close()
}