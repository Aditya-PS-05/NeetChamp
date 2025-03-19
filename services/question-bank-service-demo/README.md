```
MeetChamp
└── question-bank-service-demo/
    ├── go.mod
    ├── services/
    │   └── auth-service/
    │       ├── database/ -> Dynamo,go
    │       ├── handlers/ -> all the handlers
    │       ├── models/ -> QuestionBank.go
    │       ├── utils/ -> all the helper modules
    │       └── main.go -> main.go
    └── stacks/
        └── QuestionBankStack.ts
    └── main
    └── sst.config.ts

```

Question Bank Service - System Design

Overview

The Question Bank Service is a microservice responsible for storing, retrieving, updating, and deleting quiz questions. It is designed to be scalable, reliable, and efficient using AWS services like DynamoDB, gRPC, and AWS Lambda where applicable.

Components

1. API Gateway

Exposes HTTP/REST endpoints for admin and internal services.

Routes requests to the appropriate gRPC methods using AWS API Gateway.

2. gRPC Gateway

Acts as a bridge for internal communication.

Provides low-latency communication for the microservices.

3. Auth Service

Ensures that only admins can perform CRUD operations.

Uses JWT for authentication and role-based authorization.

4. Question Bank Service

Core microservice that handles:

Creating questions

Updating questions

Deleting questions

Fetching questions based on filters (e.g., subject, difficulty)

Built with Golang and deployed using SST.

5. DynamoDB

NoSQL database to store questions in a flexible schema.

Efficient read/write operations with auto-scaling.

6. Monitoring & Logging

AWS CloudWatch for monitoring logs and system metrics.

AWS X-Ray for distributed tracing.

Data Model

Table Name: Questions

Field

Data Type

Description

id

String

Unique identifier for the question

question

String

The actual quiz question

options

List

Multiple-choice options

answer

String

Correct answer for the question

subject

String

Subject category of the question

difficulty

String

Difficulty level (easy, medium, hard)

created_at

String

Timestamp of creation

API Endpoints

Endpoint

Method

Description

Auth Required

/question

POST

Create a new question

Yes (Admin)

/question/{id}

GET

Get question by ID

Yes (Admin)

/question

PUT

Update a question

Yes (Admin)

/question/{id}

DELETE

Delete a question by ID

Yes (Admin)

/questions?subject=&difficulty=

GET

Get questions by filter

No (Public)

Error Handling

Detailed error codes with meaningful error messages.

Logs all errors to CloudWatch.

Graceful fallback in case of DynamoDB read/write issues.

Scaling Strategy

DynamoDB Autoscaling for dynamic workloads.

gRPC Load Balancer to handle traffic spikes.

AWS Lambda for background processing (e.g., bulk question uploads).

Security Measures

AWS IAM for managing permissions.

JWT authentication with RBAC (Role-Based Access Control).

Data encryption using AWS KMS.

API rate limiting to prevent DDoS attacks.

This design ensures the Question Bank Service is scalable, secure, and highly available. Let me know if you'd like further adjustments or additional diagrams for this architecture.

aditya@aditya:~/my-work/Freelance/NeetChamp/services/question-bank-service-demo$ GOOS=linux GOARCH=arm64 go build -o main bank-service/main.go