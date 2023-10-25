package kv

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testcase struct {
	name        string
	initialData map[string]string
	key         string
	expError    error
	targetEntry *Entry
	finalData   map[string]string
}

func setupTestDB(t *testing.T, entries map[string]string) string {
	dir := t.TempDir()
	assert.NoError(t, os.Chdir(dir))
	for key, value := range entries {
		os.WriteFile(key, []byte(value), fs.ModePerm)
	}
	return dir
}

func checkTestDB(t *testing.T, path string, entries map[string]string) {
	dir, err := os.ReadDir(path)
	assert.NoError(t, err)
	assert.Equal(t, len(entries), len(dir))
	for key, data := range entries {
		assert.FileExists(t, key)
		actual, err := os.ReadFile(key)
		assert.NoError(t, err)
		assert.Equal(t, data, string(actual))
	}
}

func TestRead(t *testing.T) {
	for _, tc := range []testcase{
		{
			name:        "TestReadOK",
			key:         "foo",
			initialData: map[string]string{"foo": "0,0,bar"},
			expError:    nil,
			targetEntry: &Entry{
				Value:       "bar",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(0, 0),
			},
			finalData: map[string]string{"foo": "0,0,bar"},
		},
		{
			name:        "TestReadEmptyValue",
			key:         "foo",
			initialData: map[string]string{"foo": "0,0,"},
			expError:    ErrBadFormat,
			targetEntry: nil,
			finalData:   map[string]string{"foo": "0,0,"},
		},
		{
			name:        "TestReadNotExist",
			key:         "foo",
			initialData: nil,
			expError:    ErrNotExist,
			targetEntry: nil,
			finalData:   nil,
		},
		{
			name:        "TestReadDelimiterInValue",
			key:         "foo",
			initialData: map[string]string{"foo": "0,0,ba,r"},
			expError:    nil,
			targetEntry: &Entry{
				Value:       "ba,r",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(0, 0),
			},
			finalData: map[string]string{"foo": "0,0,ba,r"},
		},
		{
			name:        "TestReadNoValue",
			key:         "foo",
			initialData: map[string]string{"foo": "0,0"},
			expError:    ErrBadFormat,
			targetEntry: nil,
			finalData:   map[string]string{"foo": "0,0"},
		},
		{
			name:        "TestReadNothing",
			key:         "foo",
			initialData: map[string]string{"foo": ""},
			expError:    ErrBadFormat,
			targetEntry: nil,
			finalData:   map[string]string{"foo": ""},
		},
		{
			name:        "TestReadNotTimestamp",
			key:         "foo",
			initialData: map[string]string{"foo": "0,alpha,bar"},
			expError:    ErrBadFormat,
			targetEntry: nil,
			finalData:   map[string]string{"foo": "0,alpha,bar"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := setupTestDB(t, tc.initialData)

			entry, err := read(tc.key)
			assert.Equal(t, tc.targetEntry, entry)
			assert.ErrorIs(t, err, tc.expError)

			checkTestDB(t, dir, tc.finalData)
		})
	}
}

func TestWrite(t *testing.T) {
	for _, tc := range []testcase{
		{
			name:        "TestWriteOK",
			initialData: nil,
			key:         "foo",
			targetEntry: &Entry{
				Value:       "bar",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  nil,
			finalData: map[string]string{"foo": "0,1,bar"},
		},
		{
			name:        "TestOverWrite",
			initialData: map[string]string{"foo": "0,0,bar"},
			key:         "foo",
			targetEntry: &Entry{
				Value:       "baz",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  nil,
			finalData: map[string]string{"foo": "0,1,baz"},
		},
		{
			name:        "TestWriteEmptyValue",
			initialData: nil,
			key:         "foo",
			targetEntry: &Entry{
				Value:       "",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  ErrBadFormat,
			finalData: nil,
		},
		{
			name:        "TestOverWriteEmptyValue",
			initialData: map[string]string{"foo": "0,0,bar"},
			key:         "foo",
			targetEntry: &Entry{
				Value:       "",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  ErrBadFormat,
			finalData: map[string]string{"foo": "0,0,bar"},
		},
		{
			name:        "TestWriteEmptyEntry",
			initialData: nil,
			key:         "foo",
			targetEntry: &Entry{},
			expError:    ErrBadFormat,
			finalData:   nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := setupTestDB(t, tc.initialData)

			err := write(tc.key, tc.targetEntry)
			assert.ErrorIs(t, err, tc.expError)

			checkTestDB(t, dir, tc.finalData)
		})
	}
}

func TestDelete(t *testing.T) {
	for _, tc := range []testcase{
		{
			name:        "TestDeleteOK",
			initialData: map[string]string{"foo": "0,1,bar"},
			key:         "foo",
			targetEntry: &Entry{
				Value:       "bar",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  nil,
			finalData: nil,
		},
		{
			name:        "TestIdempotentDelete",
			initialData: nil,
			key:         "foo",
			targetEntry: &Entry{
				Value:       "bar",
				FirstEdited: time.Unix(0, 0),
				LastEdited:  time.Unix(1, 0),
			},
			expError:  nil,
			finalData: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testDir := setupTestDB(t, tc.initialData)

			err := delete(tc.key)
			assert.ErrorIs(t, err, tc.expError)

			checkTestDB(t, testDir, tc.finalData)
		})
	}
}
