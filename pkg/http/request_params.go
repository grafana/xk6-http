package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/grafana/sobek"

	"github.com/grafana/xk6-http/internal/interfaces"
)

// Requestparams struct is the options for request
type Requestparams struct {
	// headers represents the default headers client object would use to do requests.
	headers http.Header

	// url represents the default URL client object would use to do requests.
	url *url.URL
}

// ParseParams parses Request params and save them to it's instance
func (r *Request) ParseParams(rt *sobek.Runtime, args []sobek.Value) (interfaces.Params, error) {
	parsed := &Requestparams{
		headers: make(http.Header),
	}
	if len(args) == 0 {
		r.params = parsed
		return parsed, nil
	}
	if len(args) > 1 {
		return nil, fmt.Errorf(
			"you can't have multiple arguments when creating a new Request, but you've had %d args",
			len(args),
		)
	}

	rawParams := args[0]
	params := rawParams.ToObject(rt)
	for _, k := range params.Keys() {
		switch k {
		case "headers":
			headers := params.Get(k)
			if sobek.IsUndefined(headers) || sobek.IsNull(headers) {
				continue
			}
			headersObj := headers.ToObject(rt)
			if headersObj == nil {
				continue
			}
			for _, key := range headersObj.Keys() {
				parsed.headers.Set(key, headersObj.Get(key).String())
			}

		case "url":
			urlV := params.Get(k)
			if sobek.IsUndefined(urlV) || sobek.IsNull(urlV) {
				continue
			}
			if v, ok := urlV.Export().(string); ok {
				addr, err := url.Parse(v)
				if err != nil {
					return parsed, fmt.Errorf(
						"invalid url for Request: %s",
						v,
					)
				}
				parsed.url = addr
			} else {
				return parsed, fmt.Errorf(
					"invalid url for Request: %s",
					v,
				)
			}

		default:
			return parsed, fmt.Errorf(
				"unknown Request's option: %s",
				k,
			)
		}
	}

	r.params = parsed
	return &Requestparams{}, nil
}
