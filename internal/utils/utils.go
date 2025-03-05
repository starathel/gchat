package utils

import (
	"errors"
	"io"
)

func NilIfEOF(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
