package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	_ "recipes-api/docs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

var recipes []Recipe

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	log.Println("Connected to MongoDB")
}

// NewRecipeHandler godoc
//
//	@Summary		Add recipe
//	@Description	add by json recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe	body		Recipe		true	"Add recipe"
//	@Success		200		{object}	Recipe
//	@Router			/recipes [post]
func NewRecipeHandler(c *gin.Context) {

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := collection.InsertOne(ctx, recipe)

	if err != nil {

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler godoc
//
//	@Summary		List recipes
//	@Description	get recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		[]Recipe
//	@Router			/recipes [get]
func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	recipes = make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipesHandler godoc
//
//	@Summary		Update recipe
//	@Description	Update by json recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Recipe ID"
//	@Param			recipe	body		Recipe				true	"Update recipe"
//	@Success		200		{object}	Recipe
//	@Router			/recipes/{id} [put]
func UpdateRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)

	bsonD := bson.D{
		{Key: "name", Value: recipe.Name},
		{Key: "instructions", Value: recipe.Instructions},
		{Key: "ingredients", Value: recipe.Ingredients},
		{Key: "tags", Value: recipe.Tags},
	}

	_, err := collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{{Key: "$set", Value: bsonD}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// DeleteRecipeHandler godoc
//
//	@Summary		Delete recipe
//	@Description	Delete by recipe ID
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Recipe ID"	string
//	@Success		200	{object}	Recipe
//	@Failure		404	{object}	string
//
// @Router			/recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

// SearchRecipesHandler godoc
//
//	@Summary		Search recipes by tags
//	@Description	get recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		[]Recipe

// @Router			/recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cur, err := collection.Find(ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	listOfRecipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		listOfRecipes = append(listOfRecipes, recipe)
	}

	c.JSON(http.StatusOK, listOfRecipes)
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
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipesHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":3000")
}
