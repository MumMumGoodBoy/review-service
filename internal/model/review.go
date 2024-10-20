package model

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	RestaurantId string
	UserId       uint
	Rating       float32
	Content      string
}

type FavoriteFood struct {
	gorm.Model
	UserId       uint
	FoodId       string
	RestaurantId string
}

type ReviewEvent struct {
	Event        string  `json:"event"`
	Id           int     `json:"id"`
	ReviewerId   int     `json:"reviewer_id"`
	RestaurantId string  `json:"restaurant_id"`
	Rating       float32 `json:"rating"`
	Content      string  `json:"content"`
}

type FavoriteEvent struct {
	Event        string `json:"event"`
	UserId       int    `json:"user_id"`
	FoodId       string `json:"food_id"`
	RestaurantId string `json:"restaurant_id"`
}
