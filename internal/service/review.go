package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MumMumGoodBoy/review-service/internal/model"
	"github.com/MumMumGoodBoy/review-service/proto"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

var _ proto.ReviewServer = (*ReviewService)(nil)

type ReviewService struct {
	proto.UnimplementedReviewServer
	DB              *gorm.DB
	RabbitMQChannel *amqp091.Channel
}

func (r *ReviewService) publishReviewEvent(data model.ReviewEvent, event string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling review event: %v", err)
	}

	err = r.RabbitMQChannel.Publish(
		"review_topic",
		event,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("error publishing review event: %v", err)
	}
	return nil
}

func (r *ReviewService) publishFavoriteEvent(data model.FavoriteEvent, event string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling favorite event: %v", err)
	}

	err = r.RabbitMQChannel.Publish(
		"favorite_topic",
		event,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("error publishing favorite event: %v", err)
	}
	return nil
}

// CreateReview implements proto.ReviewServer.
func (r *ReviewService) CreateReview(ctx context.Context, review *proto.ReviewRequest) (*proto.Empty, error) {
	userReview := model.Review{
		RestaurantId: review.RestaurantId,
		UserId:       uint(review.UserId),
		Rating:       review.Rating,
		Content:      review.Content,
	}

	if err := r.DB.Create(&userReview).Error; err != nil {
		return nil, err
	}

	event := model.ReviewEvent{
		Event:        "review.create",
		Id:           int(userReview.ID),
		RestaurantId: userReview.RestaurantId,
		ReviewerId:   int(userReview.UserId),
		Rating:       userReview.Rating,
		Content:      userReview.Content,
	}
	if err := r.publishReviewEvent(event, "review.create"); err != nil {
		fmt.Println("Error publishing create review event: ", err)
	}

	return &proto.Empty{}, nil
}

// DeleteReview implements proto.ReviewServer.
func (r *ReviewService) DeleteReview(ctx context.Context, req *proto.DeleteReviewRequest) (*proto.Empty, error) {
	reviewId := req.ReviewId

	// Retrieve the to-be-deleted review
	var userReview model.Review
	if err := r.DB.First(&userReview, reviewId).Error; err != nil {
		// If the review is not found, return an error
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("[DeleteReview]: review not found")
		}
		return nil, err
	}

	if err := r.DB.Delete(&model.Review{}, reviewId).Error; err != nil {
		return nil, err
	}

	event := model.ReviewEvent{
		Event:        "review.delete",
		Id:           int(userReview.ID),
		RestaurantId: userReview.RestaurantId,
		ReviewerId:   int(userReview.UserId),
		Rating:       userReview.Rating,
		Content:      userReview.Content,
	}
	if err := r.publishReviewEvent(event, "review.delete"); err != nil {
		fmt.Println("Error publishing delete review event: ", err)
	}

	return &proto.Empty{}, nil
}

// GetReview implements proto.ReviewServer.
func (r *ReviewService) GetReview(ctx context.Context, req *proto.GetReviewRequest) (*proto.ReviewResponse, error) {
	var review model.Review
	if err := r.DB.First(&review, req.ReviewId).Error; err != nil {
		return nil, err
	}

	return &proto.ReviewResponse{
		ReviewId:     fmt.Sprintf("%d", review.ID),
		RestaurantId: review.RestaurantId,
		UserId:       int32(review.UserId),
		Rating:       review.Rating,
		Content:      review.Content,
	}, nil
}

// GetReviewsByRestaurantId implements proto.ReviewServer.
func (r *ReviewService) GetReviewsByRestaurantId(ctx context.Context, req *proto.GetReviewsRequest) (*proto.GetReviewsResponse, error) {
	var reviews []model.Review
	if err := r.DB.Where("restaurant_id = ?", req.RestaurantId).Find(&reviews).Error; err != nil {
		return nil, err
	}

	response := &proto.GetReviewsResponse{}
	for _, review := range reviews {
		response.Reviews = append(response.Reviews, &proto.ReviewResponse{
			ReviewId:     fmt.Sprintf("%d", review.ID),
			RestaurantId: review.RestaurantId,
			UserId:       int32(review.UserId),
			Rating:       review.Rating,
			Content:      review.Content,
		})
	}

	return response, nil
}

// UpdateReview implements proto.ReviewServer.
func (r *ReviewService) UpdateReview(ctx context.Context, req *proto.UpdateReviewRequest) (*proto.ReviewResponse, error) {
	var review model.Review

	if err := r.DB.First(&review, req.ReviewId).Error; err != nil {
		return nil, err
	}
	review.Rating = req.Rating
	review.Content = req.Content
	if err := r.DB.Save(&review).Error; err != nil {
		return nil, err
	}

	// send event to rabbitmq
	event := model.ReviewEvent{
		Event:        "review.update",
		Id:           int(review.ID),
		RestaurantId: review.RestaurantId,
		ReviewerId:   int(review.UserId),
		Rating:       review.Rating,
		Content:      review.Content,
	}
	if err := r.publishReviewEvent(event, "review.update"); err != nil {
		fmt.Println("Error publishing update review event: ", err)
	}

	return &proto.ReviewResponse{
		ReviewId:     fmt.Sprintf("%d", review.ID),
		RestaurantId: review.RestaurantId,
		UserId:       int32(review.UserId),
		Rating:       review.Rating,
		Content:      review.Content,
	}, nil
}

// AddFavoriteFood implements proto.ReviewServer.
func (r *ReviewService) AddFavoriteFood(ctx context.Context, req *proto.AddFavoriteFoodRequest) (*proto.Empty, error) {
	userId := req.UserId
	foodId := req.FoodId

	var existingFavorite model.FavoriteFood
	if err := r.DB.Where("user_id = ? AND food_id = ?", userId, foodId).First(&existingFavorite).Error; err == nil {
		return nil, fmt.Errorf("[AddFavoriteFood]: favorite already exists")
	}

	favorite := model.FavoriteFood{
		UserId: uint(userId),
		FoodId: foodId,
	}
	if err := r.DB.Create(&favorite).Error; err != nil {
		return nil, fmt.Errorf("failed to add favorite food: %v", err)
	}

	// Publish event to RabbitMQ
	event := model.FavoriteEvent{
		Event:  "favorite.add",
		UserId: int(userId),
		FoodId: foodId,
	}
	if err := r.publishFavoriteEvent(event, "favorite.add"); err != nil {
		fmt.Printf("Error publishing add favorite event: %v", err)
	}

	return &proto.Empty{}, nil
}

// RemoveFavoriteFood implements proto.ReviewServer.
func (r *ReviewService) RemoveFavoriteFood(ctx context.Context, req *proto.RemoveFavoriteFoodRequest) (*proto.Empty, error) {
	userId := req.UserId
	foodId := req.FoodId

	var existingFavorite model.FavoriteFood
	if err := r.DB.Where("user_id = ? AND food_id = ?", userId, foodId).First(&existingFavorite).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("[RemoveFavoriteFood]: favorite not found for user ID [%d] and food ID [%s]", userId, foodId)
		}
		return nil, fmt.Errorf("[RemoveFavoriteFood]: error retrieving favorite food: %v", err)
	}

	if err := r.DB.Delete(&existingFavorite).Error; err != nil {
		return nil, fmt.Errorf("[RemoveFavoriteFood]: failed to remove favorite food: %v", err)
	}

	// Publish event to RabbitMQ
	event := model.FavoriteEvent{
		Event:  "favorite.remove",
		UserId: int(userId),
		FoodId: foodId,
	}
	if err := r.publishFavoriteEvent(event, "favorite.remove"); err != nil {
		fmt.Printf("Warning: failed to publish favorite remove event: %v\n", err)
	}

	return &proto.Empty{}, nil
}

// GetFavoriteFoodsByUserId implements proto.ReviewServer.
func (r *ReviewService) GetFavoriteFoodsByUserId(ctx context.Context, req *proto.GetFavoriteFoodsByUserIDRequest) (*proto.GetFavoriteFoodsByUserIDResponse, error) {
	userId := req.UserId
	var favoriteFoods []model.FavoriteFood

	if err := r.DB.Where("user_id = ?", userId).Find(&favoriteFoods).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve favorite foods for user ID %d: %v", req.UserId, err)
	}

	// extract food_id from food
	var foodIds []string
	for _, favorite := range favoriteFoods {
		foodIds = append(foodIds, favorite.FoodId)
	}

	return &proto.GetFavoriteFoodsByUserIDResponse{
		FoodIds: foodIds,
	}, nil
}
