package session

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type RedisStore struct {
	// Redis client used to talk to redis server.
	Client *redis.Client
	// Used for key expiry time on redis.
	SessionDuration time.Duration
}

// NewRedisClient takes an address and returns a the pointer of new client.
func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

// NewRedisStore constructs a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	// initialize and return a new RedisStore struct
	return &RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
}

// Save saves the provided `sessionState` and associated SID to the store.
// The `sessionState` parameter is typically a pointer to a struct containing
// all the data you want to associated with the given SID.
func (rs *RedisStore) Save(sid SID, sessionState interface{}) error {
	// marshal the `sessionState` to JSON and save it in the redis database,
	// using `sid.getRedisKey()` for the key.
	// return any errors that occur along the way.
	j, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}

	rs.Client.Set(sid.getRedisKey(), j, rs.SessionDuration)
	return nil
}

// DeleteUser deletes all state data associated with the SID from the store.
func (rs *RedisStore) Delete(sid SID) error {
	// delete the data stored in redis for the provided SID
	rs.Client.Del(sid.getRedisKey())
	return nil
}

// Get populates `sessionState` with the data previously saved
// for the given SID
func (rs *RedisStore) Get(sid SID, sessionState interface{}) error {
	// get the previously-saved session state data from redis,
	// unmarshal it back into the `sessionState` parameter
	// and reset the expiry time, so that it doesn't get deleted until
	// the SessionDuration has elapsed.
	pipeline := rs.Client.Pipeline()
	pipe := pipeline.Get(sid.getRedisKey())
	pipeline.Expire(sid.getRedisKey(), rs.SessionDuration)

	if _, err := pipeline.Exec(); err != nil {
		return err
	}

	if s, err := pipe.Result(); err != nil {
		return err
	} else {
		if err = json.Unmarshal([]byte(s), sessionState); err != nil {
			// cannot unmarshal
			return err
		} else {
			return nil
		}
	}
}

// getRedisKey() returns the redis key to use for the SID
func (sid SID) getRedisKey() string {
	// convert the SID to a string and add the prefix "sid:" to keep
	// SID keys separate from other keys that might end up in this
	// redis instance
	return string(sid)
}