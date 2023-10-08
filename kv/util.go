package kv

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ErrNotExist is returned when attempting to perform an operation on a key
// which does not exist.
var ErrNotExist = errors.New("key does not exist")

// ErrBadFormat is returned when a read file contains data which cannot be
// properly parsed into an Entry.
var ErrBadFormat = errors.New("malformed data")

type Entry struct {
	FirstEdited time.Time
	LastEdited  time.Time
	Value       string
}

// delimiter is the internal symbol used to separate fields of an Entry
// in its encoded, on-disk format.
const delimiter = ','

// read is a wrapper around readFrom() which reads a key
// from the current directory.
func read(key string) (*Entry, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(wd, key)
	return readFrom(path)
}

// readFrom is an internal function that gets the contents of a file in db, and
// parses it into an Entry if possible. This function is exposed to be testable
// and swappable with other implementations in the future.
func readFrom(path string) (*Entry, error) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNotExist
	} else if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := bytes.SplitN(b, []byte{delimiter}, 3)
	if len(raw) != 3 {
		return nil, fmt.Errorf("%w, missing value or timestamp %q", ErrBadFormat, b)
	}

	firstEdited, err := strconv.ParseInt(string(raw[0]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w, could not parse first-set timestamp %q, %w",
			ErrBadFormat, raw[0], err)
	}

	lastEdited, err := strconv.ParseInt(string(raw[1]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w, could not parse last-set timestamp %q, %w",
			ErrBadFormat, raw[1], err)
	}

	value := raw[2]
	if len(value) == 0 {
		return nil, fmt.Errorf("%w, missing value %q", ErrBadFormat, b)
	}

	return &Entry{
		FirstEdited: time.Unix(firstEdited, 0),
		LastEdited:  time.Unix(lastEdited, 0),
		Value:       string(value),
	}, nil
}

// write is a wrapper around writeTo() which writes a key to the current
// directory.
func write(key string, entry *Entry) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(wd, key)
	return writeTo(path, entry)
}

// readFrom is an internal function that serializes an Entry, and then
// writes it to db in a file named key. This function is exposed to be testable
// and swappable with other implementations in the future.
func writeTo(path string, entry *Entry) error {
	if len(entry.Value) == 0 {
		return fmt.Errorf("%w, storing empty \"\" is not allowed", ErrBadFormat)
	}

	data := fmt.Sprintf(
		"%d%c%d%c%s",
		entry.FirstEdited.Unix(),
		delimiter,
		entry.LastEdited.Unix(),
		delimiter,
		entry.Value)

	return os.WriteFile(path, []byte(data), os.ModePerm)
}

func delete(key string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(wd, key)
	return deleteFile(path)
}

func deleteFile(path string) error {
	err := os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}
