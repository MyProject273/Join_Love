package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	consts "github.com/MyProject273/Join_Love/pkg/const"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(consts.Alphabet)

	for i := 0; i < n; i++ {
		c := consts.Alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
