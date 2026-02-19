package cache

import (
	"sync"
	"time"
)

type M3UCache struct {
	UpstreamETag string
	FilteredETag string
	FilteredBody []byte
	Mu           sync.Mutex
}

type EPGCache struct {
	UpstreamETag string
	FilteredETag string
	FilteredBody []byte
	LastFetch    time.Time
	Mu           sync.Mutex
}
