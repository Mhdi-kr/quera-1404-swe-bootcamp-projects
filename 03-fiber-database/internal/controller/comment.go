package controller

import "github.com/gofiber/fiber/v3"

func (ctrl Controller) HandleCreateComment(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func (ctrl Controller) HandleUpvoteComment(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
func (ctrl Controller) HandleDeleteComment(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
