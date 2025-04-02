package rand

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var all = []int32("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var alphabet = []int32("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var numeric = []int32("0123456789")
var capsAll = []int32("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// AllCapsString ...
func AllCapsString(length int) string {
	bytes := make([]int32, length)
	for i := range bytes {
		bytes[i] = capsAll[rand.Intn(len(capsAll))]
	}
	return string(bytes)
}

// String ...
func String(length int) string {
	bytes := make([]int32, length)
	for i := range bytes {
		bytes[i] = all[rand.Intn(len(all))]
	}
	return string(bytes)
}

// WordString ...
func WordString(length int) string {
	bytes := make([]int32, length)
	for i := range bytes {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(bytes)
}

// NumericString ...
func NumericString(length int) string {
	bytes := make([]int32, length)
	for i := range bytes {
		bytes[i] = numeric[rand.Intn(len(numeric))]
	}
	return string(bytes)
}

func UUID() string {
	return uuid.NewString()
}
