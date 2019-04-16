package session

import "github.com/pkg/errors"

// ErrStateNotFound is returned from Store.Get() when the requested
// session id was not found in the store
var ErrStateNotFound = errors.New("no session state was found in the session store")

// ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

// ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

// ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("invalid Session ID")