// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/current": {
            "get": {
                "description": "current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "CurrentUser",
                "operationId": "current-user",
                "responses": {}
            }
        },
        "/auth/logout": {
            "get": {
                "description": "logout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout",
                "operationId": "logout",
                "responses": {}
            }
        },
        "/auth/refresh": {
            "get": {
                "description": "refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "RefreshToken",
                "operationId": "refresh-token",
                "responses": {}
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "sign in",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignIn",
                "operationId": "sign-in",
                "parameters": [
                    {
                        "description": "body info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UserLogin"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "sign un",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignUp",
                "operationId": "sign-up",
                "parameters": [
                    {
                        "description": "body info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.User"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "request.User": {
            "type": "object",
            "required": [
                "email",
                "nickname",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "logg": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 3
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        },
        "request.UserLogin": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 6
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
