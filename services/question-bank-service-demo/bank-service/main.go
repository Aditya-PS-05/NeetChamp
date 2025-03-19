package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/database"
	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/handlers"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request: Method=%s, Path=%s, Resource=%s",
		request.HTTPMethod, request.Path, request.Resource)

	// Log path parameters for debugging
	log.Printf("Path parameters: %+v", request.PathParameters)
	log.Printf("Query string parameters: %+v", request.QueryStringParameters)

	// Set default HTTP method if empty
	httpMethod := request.HTTPMethod
	if httpMethod == "" {
		if request.Body != "" {
			httpMethod = "POST"
		} else {
			httpMethod = "GET"
		}
	}

	// Extract path for easier processing
	path := request.Path

	// Direct routing based on HTTP method and path pattern
	switch {
	// GET /api/question/{questionId}
	case httpMethod == "GET" && strings.HasPrefix(path, "/api/question/"):
		questionId := extractQuestionId(path)
		if questionId != "" {
			log.Printf("GET question with ID: %s", questionId)
			return handlers.GetQuestion(request)
		}

	// POST /api/question/add
	case httpMethod == "POST" && (strings.HasSuffix(path, "/api/question/add") ||
		strings.HasSuffix(path, "/question/add")):
		return handlers.AddQuestion(request)

	// Alternative POST handling for just /api/question
	case httpMethod == "POST" && (path == "/api/question" || path == "/question"):
		return handlers.AddQuestion(request)

	// POST with quiz_id in body (fallback)
	case httpMethod == "POST" && strings.Contains(request.Body, "quiz_id"):
		return handlers.AddQuestion(request)

	// PUT /api/question/{questionId}/update
	case httpMethod == "PUT" && strings.Contains(path, "/update"):
		return handlers.UpdateQuestion(request)

	// DELETE /api/question/{questionId}/delete
	case httpMethod == "DELETE" && strings.Contains(path, "/delete"):
		return handlers.DeleteQuestion(request)
	}

	// Check for question ID in path parameters
	if questionId, ok := request.PathParameters["questionId"]; ok && httpMethod == "GET" {
		log.Printf("Found questionId in path parameters: %s", questionId)
		// Create a new request with modified path for consistent processing
		modifiedRequest := request
		modifiedRequest.Path = "/api/question/" + questionId
		return handlers.GetQuestion(modifiedRequest)
	}

	// Default response for unmatched routes
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       fmt.Sprintf(`{"error": "Route not found", "endpoint": "%s", "method": "%s", "path": "%s"}`, request.Path, httpMethod, path),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// Extract question ID from the path
func extractQuestionId(path string) string {
	// Match patterns like /api/question/{id} or /api/question/{id}/update
	re := regexp.MustCompile(`/api/question/([^/]+)(?:/|$)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) > 1 {
		id := matches[1]
		// Don't return "add", "update", or "delete" as IDs
		if id != "add" && id != "update" && id != "delete" {
			return id
		}
	}

	// Try alternative pattern without /api prefix
	re = regexp.MustCompile(`/question/([^/]+)(?:/|$)`)
	matches = re.FindStringSubmatch(path)

	if len(matches) > 1 {
		id := matches[1]
		if id != "add" && id != "update" && id != "delete" {
			return id
		}
	}

	return ""
}

func main() {
	database.InitDynamoDB()
	fmt.Println("🚀 NeetChamp Question Bank Service Started!")
	lambda.Start(handler)
}
