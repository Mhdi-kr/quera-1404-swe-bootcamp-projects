package controller

import (
	"strconv"

	"example.com/authorization/internal/constants"
	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/repository"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleGetAllPosts(c fiber.Ctx) error {
	var req dto.ListPostsRequest
	var response dto.ListPostsResponse

	err := c.Bind().Query(&req)
	if err != nil {
		c.SendStatus(fiber.StatusBadRequest)
	}

	req.Sanitize()

	dps, err := ctrl.postSrv.ListPosts(c.Context(), domain.PostFilters{
		Page: req.Page,
		Size: req.Size,
	})

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for _, dp := range dps {
		response.Posts = append(response.Posts, dto.Post{
			Id:          int(dp.Id),
			CreatedAt:   dp.CreatedAt,
			UpdatedAt:   &dp.UpdatedAt,
			URL:         dp.URL,
			Description: dp.Description,
		})
	}

	return c.JSON(response)
}

func (ctrl Controller) HandleDeletePost(c fiber.Ctx) error {
	postIDStr := c.Params("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)

	if len(postIDStr) == 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = ctrl.postSrv.DeletePost(c.Context(), userID, postID)
	if err != nil {
		if err == repository.ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response{
				Error:   err.Error(),
				Message: "post not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Error:   err.Error(),
			Message: "could not delete post",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
