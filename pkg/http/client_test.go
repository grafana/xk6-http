package http

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils/httpmultibin"
	"go.k6.io/k6/metrics"
	"gopkg.in/guregu/null.v3"
)

// copied from xk6-websocket
type testState struct {
	tb      *httpmultibin.HTTPMultiBin
	runtime *modulestest.Runtime
	samples chan metrics.SampleContainer
	t       testing.TB

	errors chan error

	module *ModuleInstance
}

// copied from xk6-websocket
func newTestState(t testing.TB) testState {
	runtime := modulestest.NewRuntime(t)
	tb := httpmultibin.NewHTTPMultiBin(t)

	samples := make(chan metrics.SampleContainer, 1000)
	state := &lib.State{
		Dialer: tb.Dialer,
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
			UserAgent: null.StringFrom("TestUserAgent"),
		},
		Samples:        samples,
		TLSConfig:      tb.TLSClientConfig,
		BuiltinMetrics: runtime.BuiltinMetrics,
		Tags:           lib.NewVUStateTags(runtime.VU.InitEnvField.Registry.RootTagSet()),
	}

	m := new(RootModule).NewModuleInstance(runtime.VU)
	require.NoError(t, runtime.VU.RuntimeField.Set("Client", m.Exports().Named["Client"]))

	runtime.MoveToVUContext(state)
	return testState{
		runtime: runtime,
		tb:      tb,
		samples: samples,
		errors:  make(chan error, 50),
		t:       t,
		module:  m.(*ModuleInstance),
	}
}

func TestBasicGet(t *testing.T) {
	t.Parallel()
	ts := newTestState(t)
	sr := ts.tb.Replacer.Replace
	// TODO: await for get request, i don't know how to run it yet becuase it gives me unexpected token while running tests and using 'await' keyword
	_, err := ts.runtime.RunOnEventLoop(sr(`
  	const client = new Client();
  	const response = client.get('https://httpbin.test.k6.io/get');
	`))

	require.NoError(t, err)
}
