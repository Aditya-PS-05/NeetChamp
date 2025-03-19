import { StackContext, Function, Api, Table } from "sst/constructs";

export function QuestionBankStack({ stack }: StackContext) {
  // Questions Table
  const questionsTable = new Table(stack, "QuestionsTable", {
    fields: {
      question_id: "string",
      quiz_id: "string",
      question_text: "string",
      question_type: "string",
      created_at: "string",
      updated_at: "string",
    },
    primaryIndex: { partitionKey: "question_id" },
    globalIndexes: {
      quizIndex: { partitionKey: "quiz_id" },
    },
  });

  // Options Table
  const optionsTable = new Table(stack, "OptionsTable", {
    fields: {
      option_id: "string",
      question_id: "string",
      option_text: "string",
      is_correct: "boolean" as any,
    },
    primaryIndex: { partitionKey: "option_id" },
    globalIndexes: {
      questionIndex: { partitionKey: "question_id" },
    },
  });

  const questionBankFunction = new Function(stack, "QuestionBankFunction", {
    handler: "bank-service/main.go",
    runtime: "go",
    architecture: "arm_64" as const,
    memorySize: 1024,
    timeout: 600,
    permissions: [
      questionsTable,
      optionsTable
    ],
    bundling: { format: "binary" },
    environment: {
      STAGE: stack.stage,
      QUESTIONS_TABLE: questionsTable.tableName,
      OPTIONS_TABLE: optionsTable.tableName,
    },
  });

  const api = new Api(stack, "QuestionBankApi", {
    routes: {
      "POST /api/question/add": questionBankFunction,
      "GET /api/question/{questionId}": questionBankFunction,
      "PUT /api/question/{questionId}/update": questionBankFunction,
      "DELETE /api/question/{questionId}/delete": questionBankFunction,
    },
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}