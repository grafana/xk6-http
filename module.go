// Package http implements the http-protocol js module for k6.
// That module is used to make http requests.
package http

import (
	"net/http"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"

	"github.com/grafana/xk6-http/internal/helpers"
	xhttp "github.com/grafana/xk6-http/pkg/http"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/net/http", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu modules.VU
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new
// instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu: vu,
	}
}

// Exports implements the [modules.Instance] interface and returns the exports
// of the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Named: map[string]interface{}{
		"Client":  mi.initClient,
		"Request": mi.initRequest,
	}}
}

func (mi *ModuleInstance) initClient(sc sobek.ConstructorCall) *sobek.Object {
	rt := mi.vu.Runtime()

	c := &xhttp.Client{
		Vu: mi.vu,
		M:  make(map[string]sobek.Value),
	}

	helpers.Must(rt, func() error {
		_, err := c.ParseParams(rt, sc.Arguments)
		return err
	}())
	helpers.Must(rt, c.Define())

	o := rt.NewDynamicObject(c)
	return o
}

func (mi *ModuleInstance) initRequest(sc sobek.ConstructorCall) *sobek.Object {
	rt := mi.vu.Runtime()

	r := &xhttp.Request{
		Vu:      mi.vu,
		M:       make(map[string]sobek.Value),
		Request: &http.Request{},
	}

	helpers.Must(rt, func() error {
		_, err := r.ParseParams(rt, sc.Arguments)
		return err
	}())
	helpers.Must(rt, r.Define())

	// TODO: find another way to reconstruct the original Object cause this way we cannot implement other functionality to the object
	o := rt.NewDynamicObject(r)
	return o
}
