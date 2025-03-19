package models

import (
	"time"

	"github.com/google/uuid"
)

// Question types as constants
const (
	QuestionTypeMCQ         = "MCQ"
	QuestionTypeTrueFalse   = "True/False"
	QuestionTypeFillBlank   = "Fill in the Blank"
	QuestionTypeShortAnswer = "Short Answer"
)

// Question represents a question in the question bank
type Question struct {
	QuestionID   string    `json:"question_id" dynamodbav:"question_id"`
	QuizID       string    `json:"quiz_id" dynamodbav:"quiz_id"`
	QuestionText string    `json:"question_text" dynamodbav:"question_text"`
	QuestionType string    `json:"question_type" dynamodbav:"question_type"`
	Answer       string    `json:"answer" dynamodbav:"answer"`
	CreatedAt    time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// Option represents an answer choice for a question
type Option struct {
	OptionID   string    `json:"option_id" dynamodbav:"option_id"`
	QuestionID string    `json:"question_id" dynamodbav:"question_id"`
	OptionText string    `json:"option_text" dynamodbav:"option_text"`
	IsCorrect  bool      `json:"is_correct" dynamodbav:"is_correct"`
	CreatedAt  time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// NewQuestion creates a new question with default values
func NewQuestion(quizID, questionText, questionType, answer string) Question {
	now := time.Now()
	return Question{
		QuestionID:   uuid.New().String(),
		QuizID:       quizID,
		QuestionText: questionText,
		QuestionType: questionType,
		Answer:       answer,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewOption creates a new option for a question
func NewOption(questionID, optionText string, isCorrect bool) Option {
	now := time.Now()
	return Option{
		OptionID:   uuid.New().String(),
		QuestionID: questionID,
		OptionText: optionText,
		IsCorrect:  isCorrect,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// ValidateQuestionType checks if the question type is valid
func ValidateQuestionType(questionType string) bool {
	validTypes := []string{
		QuestionTypeMCQ,
		QuestionTypeTrueFalse,
		QuestionTypeFillBlank,
		QuestionTypeShortAnswer,
	}

	for _, validType := range validTypes {
		if questionType == validType {
			return true
		}
	}
	return false
}

// QuestionWithOptions represents a question with its options
type QuestionWithOptions struct {
	Question Question `json:"question"`
	Options  []Option `json:"options,omitempty"`
}
