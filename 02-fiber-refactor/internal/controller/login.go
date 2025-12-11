package controller

import (
	"example.com/authorization/internal/controller/dto"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleLogin(c fiber.Ctx) error {
	var response dto.LoginResponse
	var request dto.LoginRequest

	err := c.Bind().Body(&request)
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
	}

	// TODO: call service

	return c.JSON(response)
}
