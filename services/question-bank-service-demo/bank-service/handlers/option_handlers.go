package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/database"
	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/models"
	"github.com/aws/aws-lambda-go/events"
)

// AddOptionRequest represents the request body for adding an option
type AddOptionRequest struct {
	QuestionID string `json:"question_id"`
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// AddOption handles adding a new option
func AddOption(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing AddOption request")

	// Parse request body
	var req AddOptionRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "Invalid request format: %s"}`, err.Error()),
		}, nil
	}

	// Validate input
	if req.QuestionID == "" || req.OptionText == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "question_id and option_text are required"}`,
		}, nil
	}

	// Ensure the question exists
	_, err = database.GetQuestionByID(req.QuestionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Question not found"}`,
		}, nil
	}

	// Create and save the option
	option := models.NewOption(req.QuestionID, req.OptionText, req.IsCorrect)
	err = database.SaveOption(option)
	if err != nil {
		log.Printf("Error saving option: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to save option: %s"}`, err.Error()),
		}, nil
	}

	// Return the created option
	optionJSON, _ := json.Marshal(option)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(optionJSON),
	}, nil
}

// // DeleteOption handles deleting an option
// func DeleteOption(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	log.Println("Processing DeleteOption request")

// 	optionID := request.PathParameters["optionId"]
// 	if optionID == "" {
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       `{"error": "Option ID is required"}`,
// 		}, nil
// 	}

// 	err := database.DeleteOption(optionID)
// 	if err != nil {
// 		log.Printf("Error deleting option: %v", err)
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       fmt.Sprintf(`{"error": "Failed to delete option: %s"}`, err.Error()),
// 		}, nil
// 	}

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusOK,
// 		Body:       `{"message": "Option deleted successfully"}`,
// 	}, nil
// }

// GetOptionsByQuestion handles fetching all options for a question
func GetOptionsByQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing GetOptionsByQuestion request")

	questionID := request.PathParameters["questionId"]
	if questionID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Question ID is required"}`,
		}, nil
	}

	options, err := database.GetOptionsByQuestionID(questionID)
	if err != nil {
		log.Printf("Error fetching options: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to fetch options: %s"}`, err.Error()),
		}, nil
	}

	optionsJSON, _ := json.Marshal(options)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(optionsJSON),
	}, nil
}
