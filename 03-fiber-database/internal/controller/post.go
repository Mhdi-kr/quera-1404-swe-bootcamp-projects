package controller

import "github.com/gofiber/fiber/v3"

func (ctrl Controller) HandleGetAllPosts(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
