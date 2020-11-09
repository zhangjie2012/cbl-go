/*
specific string generated. about UUID:

- version 1 (date-time and MAC address)
- version 2 (date-time and MAC address, DCE security version)
- Versions 3 and 5 (namespace name-based)
- version 4 (random)
*/

package cbl

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// GenUUIDV1 global unique id
func GenUUIDV1() string {
	uuid, _ := uuid.NewUUID()
	return uuid.String()
}

// GenUUIDV4 wikipedia said Randomly generated UUIDs have 122 random bits.
// One's annual risk of being hit by a meteorite is estimated to be one chance in 17 billion,
// that means the probability is about 0.00000000006 (6 x 10−11),
// equivalent to the odds of creating a few tens of trillions of UUIDs in a
// year and having one duplicate.
//
// (most business scene recommand)
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

func genVerifyCode(pool string, width int) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	code := make([]byte, 0, width)
	for i := 0; i < width; i++ {
		code = append(code, pool[r.Intn(len(pool))])
	}
	return string(code)
}

// GenVerifyCodeAny generate verify code, get from pool
func GenVerifyCodeAny(pool string, width int) string {
	return genVerifyCode(pool, width)
}

// GenVerifyCodeNumber generate verify code, all is number
func GenVerifyCodeNumber(width int) string {
	pool := "0123456789"
	return genVerifyCode(pool, width)
}

// GenVerifyCodeCommon generate verify code, letter + number
func GenVerifyCodeCommon(width int) string {
	pool := "abcdefghigklmnopqrstuvwxyz0123456789"
	return genVerifyCode(pool, width)
}
