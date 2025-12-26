package controller

import (
	"errors"
	"fmt"
	"strconv"

	"example.com/authorization/internal/constants"
	"example.com/authorization/internal/controller/dto"
	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/repository"
	"github.com/gofiber/fiber/v3"
)

func (ctrl Controller) HandleListComments(c fiber.Ctx) error {
	var req dto.ListCommentsRequest
	var response dto.ListCommentsResponse

	err := c.Bind().Query(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	req.Sanitize()

	postIDStr := c.Params("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil || len(postIDStr) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cs, err := ctrl.commentSrv.ListPostComments(c.Context(), postID, domain.CommentFilters{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	for _, dc := range cs {
		response.Comments = append(response.Comments, dc.ToDTO())
	}

	return c.JSON(response)
}

func (ctrl Controller) HandleCreateComment(c fiber.Ctx) error {
	var request dto.CreateCommentRequest

	err := c.Bind().Body(&request)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	postIDStr := c.Params("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil || len(postIDStr) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	commentID, err := ctrl.commentSrv.Create(c.Context(), domain.Comment{
		UserID:  userID,
		Content: request.Content,
		PostID:  postID,
	})

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Response{
		Message: fmt.Sprintf("comment created with id %d", commentID),
	})
}

func (ctrl Controller) HandleUpvoteComment(c fiber.Ctx) error {
	commentIdStr := c.Params("commentId")
	commentID, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil || len(commentIdStr) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	upvoteState, err := ctrl.commentSrv.Upvote(c.Context(), userID, commentID)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response{
				Message: "cannot upvote non-existing comment",
			})
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(dto.Response{
		Message: strconv.FormatBool(upvoteState),
	})
}

func (ctrl Controller) HandleDeleteComment(c fiber.Ctx) error {
	commentIdStr := c.Params("commentId")
	commentID, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil || len(commentIdStr) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, ok := c.Context().Value(constants.UsrIDContextKey).(int64)
	if !ok {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = ctrl.commentSrv.Delete(c.Context(), userID, commentID)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.Response{
				Message: "cannot delete non-existing comment",
			})
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
