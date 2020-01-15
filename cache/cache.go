package cache // import jdel.org/gosspks/cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache contains packages cache
var Cache *cache.Cache

// DefaultCacheDuration is the default value
var DefaultCacheDuration time.Duration

func init() {
	DefaultCacheDuration = 24 * time.Hour
	Cache = cache.New(DefaultCacheDuration, DefaultCacheDuration)
}
