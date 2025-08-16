package cache

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCacheKey TODO
func GenerateCacheKey() string {
	return fmt.Sprintf("cache_%d_%d", time.Now().Unix(), rand.Int63())
}
