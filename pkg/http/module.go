package http

import (
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/modules"
)

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
		"Client": mi.initClient,
	}}
}

func (mi *ModuleInstance) initClient(sc sobek.ConstructorCall) *sobek.Object {
	rt := mi.vu.Runtime()

	c := &Client{
		Vu: mi.vu,
		M:  make(map[string]sobek.Value),
	}

	c.Init()

	o := rt.NewDynamicObject(c)
	return o
}
