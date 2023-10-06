package kv_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/adowair/kvdb/kv"
	"github.com/stretchr/testify/assert"
)

// testcase runs as a subtest of the top-level test in this file. Each
// testcase consists of a series of operations to be executed. The returns
// from each command are compared to some expected values.
type testcase struct {
	name  string
	steps []step
}

// step encapsulates the run of a single command, its expected outputs
// and the assertion that its outputs match the expectation.
type step struct {
	run func(*testing.T, *expReturns)
	expReturns
}

// expReturns encapsulates possible returns for any operation. All fields
// are optional depending on the operation being tested, except for `expError`.
type expReturns struct {
	expError    error
	expValue    string
	expFirstSet time.Time
	expLastSet  time.Time
}

func testGet(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		ret, err := kv.Get(key)
		assert.ErrorIs(t, err, r.expError)
		assert.Equal(t, ret, r.expValue)
	}
}

func testSet(key string, val string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		err := kv.Set(key, val)
		assert.ErrorIs(t, err, r.expError)
	}
}

func testTimestamp(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		first, last, err := kv.Timestamps(key)
		assert.ErrorIs(t, err, r.expError)
		assert.Equal(t, first, r.expFirstSet)
		assert.Equal(t, last, r.expLastSet)
	}
}

func testDelete(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		err := kv.Delete(key)
		assert.ErrorIs(t, err, r.expError)
	}
}

func TestKVOps(t *testing.T) {
	for _, tc := range []testcase{} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			assert.NoError(t, os.Chdir(dir))
			for i, step := range tc.steps {
				t.Run(fmt.Sprintf("TestStep%d", i), func(t *testing.T) {
					step.run(t, &step.expReturns)
				})
			}
		})
	}
}
