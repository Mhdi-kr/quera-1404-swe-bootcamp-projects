package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()
	app.Get("/*", static.New("./public"))
	app.Listen(":8081")
}
