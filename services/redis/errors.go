package redis

import "errors"

var ErrTimedOut = errors.New("client connection timed out")
var ErrNotInitialized = errors.New("client is not yet initialized")
var ErrNotAvailable = errors.New("client is not available, see logs")
