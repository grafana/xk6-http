package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"

	"github.com/grafana/xk6-http/internal/helpers"
	"github.com/grafana/xk6-http/internal/interfaces"
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

	id string

	// Multiple vus in k6 can create multiple Client objects so we need to have access the vu Runtime, etc.
	Vu modules.VU

	M map[string]sobek.Value

	// Params is the way to config the global params for Client object to do requests.
	params *Clientparams

	eventListeners interfaces.EventListeners
}

var _ interfaces.Object = &Client{}

func (c *Client) Delete(k string) bool {
	delete(c.M, k)
	return true
}

func (c *Client) Get(k string) sobek.Value {
	return c.M[k]
}

func (c *Client) Has(k string) bool {
	_, exists := c.M[k]
	return exists
}

func (c *Client) Keys() []string {
	keys := make([]string, 0, len(c.M))
	for k := range c.M {
		keys = append(keys, k)
	}
	return keys
}

func (c *Client) Set(k string, val sobek.Value) bool {
	c.M[k] = val
	return true
}

// Define func defines data properties on obj attatched to Client struct.
func (c *Client) Define() error {
	rt := c.Vu.Runtime()
	c.eventListeners = (&eventListeners{}).New()
	c.id = uuid.New().String()

	c.Set("get", rt.ToValue(c.getAsync))
	c.Set("on", rt.ToValue(c.addEventListener))
	return nil
}

// this function add an eventListener of type t and an fn callback to the object eventListeners
func (c *Client) addEventListener(t string, fn func(sobek.Value) (sobek.Value, error)) error {
	el, err := c.eventListeners.GetListener(t)
	if err != nil {
		return err
	}
	el.Add(fn)

	return nil
}

// this function call eventListeners of any type
func (c *Client) callEventListeners(t string, obj *sobek.Object) error {
	el, err := c.eventListeners.GetListener(t)
	if err != nil {
		return err
	}
	for _, fn := range el.All() {
		if _, err := fn(obj); err != nil {
			return err
		}
	}

	return nil
}

// this function queueResponse for handling on the main event loop that the VU has
func (c *Client) queueResponse(resp *Response) {
	rt := c.Vu.Runtime()
	enqCallback := c.Vu.RegisterCallback()

	go func() {
		enqCallback(func() error {
			return c.callEventListeners(RESPONSE, rt.NewDynamicObject(resp))
		})
	}()
}

// this function would handle any type of request and do the actuall job of requesting
func (c *Client) do(req *Request) (*Response, error) {
	rt := c.Vu.Runtime()

	resp := &Response{
		Vu:      c.Vu,
		M:       make(map[string]sobek.Value),
		Request: req,
	}

	httpResp, err := c.Do(req.Request)
	if err != nil {
		return resp, err
	}

	resp.Response = httpResp
	helpers.Must(rt, resp.Define())

	c.queueResponse(resp)

	// TODO: for now i don't set any property for response to access to it on js but i think it should be done in future
	req.Response = resp

	return resp, nil
}

// this function would handle creating request with params from input
func (c *Client) createRequest(method string, arg sobek.Value, body io.Reader) (*Request, error) {
	rt := c.Vu.Runtime()
	// add default options to requests function
	addDefault := func(req *Request) {
		for k, vlist := range c.params.headers {
			if len(vlist) == 0 {
				continue
			}
			for _, v := range vlist {
				req.Header.Add(k, v)
			}
		}
	}

	// if the input is an req object then everything has been set before so we just add defaults and return
	if v, ok := arg.Export().(*Request); ok {
		addDefault(v)
		return v, nil
	}

	if v, ok := arg.Export().(string); ok {
		r, err := http.NewRequest(method, v, body)
		req := &Request{
			Vu:      c.Vu,
			M:       make(map[string]sobek.Value),
			Request: r,
		}
		helpers.Must(rt, req.Define())

		addDefault(req)
		return req, err
	}

	return &Request{}, fmt.Errorf(
		"invalid input! couldn't make the request from argument: %+v",
		arg.Export())
}

func (c *Client) getAsync(arg sobek.Value) *sobek.Promise {
	rt := c.Vu.Runtime()

	enqCallback := c.Vu.RegisterCallback()
	p, resolve, reject := c.Vu.Runtime().NewPromise()

	req, err := c.createRequest(http.MethodGet, arg, nil)
	if err != nil {
		// TODO: find a way to handle the rejection error
		_ = reject(err)
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
			if er := resolve(rt.NewDynamicObject(res)); er != nil {
				return er
			}
			return nil
		})
	}()

	return p
}
