package main

import "github.com/labstack/echo/v4"

func main() {
	app := echo.New()
	app.GET("/", handleHome)
	app.Start(":1234")
}

func handleHome(c echo.Context) error {
	return c.String(200, "Hello, World!")
}