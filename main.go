package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	_ "recipes-api/docs"
	"recipes-api/handlers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-redis/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler

func init() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)

	ctx := context.Background()
	client, errMongo := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if errMongo != nil {
		log.Fatal("Error connection to MongoDB ")
		return
	} else {
		log.Println("Connected to MongoDB")

		users := map[string]string{
			"admin":      "fCRmh4Q2J7Rseqkz",
			"packt":      "RE4zfHB35VPtTkbT",
			"mlabouardy": "L3nSFRcZzNQ67bcc",
		}
		collection := client.Database(os.Getenv(
			"MONGO_DATABASE")).Collection("users")
		h := sha256.New()
		for username, password := range users {
			collection.InsertOne(ctx, bson.M{
				"username": username,
				"password": string(h.Sum([]byte(password))),
			})
		}
	}

	if err := client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func cleanup() {
	fmt.Println("cleanup")
}

// @title           Swagger Recipes with Mongo API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   twitter @amdev9
// @contact.url    https://x.com/amdev99

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.CreateRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipesHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	authorized.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.Run(":3000")
}
