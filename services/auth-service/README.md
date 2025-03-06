# 🚀 Auth Service - NeetChamp

## 📌 Overview
The **Auth Service** is a gRPC-based authentication microservice for **NeetChamp**. It handles **user registration, login, and logout**, with **JWT-based authentication** and **Redis-powered token blacklisting**.

---

## 📂 Folder Structure
```
NeetChamp/
│── services/
│   ├── auth-service/
│   │   ├── main.go                  # Entry point for gRPC server
│   │   ├── controllers/
│   │   │   ├── auth_controller.go   # Handles authentication logic
│   │   ├── database/
│   │   │   ├── db.go                 # Database connection (PostgreSQL)
│   │   ├── models/
│   │   │   ├── user.go               # User model
│   │   ├── utils/
│   │   │   ├── jwt.go                # JWT token generation & validation
│   │   │   ├── redis.go              # Redis token blacklisting
│   │   ├── .env                      # Environment variables
│   │   ├── Dockerfile                # Docker build instructions
│   │   ├── auth-service-deployment.yaml # Kubernetes Deployment file
│── shared-libs/
│   ├── proto/
│   │   ├── auth.proto                 # gRPC Protobuf definitions
```

---

## 🛠️ How to Build & Run

### **1⃣ Clone the repository**
```sh
git clone https://github.com/Aditya-PS-05/NeetChamp.git
cd NeetChamp/services/auth-service
```

### **2⃣ Set up Environment Variables**
Create a **.env** file inside `auth-service/`:
```env
DB_HOST=your-database-host
DB_USER=your-database-user
DB_PASSWORD=your-database-password
DB_NAME=NeetChamp
DB_PORT=5432
REDIS_HOST=splendid-newt-59229.upstash.io
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_TLS=true
```

### **3⃣ Run Locally**
```sh
go mod tidy
go run main.go
```

---

## 🐓 Running with Docker

### **1⃣ Build Docker Image**
```sh
docker build -t ghcr.io/aditya-ps-05/neetchamp/auth-service:latest -f Dockerfile .
```

### **2⃣ Run Docker Container**
```sh
docker run -p 50051:50051 --env-file .env ghcr.io/aditya-ps-05/neetchamp/auth-service:latest
```

---

## 🚀 Deploying to Kubernetes

### **1⃣ Install Kubernetes & Metrics Server**
```sh
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### **2⃣ Deploy Auth Service**
```sh
kubectl apply -f deployment/k8s/auth-service-deployment.yaml
```

### **3⃣ Expose Service**
```sh
kubectl expose deployment auth-service --type=LoadBalancer --port=80 --target-port=50051
```

### **4⃣ Enable Autoscaling**
```sh
kubectl autoscale deployment auth-service --cpu-percent=50 --min=2 --max=20
```

### **5⃣ Check Deployment Status**
```sh
kubectl get pods -w
kubectl get services
kubectl get hpa
```

---

## 📊 Load Testing with `ghz`

### **1⃣ Install `ghz`**
```sh
wget https://github.com/bojand/ghz/releases/download/v0.108.0/ghz-linux-amd64
chmod +x ghz-linux-amd64
sudo mv ghz-linux-amd64 /usr/local/bin/ghz
```

### **2⃣ Run Load Test**
```sh
ghz --insecure \
    --proto shared-libs/proto/auth.proto \
    --call auth.AuthService/LoginUser \
    -d '{ "email": "john@example.com", "password": "password123" }' \
    -n 10000 -c 100 \
    localhost:50051
```

---

## 🔾 Secrets Management in Kubernetes

### **1⃣ Delete Existing Secrets**
```sh
kubectl delete secret db-secret
```

### **2⃣ Create a New Secret**
```sh
kubectl create secret generic db-secret \
  --from-literal=DB_HOST=ep-morning-bar-a8d00toz-pooler.eastus2.azure.neon.tech \
  --from-literal=DB_USER=NeetChamp_owner \
  --from-literal=DB_PASSWORD=npg_U2NlWGCgTu8t \
  --from-literal=DB_NAME=NeetChamp \
  --from-literal=DB_PORT=5432
```

### **3⃣ Verify Secret Creation**
```sh
kubectl get secrets
kubectl describe secret db-secret
```

### **4⃣ Restart Deployment to Apply Changes**
```sh
kubectl rollout restart deployment auth-service
```

---

## 🔍 gRPC API Methods (Protobuf Definition)

```proto
syntax = "proto3";
package auth;
option go_package = "github.com/Aditya-PS-05/NeetChamp/shared-libs/proto";

service AuthService {
  rpc RegisterUser(RegisterRequest) returns (RegisterResponse);
  rpc LoginUser(LoginRequest) returns (LoginResponse);
  rpc LogoutUser(LogoutRequest) returns (LogoutResponse);
}

message RegisterRequest {
  string name = 1;
  string email = 2;
  string password = 3;
  string role = 4;
}

message RegisterResponse {
  string user_id = 1;
  string message = 2;
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginResponse {
  string token = 1;
}

message LogoutRequest {
  string token = 1;
}

message LogoutResponse {
  string message = 1;
}
```

---

## 🔥 Optimizations Implemented
👉 **Optimized gRPC Service**
- Used **transactions** in DB to prevent partial failures.
- **Cached queries** using Redis.
- **Prevented duplicate logins** using token blacklist.

👉 **Improved Performance**
- **Load Testing Results**
  - **Before:** 39 requests/sec
  - **After Optimizations:** **77 requests/sec** 🚀

👉 **Auto-Scaling & Kubernetes**
- Deployed on **Kubernetes** with **Horizontal Pod Autoscaler (HPA)**.
- Autoscaling adjusts replicas dynamically.

---

## 📈 Next Steps
🔹 Implement **refresh tokens** for better security.
🔹 Improve **database indexing** for faster queries.
🔹 Set up **CI/CD pipeline** using GitHub Actions.

---

## 📝 License
MIT License © 2025 Aditya Pratap Singh

