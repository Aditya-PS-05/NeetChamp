package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/database"
	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/models"
	"github.com/aws/aws-lambda-go/events"
)

// QuestionRequest represents the request body for adding/updating a question
type QuestionRequest struct {
	QuizID       string        `json:"quiz_id"`
	QuestionText string        `json:"question_text"`
	QuestionType string        `json:"question_type"`
	Options      []OptionInput `json:"options"`
	Answer       string        `json:"answer"`
}

// OptionInput represents the input for an option
type OptionInput struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// AddQuestion handles adding a new question with options
func AddQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing AddQuestion request")

	// Parse request body
	var req QuestionRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "Invalid request format: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Validate input
	if req.QuizID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "quiz_id is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	if req.QuestionText == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "question_text is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	if !models.ValidateQuestionType(req.QuestionType) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid question_type. Must be one of: MCQ, True/False, Fill in the Blank, Short Answer"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Validate options based on question type
	if req.QuestionType == models.QuestionTypeMCQ {
		if len(req.Options) < 2 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "MCQ questions must have at least 2 options"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		hasCorrectOption := false
		for _, opt := range req.Options {
			if opt.IsCorrect {
				hasCorrectOption = true
				break
			}
		}

		if !hasCorrectOption {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "At least one option must be marked as correct for MCQ questions"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
	} else if req.QuestionType == models.QuestionTypeTrueFalse {
		if len(req.Options) != 2 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "True/False questions must have exactly 2 options"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		// Check if only one option is correct
		correctCount := 0
		for _, opt := range req.Options {
			if opt.IsCorrect {
				correctCount++
			}
		}

		if correctCount != 1 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "True/False questions must have exactly one correct option"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
	} else if req.QuestionType == models.QuestionTypeFillBlank || req.QuestionType == models.QuestionTypeShortAnswer {
		// For these types, we expect an answer instead of options
		if req.Answer == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "Answer is required for Fill in the Blank and Short Answer questions"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
	}

	// Create and save the question
	question := models.NewQuestion(req.QuizID, req.QuestionText, req.QuestionType, req.Answer)

	// Add answer field if it's not MCQ/True-False
	if req.QuestionType == models.QuestionTypeFillBlank || req.QuestionType == models.QuestionTypeShortAnswer {
		question.Answer = req.Answer
	}

	err = database.SaveQuestion(question)
	if err != nil {
		log.Printf("Error saving question: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to save question: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Save options if applicable
	if req.QuestionType == models.QuestionTypeMCQ || req.QuestionType == models.QuestionTypeTrueFalse {
		for _, optInput := range req.Options {
			option := models.NewOption(question.QuestionID, optInput.OptionText, optInput.IsCorrect)
			err = database.SaveOption(option)
			if err != nil {
				log.Printf("Error saving option: %v", err)
				// Continue saving other options
			}
		}
	}

	// Fetch the newly created question with options
	result, err := getQuestionWithOptions(question.QuestionID)
	if err != nil {
		log.Printf("Error fetching created question: %v", err)
		// Return just the question ID if there's an error
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusCreated,
			Body:       fmt.Sprintf(`{"question_id": "%s", "message": "Question created successfully"}`, question.QuestionID),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Return the created question with options
	resultJSON, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(resultJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// GetQuestion handles fetching a question by ID with its options
func GetQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing GetQuestion request")

	// Get question ID from path parameters
	questionID := request.PathParameters["questionId"]
	if questionID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Question ID is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Fetch the question with options
	result, err := getQuestionWithOptions(questionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error": "Question not found"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to fetch question: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Return the question with options
	resultJSON, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resultJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// UpdateQuestion handles updating a question with its options
func UpdateQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing UpdateQuestion request")

	// Get question ID from path parameters
	questionID := request.PathParameters["questionId"]
	if questionID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Question ID is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Parse request body
	var req QuestionRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("Error unmarshaling request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"error": "Invalid request format: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Fetch the existing question
	question, err := database.GetQuestionByID(questionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error": "Question not found"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to fetch question: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Update question fields if provided
	updated := false

	if req.QuestionText != "" {
		question.QuestionText = req.QuestionText
		updated = true
	}

	if req.QuestionType != "" {
		if !models.ValidateQuestionType(req.QuestionType) {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "Invalid question_type. Must be one of: MCQ, True/False, Fill in the Blank, Short Answer"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		// Validate options based on new question type
		if req.QuestionType == models.QuestionTypeMCQ {
			if len(req.Options) < 2 {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       `{"error": "MCQ questions must have at least 2 options"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				}, nil
			}

			hasCorrectOption := false
			for _, opt := range req.Options {
				if opt.IsCorrect {
					hasCorrectOption = true
					break
				}
			}

			if !hasCorrectOption {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       `{"error": "At least one option must be marked as correct for MCQ questions"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				}, nil
			}
		} else if req.QuestionType == models.QuestionTypeTrueFalse {
			if len(req.Options) != 2 {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       `{"error": "True/False questions must have exactly 2 options"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				}, nil
			}

			// Check if only one option is correct
			correctCount := 0
			for _, opt := range req.Options {
				if opt.IsCorrect {
					correctCount++
				}
			}

			if correctCount != 1 {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       `{"error": "True/False questions must have exactly one correct option"}`,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				}, nil
			}
		} else if (req.QuestionType == models.QuestionTypeFillBlank || req.QuestionType == models.QuestionTypeShortAnswer) && req.Answer == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error": "Answer is required for Fill in the Blank and Short Answer questions"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		question.QuestionType = req.QuestionType
		updated = true
	}

	// Update answer field if applicable
	if req.Answer != "" && (question.QuestionType == models.QuestionTypeFillBlank || question.QuestionType == models.QuestionTypeShortAnswer) {
		question.Answer = req.Answer
		updated = true
	}

	// Only save if something changed
	if updated {
		// Update timestamp
		question.UpdatedAt = time.Now()

		// Save the updated question
		err = database.SaveQuestion(*question)
		if err != nil {
			log.Printf("Error saving question: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf(`{"error": "Failed to update question: %s"}`, err.Error()),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
	}

	// Handle options update - only if MCQ or True/False
	if (question.QuestionType == models.QuestionTypeMCQ || question.QuestionType == models.QuestionTypeTrueFalse) && len(req.Options) > 0 {
		// Delete existing options
		err = database.DeleteOptionsByQuestionID(questionID)
		if err != nil {
			log.Printf("Error deleting existing options: %v", err)
			// Continue anyway as this might be a transient error
		}

		// Save new options
		for _, optInput := range req.Options {
			option := models.NewOption(questionID, optInput.OptionText, optInput.IsCorrect)
			err = database.SaveOption(option)
			if err != nil {
				log.Printf("Error saving option: %v", err)
				// Continue saving other options
			}
		}
	}

	// Fetch the updated question with options
	result, err := getQuestionWithOptions(questionID)
	if err != nil {
		log.Printf("Error fetching updated question: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"message": "Question updated successfully, but error fetching updated data"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Return the updated question with options
	resultJSON, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resultJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// DeleteQuestion handles deleting a question and its options
func DeleteQuestion(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing DeleteQuestion request")

	// Get question ID from path parameters
	questionID := request.PathParameters["questionId"]
	if questionID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Question ID is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Check if the question exists
	_, err := database.GetQuestionByID(questionID)
	if err != nil {
		log.Printf("Error fetching question: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error": "Question not found"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to fetch question: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Delete associated options first
	err = database.DeleteOptionsByQuestionID(questionID)
	if err != nil {
		log.Printf("Error deleting options: %v", err)
		// Continue with question deletion anyway
	}

	// Delete the question
	err = database.DeleteQuestion(questionID)
	if err != nil {
		log.Printf("Error deleting question: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to delete question: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Return success response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "Question and associated options deleted successfully"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// GetQuestionsByQuiz handles fetching all questions for a quiz
func GetQuestionsByQuiz(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Processing GetQuestionsByQuiz request")

	// Get quiz ID from path parameters
	quizID := request.PathParameters["quizId"]
	if quizID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Quiz ID is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Fetch all questions for the quiz
	questions, err := database.GetQuestionsByQuizID(quizID)
	if err != nil {
		log.Printf("Error fetching questions: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Failed to fetch questions: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// For each question, fetch its options
	var questionsWithOptions []models.QuestionWithOptions
	for _, q := range questions {
		options, err := database.GetOptionsByQuestionID(q.QuestionID)
		if err != nil {
			log.Printf("Error fetching options for question %s: %v", q.QuestionID, err)
			// Continue with next question
			continue
		}

		questionsWithOptions = append(questionsWithOptions, models.QuestionWithOptions{
			Question: q,
			Options:  options,
		})
	}

	// Return the questions with options
	resultJSON, _ := json.Marshal(questionsWithOptions)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(resultJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// Helper function to get a question with its options
func getQuestionWithOptions(questionID string) (*models.QuestionWithOptions, error) {
	// Fetch the question
	question, err := database.GetQuestionByID(questionID)
	if err != nil {
		return nil, err
	}

	// Fetch options if applicable
	var options []models.Option
	if question.QuestionType == models.QuestionTypeMCQ || question.QuestionType == models.QuestionTypeTrueFalse {
		options, err = database.GetOptionsByQuestionID(questionID)
		if err != nil {
			log.Printf("Error fetching options: %v", err)
			// Continue even if options can't be fetched
		}
	}

	// Return the question with options
	return &models.QuestionWithOptions{
		Question: *question,
		Options:  options,
	}, nil
}
