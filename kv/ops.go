package kv

import (
	"errors"
	"fmt"
	"time"
)

func Get(key string) (string, error) {
	entry, err := read(key)
	if err != nil {
		return "", fmt.Errorf("error getting key %s, %w", key, err)
	}
	return entry.Value, nil
}

func Set(key, val string, now time.Time) error {
	entry := Entry{
		LastEdited: now,
		Value:      val,
	}

	if oldEntry, err := read(key); errors.Is(err, ErrNotExist) {
		entry.FirstEdited = entry.LastEdited
	} else {
		entry.FirstEdited = oldEntry.FirstEdited
	}

	err := write(key, &entry)
	if err != nil {
		return fmt.Errorf("error setting key %s, %w", key, err)
	}
	return nil
}

func Timestamps(key string) (first, last time.Time, err error) {
	entry, err := read(key)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf(
			"error getting timestamps for %s, %w", key, err)
	}
	return entry.FirstEdited, entry.LastEdited, nil
}

func Delete(key string) error {
	err := delete(key)
	if err != nil {
		return fmt.Errorf("error deleting key %s, %w", key, err)
	}
	return nil
}
