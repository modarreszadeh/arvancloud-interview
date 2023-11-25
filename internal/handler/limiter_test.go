package handler

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	w1Ip = "10.0.0.1"
	w2Ip = "10.0.0.2"
	b1Ip = "10.0.0.3"
	b2Ip = "10.0.0.4"
)

func TestLimiter(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:      60000,
		Burst:     10,
		WhiteList: []string{w1Ip, w2Ip},
		BlackList: []string{b1Ip, b2Ip},
	}
	rl := NewRateLimiter(&cfg)
	fail := 0
	success := 0
	for i := 0; i < 20; i++ {
		if rl.Allow("2") == false {
			fail++
		} else {
			success++
		}
	}
	require.Equal(t, 9, fail)
	require.Equal(t, 11, success)

	require.False(t, rl.Allow(b1Ip))
	require.False(t, rl.Allow(b2Ip))
}
