package service

import (
	"context"

	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepo
}

func NewCommentService(commentRepo repository.CommentRepo) CommentService {
	return CommentService{
		commentRepo: commentRepo,
	}
}

func (us CommentService) Create(ctx context.Context, comment domain.Comment) (int64, error) {
	return us.commentRepo.Insert(ctx, comment.ToEntity())
}

func (us CommentService) ListPostComments(ctx context.Context, postID int64, filters domain.CommentFilters) ([]domain.Comment, error) {
	cs, err := us.commentRepo.List(ctx, postID, filters.Size, filters.Page)
	if err != nil {
		return make([]domain.Comment, 0), err
	}

	return domain.NewCommentsFromEntities(cs), nil
}

func (us CommentService) Delete(ctx context.Context, userID int64, commentID int64) error {
	return us.commentRepo.DeleteByID(ctx, userID, commentID)
}

func (us CommentService) Upvote(ctx context.Context, userID int64, commentID int64) (bool, error) {
	return us.commentRepo.Upvote(ctx, userID, commentID)
}
