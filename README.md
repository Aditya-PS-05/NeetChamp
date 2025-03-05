# **Neet Quiz Monorepo**

## **Overview**
Neet Quiz is a **microservices-based** quiz platform designed to provide an interactive and gamified experience for students preparing for NEET. This monorepo contains all **Golang microservices** and supports **gRPC, REST APIs, and AI-driven adaptive learning**.

## **Project Structure**
```
NeetChamp/
│── apps/                          # Empty UI folder (Flutter UI handled separately)
│
│── services/                      # Microservices Directory
│   ├── auth-service/               # Authentication & User Management
│   ├── user-service/               # Handles User Profiles
│   ├── quiz-service/               # Manages Quizzes
│   ├── question-bank-service/      # Stores & Generates Questions
│   ├── live-quiz-service/          # WebSockets for Real-time Quizzes
│   ├── leaderboard-service/        # Ranking & Performance Analytics
│   ├── gamification-service/       # XP, Badges, Streaks
│   ├── results-report-service/     # Generates Reports & Analytics
│   ├── notification-service/       # Email, Push Notifications
│   ├── payment-service/            # Subscription & Monetization
│   ├── ai-service/                 # AI-driven Adaptive Learning
│   ├── grpc-gateway/               # API Gateway for gRPC (Optional)
│
│── shared-libs/                    # Shared Code & Utilities
│   ├── proto/                       # gRPC Protobuf Definitions
│   │   ├── auth.proto
│   │   ├── quiz.proto
│   ├── common/                      # Shared utility functions
│   ├── database/                    # DB connection pools for reuse
│   ├── logger/                      # Shared logging functions
│
│── deployment/                      # DevOps & CI/CD Configs
│   ├── docker-compose.yml
│   ├── k8s/                          # Kubernetes YAMLs
│   ├── terraform/                    # Infrastructure as Code
│
│── docs/                             # Documentation & API References
│   ├── README.md
│   ├── API.md
│
│── .gitignore
│── package.json
│── go.mod
│── README.md
```

## **Installation & Setup**
### **Prerequisites**
- Go 1.17+
- Docker & Docker Compose
- Kubernetes (Optional for deployment)
- Mage (Task runner for automation)

### **Clone the Repository**
```sh
git clone https://github.com/your-username/NeetChamp.git
cd NeetChamp
```

### **Initialize the Project**
```sh
mage InitProject
mage InitGoModules
```

### **Build All Services**
```sh
mage Build
```

### **Run All Services Using Docker Compose**
```sh
mage Run
```

### **Run Tests for All Services**
```sh
mage Test
```

### **Deploy to Kubernetes**
```sh
mage Deploy
```

## **Microservices Overview**

| Service | Description |
|---------|------------|
| **Auth Service** | Handles user authentication (JWT, OAuth) |
| **User Service** | Manages user profiles and roles |
| **Quiz Service** | Creates and manages quizzes |
| **Question Bank Service** | Stores and generates quiz questions |
| **Live Quiz Service** | Manages real-time quiz sessions using WebSockets |
| **Leaderboard Service** | Tracks rankings and performance metrics |
| **Gamification Service** | XP, achievements, streaks, and rewards |
| **Results & Reporting Service** | Generates student and teacher reports |
| **Notification Service** | Sends push notifications and emails |
| **Payment Service** | Handles subscriptions and premium purchases |
| **AI Service** | Provides adaptive learning features |
| **gRPC Gateway** | API gateway for internal gRPC communication |

## **API Documentation**
Refer to the [API.md](docs/API.md) for detailed API specifications.

## **Contributing**
1. Fork the repository.
2. Create a new feature branch.
3. Commit changes and push to your branch.
4. Open a Pull Request for review.

## **License**
MIT License. See [LICENSE](LICENSE) for details.
