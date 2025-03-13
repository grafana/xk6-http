package http

import (
	"testing"

	"go.uber.org/goleak"
)

// copied from xk6-websocket
func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
