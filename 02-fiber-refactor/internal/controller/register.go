package controller

import (
	"example.com/authorization/internal/controller/dto"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleRegister(c fiber.Ctx) error {
	var response dto.Response
	var req dto.RegisterRequest

	err := c.Bind().Body(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// call service
	return c.JSON(response)
}
