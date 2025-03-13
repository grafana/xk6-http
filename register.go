// Package http implements the http-protocol js module for k6.
// That module is used to make http requests.
package http

import (
	"go.k6.io/k6/js/modules"

	xhttp "github.com/grafana/xk6-http/pkg/http"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/net/http", new(xhttp.RootModule))
}
