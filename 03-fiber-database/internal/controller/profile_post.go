package controller

import (
	"fmt"

	"example.com/authorization/internal/constants"
	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/domain"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleListProfilePosts(c fiber.Ctx) error {
	var response dto.ProfilePostResponse

	// TODO: call service
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ctrl Controller) HandleCreateProfilePost(c fiber.Ctx) error {
	var request dto.CreateProfilePostRequest

	err := c.Bind().Body(&request)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	postID, err := ctrl.postSrv.Create(c.Context(), domain.Post{
		Description: request.Description,
		URL:         request.URL,
		UserID:      userID,
	})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// TODO: call service
	return c.Status(fiber.StatusOK).JSON(dto.Response{
		Message: fmt.Sprintf("post created with Id %d", postID),
	})
}
