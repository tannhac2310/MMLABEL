package idutil

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// ULID generates a string representation of a ULID.
// rand.New is not threadsafe, so we create a pool of rand to speed up the id generation.
var randPool = sync.Pool{
	New: func() interface{} {
		// Note that this implementation to create the entropy:
		//	t := time.Now()
		//	return ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		// is not correct if there have *extremely* multiple concurrent calls.
		return ulid.Monotonic(rand.Reader, 0)
	},
}

func ULID(t time.Time) string {
	entropy := randPool.Get().(*ulid.MonotonicEntropy)
	defer randPool.Put(entropy)

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// ULIDNow returns a new ULID.
func ULIDNow() string {
	return ULID(time.Now())
}

// RandomString generates a random alphanumeric string of the given length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomByte, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randomByte.Int64()]
	}
	return string(b)
}
