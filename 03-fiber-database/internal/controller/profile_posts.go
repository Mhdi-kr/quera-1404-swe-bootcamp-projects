package controller

import (
	"example.com/authorization/internal/controller/dto"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleProfilePosts(c fiber.Ctx) error {
	var response dto.ProfilePostResponse

	// TODO: call service
	return c.Status(fiber.StatusOK).JSON(response)
}
