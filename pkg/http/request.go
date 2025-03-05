package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"

	"github.com/grafana/xk6-http/internal/interfaces"
)

// Request object
type Request struct {
	*http.Request

	id string

	Vu modules.VU

	M map[string]sobek.Value

	// Response *Response

	params *Requestparams
}

var _ interfaces.Object = &Request{}

func (r *Request) Delete(k string) bool {
	delete(r.M, k)
	return true
}

func (r *Request) Get(k string) sobek.Value {
	return r.M[k]
}

func (r *Request) Has(k string) bool {
	_, exists := r.M[k]
	return exists
}

func (r *Request) Keys() []string {
	keys := make([]string, 0, len(r.M))
	for k := range r.M {
		keys = append(keys, k)
	}
	return keys
}

func (r *Request) Set(k string, val sobek.Value) bool {
	r.M[k] = val
	return true
}

// Define func defines data properties on obj attatched to Request struct.
func (r *Request) Define() error {
	rt := r.Vu.Runtime()

	r.id = uuid.New().String()
	if r.params != nil {
		r.Request.URL = r.params.url
		r.Request.Header = r.params.headers.Clone()
	}
	r.Set("id", rt.ToValue(r.id))
	return nil
}
