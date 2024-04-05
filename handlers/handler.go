package handlers

import (
	"context"
	"fmt"
	"net/http"
	"recipes-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.
	Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// NewRecipeHandler godoc
//
//	@Summary		Add recipe
//	@Description	add by json recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe	body		models.Recipe		true	"Add recipe"
//	@Success		200		{object}	models.Recipe
//	@Router			/recipes [post]
func (handler *RecipesHandler) CreateRecipeHandler(c *gin.Context) {

	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := handler.collection.InsertOne(handler.ctx, recipe)

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
//	@Success		200	{array}		[]models.Recipe
//	@Router			/recipes [get]
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)
	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
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
//	@Param			recipe	body		models.Recipe				true	"Update recipe"
//	@Success		200		{object}	models.Recipe
//	@Router			/recipes/{id} [put]
func (handler *RecipesHandler) UpdateRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
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

	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
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
//	@Success		200	{object}	models.Recipe
//	@Failure		404	{object}	string
//
// @Router			/recipes/{id} [delete]
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectId})

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
//	@Success		200	{array}		[]models.Recipe

// @Router			/recipes/search [get]
func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cur, err := handler.collection.Find(handler.ctx, bson.M{
		"tags": tag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)
	listOfRecipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		listOfRecipes = append(listOfRecipes, recipe)
	}

	c.JSON(http.StatusOK, listOfRecipes)
}
