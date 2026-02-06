package grpcserver

import (
	"context"

	"example.com/authorization/internal/domain"
	"example.com/authorization/internal/service"
	postv1 "example.com/authorization/protos-gen/post/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPageSize = 4
	maxPageSize     = 100
)

type PostServiceServer struct {
	postv1.UnimplementedPostServiceServer
	postSrv service.PostService
}

func (s *PostServiceServer) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	page, size := sanitizeListPostsRequest(req)

	posts, err := s.postSrv.ListPosts(ctx, domain.PostFilters{
		Page: page,
		Size: size,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "could not list posts")
	}

	response := &postv1.ListPostsResponse{
		Posts: make([]*postv1.Post, 0, len(posts)),
	}
	for i := range posts {
		response.Posts = append(response.Posts, postToProto(&posts[i]))
	}

	return response, nil
}

func sanitizeListPostsRequest(req *postv1.ListPostsRequest) (uint64, uint64) {
	var page uint64
	var size uint64
	if req != nil {
		page = req.GetPage()
		size = req.GetSize()
	}

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = defaultPageSize
	}

	if size > maxPageSize {
		size = maxPageSize
	}

	return page, size
}

func postToProto(post *domain.Post) *postv1.Post {
	if post == nil {
		return nil
	}

	var createdAt *timestamppb.Timestamp
	if !post.CreatedAt.IsZero() {
		createdAt = timestamppb.New(post.CreatedAt)
	}

	var updatedAt *timestamppb.Timestamp
	if !post.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(post.UpdatedAt)
	}

	return &postv1.Post{
		Id:               post.Id,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Url:              post.URL,
		Description:      post.Description,
		NumberOfComments: post.CommentsCount,
		NumberOfUpvotes:  post.VoteCount,
	}
}
