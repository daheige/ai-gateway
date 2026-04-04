package errs

import (
	"errors"
)

var (
	// ErrServerInternal server inner error
	ErrServerInternal = errors.New("server internal error")

	// ErrApiKeyNotFound apikey not found
	ErrApiKeyNotFound = errors.New("api key not found")
)
