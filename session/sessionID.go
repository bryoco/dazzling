package session

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// signedLength is the full length of the signed session ID
// (ID portion plus signature)
// const signedLength = idLength + sha256.Size

// SID represents a valid, digitally-signed session ID.
// This is a base64 URL encoded string created from a byte slice
// where the first `idLength` bytes are cryptographically random
// bytes representing the unique session ID, and the remaining bytes
// are an HMAC hash of those ID bytes (i.e., a digital signature).
// The byte slice layout is like so:
// +-----------------------------------------------------+
// |...32 crypto random bytes...|HMAC hash of those bytes|
// +-----------------------------------------------------+

const (
	// InvalidSessionID represents an empty, invalid session ID
	InvalidSessionID = ""
	// idLength is the length of the ID portion
	idLength = 32
)

var (
	// SID for current session
	generatedSID = ""
)

type SID string

// NewSessionID creates and returns a new digitally-signed session ID,
// using `signingKey` as the HMAC signing key. An error is returned only
// if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SID, error) {
	// if `signingKey` is zero-length, return InvalidSessionID
	// and an error indicating that it may not be empty
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("not signing key")
	}

	// Generate a new digitally-signed SID by doing the following:
	// - create a byte slice where the first `idLength` of bytes
	//   are cryptographically random bytes for the new session ID,
	//   and the remaining bytes are an HMAC hash of those ID bytes,
	//   using the provided `signingKey` as the HMAC key.
	// - encode that byte slice using base64 URL Encoding and return
	//   the result as a SID type

	randByte, err := GenerateRandomBytes(idLength)
	if err != nil {
		return InvalidSessionID, err
	}

	combined := append(randByte, GenerateRandomHMAC(signingKey)...)
	sid := base64.URLEncoding.EncodeToString(combined)

	generatedSID = sid
	return SID(sid), nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b);
		err != nil {
		// Note that err == nil only if we read len(b) bytes.
		return nil, err
	} else {
		return b, nil
	}
}

func GenerateRandomHMAC(signingKey string) []byte {
	buf := new(bytes.Buffer)
	b := buf.Bytes()
	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write(b)

	return h.Sum(nil)
}

func DecodeSessionID(sid string) ([]byte, error) {
	d, err := base64.URLEncoding.DecodeString(sid)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func IdHasMutated(id string) bool {
	return generatedSID == id
}

// ValidateID validates the string in the `id` parameter
// using the `signingKey` as the HMAC signing key
// and returns an error if invalid, or a SID if valid
func ValidateID(id string, signingKey string) (SID, error) {

	// validate the `id` parameter using the provided `signingKey`.
	// base64 decode the `id` parameter, HMAC hash the
	// ID portion of the byte slice, and compare that to the
	// HMAC hash stored in the remaining bytes. If they match,
	// return the entire `id` parameter as a SID type.
	// If not, return InvalidSessionID and ErrInvalidID.

	//if !IdHasMutated(id) {
	//	return InvalidSessionID, errors.New("id has mutated")
	//}

	decodedID, err := DecodeSessionID(id)
	if err != nil {
		return InvalidSessionID, err
	}

	// ID portion of the byte slice
	decodedID = decodedID[idLength:]

	match := hmac.Equal(decodedID, GenerateRandomHMAC(signingKey))
	if match {
		return SID(id), nil
	} else {
		return InvalidSessionID, errors.New("hmac not equal")
	}
}
