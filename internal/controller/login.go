package controller

import (
	"errors"

	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/service"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleLogin(c fiber.Ctx) error {
	var response dto.LoginResponse
	var request dto.LoginRequest

	err := c.Bind().Body(&request)
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := ctrl.userSrv.Login(c.Context(), request.Username, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response = dto.LoginResponse{
		Token: string(token),
		Response: dto.Response{
			Message: "ok",
			Error:   "",
		},
	}

	return c.JSON(response)
}
