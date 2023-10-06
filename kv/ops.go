package kv

import (
	"errors"
	"time"
)

func Get(key string) (string, error) {
	return "", errors.ErrUnsupported
}

func Set(key, val string) error {
	return errors.ErrUnsupported
}

func Timestamps(key string) (first, last time.Time, err error) {
	return time.Time{}, time.Time{}, errors.ErrUnsupported
}

func Delete(key string) error {
	return errors.ErrUnsupported
}
