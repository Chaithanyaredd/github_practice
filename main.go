package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/shorten", shortenURLHandler)
	e.GET("/metrics", getMetricsHandler)
	e.GET("/:shortURL", redirectURLHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
