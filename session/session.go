package session

import (
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

// BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
// Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState State, w http.ResponseWriter) (SID, error) {
	// - create a new SessionID
	// - save the sessionState to the store
	//   where "<sessionID>" is replaced with the newly-created SessionID
	//   (note the constants declared for you above, which will help you avoid typos)
	// - add a header to the ResponseWriter that looks like this:
	//     "Authorization: Bearer <sessionID>"

	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	if err := store.Save(sid, sessionState); err != nil {
		return InvalidSessionID, err
	}

	w.Header().Add("Authorization", schemeBearer + string(sid))
	return sid, nil
}

// GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SID, error) {

	// get the value of the Authorization header,
	id := r.Header.Get(headerAuthorization)

	// or the "auth" query string parameter if no Authorization header is present,
	if len(id) == 0 {
		id = r.URL.Query().Get(paramAuthorization)
	}

	s := strings.Split(id, " ")
	switch len(s) {
	// got nothing; no scheme
	case 0, 1:
		return InvalidSessionID, ErrNoSessionID
	case 2:
		// invalid scheme
		if scheme := s[0]; scheme != "Bearer" {
			return InvalidSessionID, ErrInvalidScheme
		}

		// If it's valid, return the SessionID. If not return the validation error.
		id = s[len(s)-1]
		sid, err := ValidateID(id, signingKey)
		if err != nil {
			return InvalidSessionID, err
		} else {
			return sid, nil
		}
	// unexpected weird cases
	default:
		return InvalidSessionID, ErrInvalidScheme
	}
}

// GetState extracts the SessionID from the request,
// gets the associated state from the provided store into
// the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store RedisStore, sessionState interface{}) (SID, error) {
	// get the SessionID from the request, and get the data
	// associated with that SessionID from the store.
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	if err := store.Get(sid, sessionState); err != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}

// EndSession extracts the SessionID from the request,
// and deletes the associated data in the provided store, returning
// the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SID, error) {
	// get the SessionID from the request, and delete the
	// data associated with it in the store.

	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}

	if err := store.Delete(sid); err != nil {
		return InvalidSessionID, err
	}

	return sid, nil
}