// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "taxsi2 is a WAF server application The base path for all the APIs is \"/api/v1\".\n",
    "title": "taxsi2",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/health": {
      "get": {
        "description": "Check if taxsi2 is healthy",
        "tags": [
          "health"
        ],
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "status of health check",
            "schema": {
              "$ref": "#/definitions/health"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string",
          "minLength": 1
        }
      }
    },
    "health": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string"
        }
      }
    }
  },
  "tags": [
    {
      "description": "Check if taxsi2 is healthy",
      "name": "health"
    }
  ],
  "x-tagGroups": [
    {
      "name": "taxsi2 Management",
      "tags": [
        "app"
      ]
    },
    {
      "name": "Health Check",
      "tags": [
        "health"
      ]
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "taxsi2 is a WAF server application The base path for all the APIs is \"/api/v1\".\n",
    "title": "taxsi2",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/health": {
      "get": {
        "description": "Check if taxsi2 is healthy",
        "tags": [
          "health"
        ],
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "status of health check",
            "schema": {
              "$ref": "#/definitions/health"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string",
          "minLength": 1
        }
      }
    },
    "health": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string"
        }
      }
    }
  },
  "tags": [
    {
      "description": "Check if taxsi2 is healthy",
      "name": "health"
    }
  ],
  "x-tagGroups": [
    {
      "name": "taxsi2 Management",
      "tags": [
        "app"
      ]
    },
    {
      "name": "Health Check",
      "tags": [
        "health"
      ]
    }
  ]
}`))
}
