package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/MumMumGoodBoy/review-service/internal/model"
	"github.com/MumMumGoodBoy/review-service/internal/service"
	"github.com/MumMumGoodBoy/review-service/proto"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open("host=localhost user=user password=pass dbname=review port=5432"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Review{})
	db.AutoMigrate(&model.FavoriteFood{})

	// Connect to RabbitMQ
	rabbitMQConn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQConn.Close()
	rabbitMQChannel, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQChannel.Close()
	fmt.Println("Connected to RabbitMQ")

	// service
	reviewService := service.ReviewService{
		DB:              db,
		RabbitMQChannel: rabbitMQChannel,
	}

	// grpc
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
