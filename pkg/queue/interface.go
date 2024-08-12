package queue

// Both CLient and Server expect queue implementation to match this interface
type Queue interface {
    Send(message []byte) error
    ReceiveCh() (<-chan []byte, error)
    Close() error
}
