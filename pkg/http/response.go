package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"

	"github.com/grafana/xk6-http/internal/helpers"
	"github.com/grafana/xk6-http/internal/interfaces"
)

// Response object
type Response struct {
	*http.Response

	id string

	Vu modules.VU

	M map[string]sobek.Value

	Request *Request

	params *Responseparams
}

var _ interfaces.Object = &Response{}

func (r *Response) Delete(k string) bool {
	delete(r.M, k)
	return true
}

func (r *Response) Get(k string) sobek.Value {
	return r.M[k]
}

func (r *Response) Has(k string) bool {
	_, exists := r.M[k]
	return exists
}

func (r *Response) Keys() []string {
	keys := make([]string, 0, len(r.M))
	for k := range r.M {
		keys = append(keys, k)
	}
	return keys
}

func (r *Response) Set(k string, val sobek.Value) bool {
	r.M[k] = val
	return true
}

// Define func defines data properties on obj attatched to Client struct.
func (r *Response) Define() error {
	rt := r.Vu.Runtime()
	r.id = uuid.New().String()

	r.Set("json", rt.ToValue(r.jsonAsync))
	r.Set("request", rt.NewDynamicObject(r.Request))
	r.Set("id", rt.ToValue(r.id))
	return nil
}

type jsonResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	// TODO: when it's in []byte the output would be base64 encoded and when it's in string it's ugly in the output
	// make a workaround for it
	Body []byte `json:"body"`
}

func (r *Response) json() ([]byte, error) {
	// TODO: make timeout configurable
	_, body, err := helpers.DynamicRead(r.Body.Read, 1*time.Second)
	if err != nil {
		return []byte{}, err
	}
	err = r.Body.Close()
	if err != nil {
		// maybe do something better?
		log.Printf("Response body couldn't be closed!: %s\n", err.Error())
	}

	res := &jsonResponse{
		Header:     r.Header,
		Body:       body,
		Status:     r.Status,
		StatusCode: r.StatusCode,
	}
	return json.MarshalIndent(res, "", "    ")
}

func (r *Response) jsonAsync() *sobek.Promise {
	enqCallback := r.Vu.RegisterCallback()
	p, resolve, reject := r.Vu.Runtime().NewPromise()

	go func() {
		res, err := r.json()
		enqCallback(func() error {
			if err != nil {
				if er := reject(err); er != nil {
					return er
				}
			}
			if er := resolve(string(res)); er != nil {
				return er
			}
			return nil
		})
	}()

	return p
}
