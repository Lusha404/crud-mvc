package main

import (
	"context"
	"crud-gin/controllers"
	"crud-gin/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var (
	server         *gin.Engine
	userService    services.UserService
	userController controllers.UserController
	ctx            context.Context
	userCollection *mongo.Collection
	mongoClient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()
	mongoConn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatalf("failed to establish a database connection: %v\n", err)
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("mongo connection established")

	userCollection = mongoClient.Database("userdb").Collection("users")
	userService = services.NewUserService(ctx, userCollection)
	userController = controllers.New(userService)
	server = gin.Default()
}

// .../v1/user/create
func main() {
	defer mongoClient.Disconnect(ctx)
	basePath := server.Group("/v1")
	userController.RegisterUserRoutes(basePath)
	log.Fatal(server.Run(":9000"))
}
