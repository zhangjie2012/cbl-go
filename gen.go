package cbl

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
)

// about UUID:
//   Version 1 (date-time and MAC address)
//   Version 2 (date-time and MAC address, DCE security version)
//   Versions 3 and 5 (namespace name-based)
//   Version 4 (random)

// global unique id
func GenUUIDV1() string {
	uuid, _ := uuid.NewUUID()
	return uuid.String()
}

// wikipedia
// Randomly generated UUIDs have 122 random bits.  One's annual risk of being
// hit by a meteorite is estimated to be one chance in 17 billion, that
// means the probability is about 0.00000000006 (6 x 10−11),
// equivalent to the odds of creating a few tens of trillions of UUIDs in a
// year and having one duplicate.
//
//  most business scene recommand
func GenUUIDV4() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

// GenSessionID generate a unique session id
func GenUSessionID() string {
	id := GenUUIDV1()
	h := sha256.New()
	h.Write([]byte(id))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GenRSessionID generate a random session id, most business scene recommand
func GenRSessionID() string {
	return GenUUIDV4()
}
