package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

// Client struct is the Client object type that users is going to use in js like this:
//
// const client = new Client();
// const response = await client.get('https://httpbin.test.k6.io/get');
//
// you can see more usage examples in js through examples dir.
type Client struct {
	// The http.Client struct to have all the functionalities of a http.Client in Client struct
	http.Client

	// id is the unique identifier for each instance of Client object in js script
	id string

	// Multiple vus in k6 can create multiple Client objects so we need to have access the vu Runtime, etc.
	Vu modules.VU

	M map[string]sobek.Value

	// Params is the way to config the global params for Client object to do requests.
	// params *Clientparams

	// eventListeners interfaces.EventListeners
}

var _ sobek.DynamicObject = &Client{}

// Delete func is the required implementation of Delete function by sobek.DynamicObject type
func (c *Client) Delete(k string) bool {
	delete(c.M, k)
	return true
}

// Get func is the required implementation of Get function by sobek.DynamicObject type
func (c *Client) Get(k string) sobek.Value {
	return c.M[k]
}

// Has func is the required implementation of Has function by sobek.DynamicObject type
func (c *Client) Has(k string) bool {
	_, exists := c.M[k]
	return exists
}

// Keys func is the required implementation of Keys function by sobek.DynamicObject type
func (c *Client) Keys() []string {
	keys := make([]string, 0, len(c.M))
	for k := range c.M {
		keys = append(keys, k)
	}
	return keys
}

// Set func is the required implementation of Set function by sobek.DynamicObject type
func (c *Client) Set(k string, val sobek.Value) bool {
	c.M[k] = val
	return true
}

// Init func defines data properties on Client struct as DynamicObject, Also can initialize other things.
func (c *Client) Init() {
	rt := c.Vu.Runtime()
	c.id = uuid.New().String()

	c.Set("get", rt.ToValue(c.getAsync))
}

// This function will do the actuall request and act as a wrapper to create custom Request and Response objects
func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.Do(req)
}

// This function will act as a wrapper for each request we want to create from any function with any method,
// so we can verify input of each function and wrap any global header defined in client, etc.
// NOT VERY USEFUL FOR NOW BUT IT WILL BE IN FUTURE.
func (c *Client) createRequest(
	method string,
	arg sobek.Value,
	body io.Reader,
) (*http.Request, error) {
	addDefault := func(req *http.Request) {
		// here we will add global/default settings for client, not included in this PR.
	}

	if v, ok := arg.Export().(string); ok {
		req, err := http.NewRequest(method, v, body)
		addDefault(req)

		return req, err
	}

	return nil, fmt.Errorf(
		"invalid input! couldn't make the request from argument: %+v",
		arg.Export())
}

// This is a temp function just to end purpose of this PR and will be moved to Response object in future
func (c *Client) parseBody(body io.Reader) ([]byte, error) {
	_, res, err := dynamicRead(body.Read, 1*time.Second)
	return res, err
}

// This function now just do simple get requests using url.
func (c *Client) getAsync(arg sobek.Value) *sobek.Promise {
	enqCallback := c.Vu.RegisterCallback()
	p, resolve, reject := c.Vu.Runtime().NewPromise()

	req, err := c.createRequest(http.MethodGet, arg, nil)
	if err != nil {
		enqCallback(func() error {
			if er := reject(err); er != nil {
				return er
			}
			return nil
		})
		return p
	}

	go func() {
		res, err := c.do(req)
		enqCallback(func() error {
			if err != nil {
				if er := reject(err); er != nil {
					return er
				}
			}
			// this is a temp behavior to showcase the behavior of get it will be replaced by Response object in future
			body, er := c.parseBody(res.Body)
			if er != nil {
				if e := reject(er); e != nil {
					return e
				}
			}
			if er = resolve(string(body)); er != nil {
				return er
			}
			return nil
		})
	}()

	return p
}
