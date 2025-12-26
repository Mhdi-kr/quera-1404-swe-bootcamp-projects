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

	err = ctrl.userSrv.Register(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response = dto.Response{
		Message: "ok",
		Error:   "",
	}

	return c.JSON(response)
}
