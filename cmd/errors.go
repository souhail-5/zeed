package cmd

import (
	"errors"
	"strings"
)

// An errInitBus is where to put errors that occurred during initialization.
// When the real execution begins, the code can exit gracefully if errors are found in the bus.
type errInitBus struct {
	errors []error
}

func newErrInitBus() *errInitBus {
	return &errInitBus{}
}

func (b *errInitBus) NewError(s string) {
	b.errors = append(b.errors, errors.New(s))
}

func (b *errInitBus) AppendError(e error) {
	b.errors = append(b.errors, e)
}

func (b *errInitBus) Error() string {
	var s []string
	for _, err := range b.errors {
		s = append(s, err.Error())
	}

	return strings.Join(s, "\n")
}
