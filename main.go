package main

import (
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
	postgresUrl := os.Getenv("POSTGRES_URL")
	port := os.Getenv("PORT")

	db, err := gorm.Open(postgres.Open(postgresUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("Connected to database")
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

	log.Println("Connected to RabbitMQ")
	// service
	reviewService := service.ReviewService{
		DB:              db,
		RabbitMQChannel: rabbitMQChannel,
	}

	// grpc
	grpcServer := grpc.NewServer()
	proto.RegisterReviewServer(grpcServer, &reviewService)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Review service is running on port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
