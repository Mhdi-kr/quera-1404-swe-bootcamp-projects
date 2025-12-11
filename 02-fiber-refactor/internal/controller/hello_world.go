package controller

import "github.com/gofiber/fiber/v3"

func (ctrl Controller) HandleHello(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
