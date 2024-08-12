package queue

import (
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   amqp.Queue
}

func NewRabbitMQ(url, queueName string) (Queue, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    q, err := ch.QueueDeclare(
        queueName,
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, err
    }

    return &RabbitMQ{
        conn:    conn,
        channel: ch,
        queue:   q,
    }, nil
}

func (r *RabbitMQ) Send(message []byte) error {
    return r.channel.Publish(
        "",
        r.queue.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )
}

func (r *RabbitMQ) ReceiveCh() (<-chan []byte, error) {
    msgs, err := r.channel.Consume(
        r.queue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }

    ch := make(chan []byte)
    go func() {
        for msg := range msgs {
            ch <- msg.Body
        }
    }()

    return ch, nil
}

func (r *RabbitMQ) Close() error {
    if err := r.channel.Close(); err != nil {
        return err
    }
    return r.conn.Close()
}
