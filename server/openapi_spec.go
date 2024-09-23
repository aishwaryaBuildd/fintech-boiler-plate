package server

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

var openapiYAML = `
swagger: "2.0"
info:
  title: "Gin API"
  description: "This is a sample server for a Gin application."
  version: "1.0"
host: "localhost:8080"
basePath: "/api/v1"
paths:
  /register:
    post:
      summary: Register a new user
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: user
          required: true
          schema:
            type: object
            properties:
              email:
                type: string
              password:
                type: string
      responses:
        200:
          description: User registered successfully
        400:
          description: Invalid request body
        500:
          description: Internal server error
  /login:
    post:
      summary: Login with email and password
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: user
          required: true
          schema:
            type: object
            properties:
              email:
                type: string
              password:
                type: string
      responses:
        200:
          description: Login successful
        400:
          description: Invalid request body
        401:
          description: Unauthorized
        500:
          description: Internal server error
`

func LoadSwaggerSpec() (*loads.Document, error) {
	var specObj spec.Swagger
	err := yaml.Unmarshal([]byte(openapiYAML), &specObj)
	if err != nil {
		log.Fatalf("Error unmarshalling OpenAPI YAML: %v", err)
		return nil, err
	}

	specJSON, err := json.Marshal(specObj)
	if err != nil {
		log.Fatalf("Error marshalling OpenAPI JSON: %v", err)
		return nil, err
	}

	return loads.Analyzed(json.RawMessage(specJSON), "")
}
