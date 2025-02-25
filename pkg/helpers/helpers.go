package helpers

import (
	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
)

// Must panics and throw error on rt if err is not nil
func Must(rt *sobek.Runtime, err error) {
	if err != nil {
		common.Throw(rt, err)
	}
}
