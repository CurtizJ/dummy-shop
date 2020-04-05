// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-03-29 01:37:12.674757435 +0300 MSK m=+0.034152139

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/item": {
            "put": {
                "description": "Update single item",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "Update item",
                "parameters": [
                    {
                        "description": "New item to update. Item with updated id should be already added.",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/items.Item"
                        }
                    }
                ],
                "responses": {
                    "200": {}
                }
            },
            "post": {
                "description": "Add single item",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "Add item",
                "parameters": [
                    {
                        "description": "Item to add",
                        "name": "item",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/items.Item"
                        }
                    }
                ],
                "responses": {
                    "200": {}
                }
            }
        },
        "/item/{id}": {
            "get": {
                "description": "Get single item by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "Get item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/items.Item"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete single item by id",
                "tags": [
                    "items"
                ],
                "summary": "Delete item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Item id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {}
                }
            }
        },
        "/items": {
            "get": {
                "description": "List all items with optional pagination",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "List items",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Return at most 'length' items.",
                        "name": "length",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Skip first 'offset' items. Must be specified with 'length'",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/items.Item"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "items.Item": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8181",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Shop API",
	Description: "Educational online-shop API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}