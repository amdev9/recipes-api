// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "twitter @amdev9",
            "url": "https://x.com/amdev99"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/recipes": {
            "get": {
                "description": "get recipes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "List recipes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/main.Recipe"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "add by json recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Add recipe",
                "parameters": [
                    {
                        "description": "Add recipe",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Recipe"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Recipe"
                        }
                    }
                }
            }
        },
        "/recipes/search": {
            "get": {
                "responses": {}
            }
        },
        "/recipes/{id}": {
            "put": {
                "description": "Update by json recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Update recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update recipe",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Recipe"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Recipe"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete by recipe ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Delete recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Recipe"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Recipe": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "ingredients": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "instructions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "publishedAt": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger Recipes with Mongo API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
