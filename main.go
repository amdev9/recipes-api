package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	_ "recipes-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json: "ingredients"`
	Instructions []string  `json: "instructions"`
	PublishedAt  time.Time `json: "publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

// NewRecipeHandler godoc
//
//	@Summary		Add recipe
//	@Description	add by json recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			account	body		Recipe		true	"Add recipe"
//	@Success		200		{object}	Recipe
//	@Router			/recipes [post]
func NewRecipeHandler(c *gin.Context) {

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
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
//	@Param			account	body		Recipe				true	"Update recipe"
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
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			recipe.ID = id
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
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
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
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
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

// @title           Swagger Recipes API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   twitter @amdev9
// @contact.url    https://x.com/amdev99

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipesHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
