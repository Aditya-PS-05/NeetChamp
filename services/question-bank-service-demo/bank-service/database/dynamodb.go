package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/Aditya-PS-05/NeetChamp-question-bank-service/bank-service/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var db *dynamodb.DynamoDB

func InitDynamoDB() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	db = dynamodb.New(sess)
	fmt.Println("DynamoDB initialized")
}

// Question-related functions

// SaveQuestion saves a question to the database
func SaveQuestion(question models.Question) error {
	tableName := os.Getenv("QUESTIONS_TABLE")
	if tableName == "" {
		tableName = "QuestionsTable"
	}

	if question.QuestionID == "" {
		question.QuestionID = uuid.New().String()
	}

	now := time.Now()
	if question.CreatedAt.IsZero() {
		question.CreatedAt = now
	}
	question.UpdatedAt = now

	av, err := dynamodbattribute.MarshalMap(question)
	if err != nil {
		log.Printf("Error marshaling question: %v", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		log.Printf("Error putting item in DynamoDB: %v", err)
		return err
	}

	return nil
}

// GetQuestionByID retrieves a question by its ID
func GetQuestionByID(questionID string) (*models.Question, error) {
	tableName := os.Getenv("QUESTIONS_TABLE")
	if tableName == "" {
		tableName = "QuestionsTable"
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"question_id": {S: aws.String(questionID)},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("question not found")
	}

	question := &models.Question{}
	err = dynamodbattribute.UnmarshalMap(result.Item, question)
	if err != nil {
		return nil, err
	}

	return question, nil
}

// DeleteQuestion deletes a question by its ID
func DeleteQuestion(questionID string) error {
	tableName := os.Getenv("QUESTIONS_TABLE")
	if tableName == "" {
		tableName = "QuestionsTable"
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"question_id": {S: aws.String(questionID)},
		},
	}

	_, err := db.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

// GetQuestionsByQuizID retrieves all questions for a given quiz
func GetQuestionsByQuizID(quizID string) ([]models.Question, error) {
	tableName := os.Getenv("QUESTIONS_TABLE")
	if tableName == "" {
		tableName = "QuestionsTable"
	}

	// Create the expression for the filter
	filt := expression.Name("quiz_id").Equal(expression.Value(quizID))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}

	questions := []models.Question{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// Option-related functions

// SaveOption saves an option to the database
func SaveOption(option models.Option) error {
	tableName := os.Getenv("OPTIONS_TABLE")
	if tableName == "" {
		tableName = "OptionsTable"
	}

	if option.OptionID == "" {
		option.OptionID = uuid.New().String()
	}

	av, err := dynamodbattribute.MarshalMap(option)
	if err != nil {
		log.Printf("Error marshaling option: %v", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = db.PutItem(input)
	if err != nil {
		log.Printf("Error putting item in DynamoDB: %v", err)
		return err
	}

	return nil
}

// GetOptionsByQuestionID retrieves all options for a given question
func GetOptionsByQuestionID(questionID string) ([]models.Option, error) {
	tableName := os.Getenv("OPTIONS_TABLE")
	if tableName == "" {
		tableName = "OptionsTable"
	}

	// Create the expression for the filter
	filt := expression.Name("question_id").Equal(expression.Value(questionID))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}

	options := []models.Option{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}

// DeleteOptionsByQuestionID deletes all options for a given question
func DeleteOptionsByQuestionID(questionID string) error {
	options, err := GetOptionsByQuestionID(questionID)
	if err != nil {
		return err
	}

	tableName := os.Getenv("OPTIONS_TABLE")
	if tableName == "" {
		tableName = "OptionsTable"
	}

	for _, option := range options {
		input := &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"option_id": {S: aws.String(option.OptionID)},
			},
		}

		_, err := db.DeleteItem(input)
		if err != nil {
			return err
		}
	}

	return nil
}
