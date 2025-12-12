package controller

import (
	"example.com/authorization/internal/controller/dto"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleListUsers(c fiber.Ctx) error {
	var response dto.ListUsersResponse

	dusers, err := ctrl.userSrv.List()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for _, du := range dusers {
		response.Users = append(response.Users, dto.NewUserFromDomain(du))
	}

	return c.JSON(response)
}
