//go:build !golangdesign

package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func Unsupported() bool {
	return clipboard.Unsupported
}

func Write(data []byte) error {
	if Unsupported() {
		return ErrUnsupported
	}
	if err := clipboard.WriteAll(string(data)); err != nil {
		return fmt.Errorf("clipboard.WriteAll: %w", err)
	}
	return nil
}

func Read() ([]byte, error) {
	if Unsupported() {
		return nil, ErrUnsupported
	}
	data, err := clipboard.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("clipboard.ReadAll: %w", err)
	}
	return []byte(data), nil
}
