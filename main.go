package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/cars", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"hello": "world",
		})
	})

	e.POST("/cars", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"hello": "world",
		})

	})

	e.PATCH("/cars/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(200, map[string]interface{}{
			"hello": "world",
			"id":    id,
		})
	})

	e.DELETE("/cars/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(200, map[string]interface{}{
			"hello":   "world",
			"message": "cars " + id + " has been deleted.",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
