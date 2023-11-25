package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/modarreszadeh/arvancloud-interview/pkg/constants"
	"github.com/modarreszadeh/arvancloud-interview/pkg/helper"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type RateLimit struct {
	Limiters  *cache.Cache
	MaxTime   time.Duration
	TimeStep  time.Duration
	Config    *RateLimiterConfig
	WhiteList map[string]interface{}
	BlackList map[string]interface{}
}

type RateLimiterConfig struct {
	Burst     int      `json:"burst"`
	Rate      int      `json:"rate"`
	WhiteList []string `json:"whitelist"`
	BlackList []string `json:"blacklist"`
}

type Limiter struct {
	Size       time.Duration
	LastUpdate time.Time
	Mutex      *sync.Mutex
}

func NewRateLimiter(config *RateLimiterConfig) *RateLimit {
	rl := &RateLimit{
		Config: config,
	}
	rl.Limiters = cache.New(time.Minute, time.Minute*10)
	rl.WhiteList = make(map[string]interface{})
	for _, x := range config.WhiteList {
		rl.WhiteList[x] = nil
	}
	rl.BlackList = make(map[string]interface{})
	for _, x := range config.BlackList {
		rl.BlackList[x] = nil
	}
	rl.TimeStep = time.Duration(60000/config.Rate) * time.Millisecond
	rl.MaxTime = rl.TimeStep * time.Duration(config.Burst)
	return rl
}

func (rl *RateLimit) Allow(key string) bool {

	if _, exist := rl.BlackList[key]; exist {
		return false
	}
	if _, exist := rl.WhiteList[key]; exist {
		return true
	}
	var (
		res bool
		l   *Limiter
	)
	value, found := rl.Limiters.Get(key)
	if found {
		l = value.(*Limiter)
		l.Mutex.Lock()

		l.Size -= time.Since(l.LastUpdate)
		l.LastUpdate = time.Now()

		if l.Size < 0 {
			l.Size = 0
		}
		if l.Size > rl.MaxTime {
			res = false
		} else {
			l.Size += rl.TimeStep
			res = true
		}
		rl.Limiters.Set(key, l, time.Minute)
		l.Mutex.Unlock()
	} else {
		l = &Limiter{
			Size:       rl.TimeStep,
			LastUpdate: time.Now(),
			Mutex:      &sync.Mutex{},
		}
		res = true
		rl.Limiters.Set(key, l, time.Minute)
	}
	return res
}

func (rl *RateLimit) Use(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, _ := strconv.Atoi(helper.GetQueryParameters(c.QueryString())["userid"])
		if rl.Allow(strconv.Itoa(userId)) {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusBadRequest, constants.TooMeanyRequest)
	}
}
