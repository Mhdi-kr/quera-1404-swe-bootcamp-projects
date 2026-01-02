package service

import (
	"context"

	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/repository"
	"example.com/authorization/internal/repository/entity"
)

type PostService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return PostService{
		postRepo: postRepo,
	}
}

func (us PostService) CreateProfilePost(ctx context.Context, post domain.Post) (int64, error) {
	return us.postRepo.Insert(ctx, entity.Post{
		Description: post.Description,
		URL:         post.URL,
		UserID:      post.UserID,
	})
}

func (us PostService) ListProfilePosts(ctx context.Context, userID int64, filters domain.PostFilters) ([]domain.Post, error) {
	ps, err := us.postRepo.List(ctx, &userID, filters.Size, filters.Page)
	if err != nil {
		return make([]domain.Post, 0), err
	}

	return domain.NewPostsFromEntities(ps), nil
}

func (us PostService) ListPosts(ctx context.Context, filters domain.PostFilters) ([]domain.Post, error) {
	ps, err := us.postRepo.List(ctx, nil, filters.Size, filters.Page)
	if err != nil {
		return make([]domain.Post, 0), err
	}

	return domain.NewPostsFromEntities(ps), nil
}

func (us PostService) DeletePost(ctx context.Context, userID int64, postID int64) error {
	return us.postRepo.DeleteByID(ctx, userID, postID)
}
