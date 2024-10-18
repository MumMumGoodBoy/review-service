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

type FavouriteFood struct {
	gorm.Model
	UserId uint
	FoodId string
}

type ReviewEvent struct {
	Event        string  `json:"event"`
	Id           int     `json:"id"`
	ReviewerId   int     `json:"reviewer_id"`
	RestaurantId string  `json:"restaurant_id"`
	Rating       float32 `json:"rating"`
	Content      string  `json:"content"`
}
