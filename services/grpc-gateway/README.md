+------------+         +-------------+        +------------+        +------------+
|  Frontend  | ----->  | API Gateway | -----> | QuizService| -----> |  Database  |
| (React)    | (gRPC)  | (Optional)  | (gRPC) | (Golang)   | (SQL)   | (Postgres) |
+------------+         +-------------+        +------------+        +------------+
                              |
                              v
                     +-----------------+
                     | Redis Cache      |
                     | (Quiz Sessions)  |
                     +-----------------+
