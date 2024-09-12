package service

import (
	"context"

	"github.com/MumMumGoodBoy/review-service/internal/model"
	"github.com/MumMumGoodBoy/review-service/proto"
	"gorm.io/gorm"
)

var _ proto.ReviewServer = (*ReviewService)(nil)

type ReviewService struct {
	proto.UnimplementedReviewServer
	DB *gorm.DB
}
// CreateReview implements proto.ReviewServer.
func (r *ReviewService) CreateReview(ctx context.Context, review *proto.ReviewRequest) (*proto.Empty, error) {
	userReview := model.Review{
		RestaurantId: review.RestaurantId,
		UserId:       review.UserId,
		Rating:       review.Rating,
		Content:      review.Content,
	}

	if err := r.DB.Create(&userReview).Error; err != nil {
		return nil, err
	}
	
	return nil, nil
}

// DeleteReview implements proto.ReviewServer.
func (r *ReviewService) DeleteReview(context.Context, *proto.DeleteReviewRequest) (*proto.Empty, error) {
	panic("unimplemented")
}

// GetReview implements proto.ReviewServer.
func (r *ReviewService) GetReview(context.Context, *proto.GetReviewRequest) (*proto.ReviewResponse, error) {
	panic("unimplemented")
}

// GetReviewsByRestaurantId implements proto.ReviewServer.
func (r *ReviewService) GetReviewsByRestaurantId(context.Context, *proto.GetReviewsRequest) (*proto.GetReviewsResponse, error) {
	panic("unimplemented")
}

// UpdateReview implements proto.ReviewServer.
func (r *ReviewService) UpdateReview(context.Context, *proto.UpdateReviewRequest) (*proto.ReviewResponse, error) {
	panic("unimplemented")
}