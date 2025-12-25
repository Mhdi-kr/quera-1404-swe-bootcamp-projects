package controller

import (
	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/domain"
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
