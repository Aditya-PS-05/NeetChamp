import { SSTConfig } from "sst";
import { QuestionBankStack } from "./stacks/QuestionBankStack";

export default {
  config(_input: any) {
    return {
      name: "neetChamp-question-bank-service",
      region: "us-east-1",
    };
  },
  stacks(app: any) {
    app.stack(QuestionBankStack);
  },
} satisfies SSTConfig;