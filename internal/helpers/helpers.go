package helpers

import (
	"bytes"
	"context"
	"time"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"

	"github.com/grafana/xk6-http/internal/events"
)

// Must panics and throw error on rt if err is not nil
func Must(rt *sobek.Runtime, err error) {
	if err != nil {
		common.Throw(rt, err)
	}
}

// DynamicRead function helps to read dynamically when you don't know the size of []byte you would receive
func DynamicRead(read func([]byte) (int, error), timeout time.Duration) (int, []byte, error) {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	total := 0
	buffer := bytes.NewBuffer(nil)
	for ctx.Err() == nil {
		// TODO: add receive chunk size?
		chunk := make([]byte, 8192)
		n, err := read(chunk)
		if n > 0 {
			total += n
			buffer.Write(chunk[:n])
		}
		if err != nil && err.Error() != "EOF" {
			return total, buffer.Bytes(), err
		}

		if n < 8192 {
			break
		}
	}

	return total, buffer.Bytes(), nil
}

// NewEventListener function helps you to create a fress EventListener
func NewEventListener(t string) *events.EventListener {
	return &events.EventListener{
		EventType: t,
	}
}
