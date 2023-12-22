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
        "/auth/admin": {
            "get": {
                "description": "Generate an admin token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "AUTH"
                ],
                "summary": "Generate admin token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.TokenResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    }
                }
            }
        },
        "/auth/user": {
            "get": {
                "description": "Generate an user token.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "AUTH"
                ],
                "summary": "Generate user token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.TokenResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    }
                }
            }
        },
        "/tasks/multiple": {
            "post": {
                "description": "Create, update and delete tasks with single request.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TASK"
                ],
                "summary": "Process Multiple Tasks",
                "parameters": [
                    {
                        "description": "Input",
                        "name": "processMultipleTasks",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.PayloadCollection"
                        }
                    },
                    {
                        "type": "string",
                        "default": "Bearer eyJblabla",
                        "description": "Insert your access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.ResponseOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.OperationType": {
            "type": "string",
            "enum": [
                "OpPut",
                "OpCreate",
                "OpDelete"
            ],
            "x-enum-varnames": [
                "OpPut",
                "OpCreate",
                "OpDelete"
            ]
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.Payload": {
            "type": "object",
            "required": [
                "data",
                "operationType"
            ],
            "properties": {
                "data": {},
                "operationType": {
                    "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.OperationType"
                }
            }
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.PayloadCollection": {
            "type": "object",
            "required": [
                "payloads"
            ],
            "properties": {
                "payloads": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.Payload"
                    }
                }
            }
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.ResponseOK": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "string"
                }
            }
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_dto.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIErrors": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_MehmetTalhaSeker_concurrent-web-service_internal_utils_errorutils.APIError"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.01",
	Host:             "localhost:3333",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API",
	Description:      "Concurrent Data Processing Web Service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
