package model

import (
	"gorm.io/gorm"
)

type Review struct {
  gorm.Model
  RestaurantId string
  UserId string
  Rating float32
  Content string
}