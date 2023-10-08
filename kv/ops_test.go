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
// from each operation are compared to some expected values.
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
	err      error
	value    string
	firstSet time.Time
	lastSet  time.Time
}

var (
	time0 = time.Unix(0, 0)
	time1 = time.Unix(1, 0)
)

func testGet(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		value, err := kv.Get(key)
		assert.ErrorIs(t, err, r.err)
		assert.Equal(t, r.value, value)
	}
}

func testSet(key, val string, now time.Time) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		err := kv.Set(key, val, now)
		assert.ErrorIs(t, err, r.err)
	}
}

func testTimestamp(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		firstSet, lastSet, err := kv.Timestamps(key)
		assert.ErrorIs(t, err, r.err)
		assert.Equal(t, r.firstSet, firstSet, "firstSet")
		assert.Equal(t, r.lastSet, lastSet, "lastSet")
	}
}

func testDelete(key string) func(*testing.T, *expReturns) {
	return func(t *testing.T, r *expReturns) {
		err := kv.Delete(key)
		assert.ErrorIs(t, err, r.err)
	}
}

func TestKVOps(t *testing.T) {
	for _, tc := range []testcase{
		{
			name: "TestGetNotExist",
			steps: []step{
				{
					run:        testGet("a"),
					expReturns: expReturns{err: kv.ErrNotExist},
				},
			},
		},
		{
			name: "TestDoubleGetNotExist",
			steps: []step{
				{
					run:        testGet("a"),
					expReturns: expReturns{err: kv.ErrNotExist},
				},
				{
					run:        testGet("b"),
					expReturns: expReturns{err: kv.ErrNotExist},
				},
			},
		},
		{
			name: "TestSetGet",
			steps: []step{
				{
					run:        testSet("a", "12", time0),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testGet("a"),
					expReturns: expReturns{err: nil, value: "12"},
				},
				{
					run:        testSet("b", "13", time0),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testGet("b"),
					expReturns: expReturns{err: nil, value: "13"},
				},
				{
					run:        testSet("a", "14", time1),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testGet("a"),
					expReturns: expReturns{err: nil, value: "14"},
				},
			},
		},
		{
			name: "TestSetTS",
			steps: []step{
				{
					run:        testSet("a", "12", time0),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testTimestamp("a"),
					expReturns: expReturns{err: nil, firstSet: time0, lastSet: time0},
				},
				{
					run:        testSet("b", "13", time1),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testTimestamp("b"),
					expReturns: expReturns{err: nil, firstSet: time1, lastSet: time1},
				},
				{
					run:        testSet("a", "13", time1),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testTimestamp("a"),
					expReturns: expReturns{err: nil, firstSet: time0, lastSet: time1},
				},
			},
		},
		{
			name: "TestDeleteTS",
			steps: []step{
				{
					run:        testSet("a", "12", time0),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testDelete("a"),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testTimestamp("a"),
					expReturns: expReturns{err: kv.ErrNotExist},
				},
			},
		},
		{
			name: "TestDeleteGet",
			steps: []step{
				{
					run:        testSet("a", "12", time0),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testDelete("a"),
					expReturns: expReturns{err: nil},
				},
				{
					run:        testGet("a"),
					expReturns: expReturns{err: kv.ErrNotExist},
				},
			},
		},
	} {
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
