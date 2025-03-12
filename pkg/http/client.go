package http

import (
	"net/http"

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
func (c *Client) Init() error {
	c.id = uuid.New().String()

	return nil
}
