package http

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateLongContent() ([]byte, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1e6))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, nBig.Int64())
	_, err = rand.Read(buf)
	return buf, err
}

func TestDynamicRead(t *testing.T) {
	t.Parallel()
	type testCase struct {
		readFn  func([]byte) (int, error)
		timeout time.Duration

		expectedLen int
		expectedRes []byte
		error       error
	}

	t.Run("Read long and random len content", func(t *testing.T) {
		t.Parallel()
		tests := []testCase{}
		for range 3 {
			content, err := generateLongContent()
			// TODO: maby this is wrong and should be handled another way?
			if err != nil {
				assert.Error(t, err)
			}
			reader := bytes.NewReader(content)

			tests = append(tests, testCase{
				readFn:  reader.Read,
				timeout: time.Second * 1,

				expectedRes: content,
				expectedLen: reader.Len(),
				error:       nil,
			})
		}

		for _, test := range tests {
			actualLen, actualRes, err := dynamicRead(test.readFn, test.timeout)
			t.Logf("The actual content len is %d", actualLen)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedLen, actualLen)
			assert.Equal(t, test.expectedRes, actualRes)
		}
	})
}
