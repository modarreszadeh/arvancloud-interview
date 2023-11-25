package main

import (
	"github.com/labstack/echo/v4"
	"github.com/modarreszadeh/arvancloud-interview/internal/handler"
	"github.com/modarreszadeh/arvancloud-interview/internal/setting"
)

var rateLimiter = handler.NewRateLimiter(&handler.RateLimiterConfig{Rate: 60,
	Burst: 10, BlackList: []string{}, WhiteList: []string{}})

func handleRequests(e *echo.Echo) {
	e.POST("/object-storage", handler.HandleStorageRequest)
}

func handleFilters(e *echo.Echo) {
	e.Use(rateLimiter.Use)
}

func main() {
	e := echo.New()
	handleFilters(e)
	handleRequests(e)
	e.Logger.Fatal(e.Start(setting.Port))
}
