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

func (us PostService) ListProfilePosts(ctx context.Context, post domain.Post) (int64, error) {
	return us.postRepo.Insert(ctx, entity.Post{
		Description: post.Description,
		URL:         post.URL,
		UserID:      post.UserID,
	})
}
