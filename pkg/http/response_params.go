package http

import (
	"github.com/grafana/sobek"

	"github.com/grafana/xk6-http/internal/interfaces"
)

// Responseparams struct is the default global options for Client struct
type Responseparams struct{}

var _ interfaces.Params = &Responseparams{}

func (r *Response) ParseParams(rt *sobek.Runtime, args []sobek.Value) (interfaces.Params, error) {
	r.params = &Responseparams{}
	return &Responseparams{}, nil
}
