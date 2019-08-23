package pubsub

import (
	"context"
	"encoding/json"
	"time"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const bashTopic = "logs_topic"
const bashResult = "logs_result"

type BashResults struct {
	// connection pool
	pool *redis.Pool
}

func NewBashResults(pool *redis.Pool) *BashResults {
	return &BashResults{pool: pool}
}

func (b *BashResults) Write(bs []json.RawMessage) error {
	conn := b.pool.Get()
	defer conn.Close()

	for _, v := range bs {
	    n, err := redis.Int(conn.Do("PUBLISH", bashTopic, string(v)))
	    if err != nil {
	    	return errors.Wrap(err, "PUBLISH failed to channel "+bashTopic)
	    }
	    if n == 0 {
	    	return noSubscriberError{bashTopic}
	    }
	}

	return nil
}

func (b *BashResults) ReadChannel(ctx context.Context) (<-chan interface{}, error) {
	outChannel := make(chan interface{})

	conn := redis.PubSubConn{Conn: b.pool.Get()}
	conn.Subscribe(bashResult)
	msgChannel := make(chan interface{})
	// Run a separate goroutine feeding redis messages into
	// msgChannel
	go receiveMessages(&conn, msgChannel)

	go func() {
		defer close(outChannel)
		defer conn.Close()

		for {
			// Loop reading messages from conn.Receive() (via
			// msgChannel) until the context is cancelled.
			select {
			case msg, ok := <-msgChannel:
				if !ok {
					conn.Close()
					time.Sleep(10 * time.Second)
					msgChannel = make(chan interface{})
					conn = redis.PubSubConn{Conn: b.pool.Get()}
					conn.Subscribe(bashResult)
					go receiveMessages(&conn, msgChannel)
				}
				switch msg := msg.(type) {
				case string:
					outChannel <- msg
				case error:
					outChannel <- errors.Wrap(msg, "reading from redis")
				}

			case <-ctx.Done():
				conn.Unsubscribe()

			}
		}

	}()
	return outChannel, nil
}