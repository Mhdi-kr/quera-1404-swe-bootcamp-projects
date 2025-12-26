package controller

import (
	"example.com/authorization/internal/constants"
	"example.com/authorization/internal/controller/dto"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleSelf(c fiber.Ctx) error {
	var response dto.SelfResponse

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	user, err := ctrl.userSrv.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	response = dto.SelfResponse{
		User:   user.ToDTO(),
		Status: "ok",
	}
	// TODO: call service
	return c.Status(fiber.StatusOK).JSON(response)
}
