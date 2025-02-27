package interfaces

import (
	"github.com/grafana/sobek"
)

// Params interface defines the common behavior/functionalities each object params exported on the main API should have.
// Client, Request, TCP, etc. are objects and params you pass to them while initializing them are object params.
type Params interface{}

// Object interface defines the common behavior/functionalities each object exported on the main API should have.
// Client, Request, TCP, etc. are objects
type Object interface {
	sobek.DynamicObject
	Define() error
	ParseParams(*sobek.Runtime, []sobek.Value) (Params, error)
}
