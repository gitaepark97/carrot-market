package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func CreateRandomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

func CreateRandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func CreateRandomString(n int) string {
	var sb strings.Builder
	k := len(ALPHABET)

	for i := 0; i < n; i++ {
		c := ALPHABET[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func CreateRandomEmail() string {
	return fmt.Sprintf("%s@email.com", CreateRandomString(6))
}
