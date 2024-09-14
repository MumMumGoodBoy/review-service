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
	
	return &proto.Empty{}, nil
}

// DeleteReview implements proto.ReviewServer.
func (r *ReviewService) DeleteReview(ctx context.Context, req *proto.DeleteReviewRequest) (*proto.Empty, error) {
	reviewId := req.ReviewId

	if err := r.DB.Delete(&model.Review{}, reviewId).Error; err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// GetReview implements proto.ReviewServer.
func (r *ReviewService) GetReview(ctx context.Context,req *proto.GetReviewRequest) (*proto.ReviewResponse, error) {
	var review model.Review
	if err := r.DB.First(&review, req.ReviewId).Error; err != nil {
		return nil, err
	}

	return &proto.ReviewResponse{
		ReviewId:    	review.ID,
		RestaurantId:	review.RestaurantId,
		UserId:      	review.UserId,
		Rating:      	review.Rating,
		Content:     	review.Content,
	}, nil
}

// GetReviewsByRestaurantId implements proto.ReviewServer.
func (r *ReviewService) GetReviewsByRestaurantId(ctx context.Context,req *proto.GetReviewsRequest) (*proto.GetReviewsResponse, error) {
	var reviews []model.Review
	if err := r.DB.Where("restaurant_id = ?", req.RestaurantId).Find(&reviews).Error; err != nil {
		return nil, err
	}

	response := &proto.GetReviewsResponse{}
	for _, review := range reviews {
		response.Reviews = append(response.Reviews, &proto.ReviewResponse{
			ReviewId:    review.ID,
			RestaurantId: review.RestaurantId,
			UserId:      review.UserId,
			Rating:      review.Rating,
			Content:     review.Content,
		})
	}

	return response, nil
}

// UpdateReview implements proto.ReviewServer.
func (r *ReviewService) UpdateReview(ctx context.Context,req *proto.UpdateReviewRequest) (*proto.ReviewResponse, error) {
	var review model.Review

	if err := r.DB.First(&review, req.ReviewId).Error; err != nil {
		return nil, err
	}
	review.Rating = req.Rating
	review.Content = req.Content
	if err := r.DB.Save(&review).Error; err != nil {
		return nil, err
	}
	return &proto.ReviewResponse{
		ReviewId:    review.ID,
		RestaurantId: review.RestaurantId,
		UserId:      review.UserId,
		Rating:      review.Rating,
		Content:     review.Content,
	}, nil
}