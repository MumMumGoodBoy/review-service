package main

import (
	"log"
	"net"

	"github.com/MumMumGoodBoy/review-service/internal/model"
	"github.com/MumMumGoodBoy/review-service/internal/service"
	"github.com/MumMumGoodBoy/review-service/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
  db, err := gorm.Open(postgres.Open("host=localhost user=user password=pass dbname=review port=5432"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  db.AutoMigrate(&model.Review{})
  
	reviewService := service.ReviewService{
		DB: db,
  }

	grpcServer := grpc.NewServer()
	proto.RegisterReviewServer(grpcServer, &reviewService)

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}