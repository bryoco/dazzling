package session

// Store represents a session data store.
// This is an abstract interface that can be implemented
// against several different types of data stores. For example,
// session data could be stored in memory in a concurrent map,
// or more typically in a shared key/value server store like redis.
type Store interface {
	// Save saves the provided `sessionState` and associated SessionID to the store.
	// The `sessionState` parameter is typically a pointer to a struct containing
	// all the data you want to associated with the given SessionID.
	Save(sid SID, sessionState interface{}) error

	// Get populates `sessionState` with the data previously saved
	// for the given SessionID
	Get(sid SID, sessionState interface{}) error

	// DeleteUser deletes all state data associated with the SessionID from the store.
	Delete(sid SID) error
}
