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
  --from-literal=DB_HOST= \
  --from-literal=DB_USER= \
  --from-literal=DB_PASSWORD= \
  --from-literal=DB_NAME=NeetChamp \
  --from-literal=DB_PORT=5432
```

```sh
kubectl create secret generic redis-secret \
  --from-literal=REDIS_HOST= \
  --from-literal=REDIS_PORT=6379 \
  --from-literal=REDIS_PASSWORD= \
  --from-literal=REDIS_TLS=true
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

## 🏗️ Deploying to AWS EKS

### **1️⃣ Install AWS CLI & eksctl**
```sh
curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
sudo installer -pkg AWSCLIV2.pkg -target /

curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_Linux_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

### **2️⃣ Configure AWS CLI**
```sh
aws configure
```
Enter AWS credentials and verify:
```sh
aws sts get-caller-identity
```

### **3️⃣ Create an EKS Cluster**
```sh
eksctl create cluster --name auth-cluster --region us-east-1 --nodegroup-name auth-nodes --node-type t3.medium --nodes 3
```

### **4️⃣ Deploy Auth Service on EKS**
```sh
kubectl apply -f deployment/k8s/auth-service-deployment.yaml
kubectl expose deployment auth-service --type=LoadBalancer --port=50051 --target-port=50051
```

### **5️⃣ Get Load Balancer External IP**
```sh
kubectl get services
```
Use the **EXTERNAL-IP** in load testing.

### **6️⃣ Run Load Test with `ghz`**
```sh
ghz --insecure --proto shared-libs/proto/auth.proto --call auth.AuthService/LoginUser -d '{ "email": "john@example.com", "password": "password123" }' -n 10000 -c 100 --connections=100 <EXTERNAL-IP>:50051
```

### **7️⃣ Check Load Balancing & Resource Usage**
```sh
kubectl top pods
```

### **8️⃣ Enable Auto-Scaling on EKS**
```sh
kubectl autoscale deployment auth-service --cpu-percent=50 --min=2 --max=10
kubectl get hpa
```

### **9️⃣ Delete Cluster (When Done)**
```sh
eksctl delete cluster --name auth-cluster
```

---

## 📊 Load Testing with `ghz`

### **1️⃣ Install `ghz`**
```sh
wget https://github.com/bojand/ghz/releases/download/v0.108.0/ghz-linux-amd64
chmod +x ghz-linux-amd64
sudo mv ghz-linux-amd64 /usr/local/bin/ghz
```

### **2️⃣ Run Load Test**
```sh
ghz --insecure --proto shared-libs/proto/auth.proto --call auth.AuthService/LoginUser -d '{ "email": "john@example.com", "password": "password123" }' -n 10000 -c 100 --connections=100 <EXTERNAL-IP>:50051
```

---

## 📌 Next Steps
🔹 Ensure **all replicas are utilized** using `--connections=100` in `ghz`.
🔹 Deploy on **AWS EKS** to leverage real **multi-node scaling**.
🔹 Implement **refresh tokens** for better security.
🔹 Set up **CI/CD pipeline** using GitHub Actions.

---

## 📄 License
MIT License © 2025 Aditya Pratap Singh

---


### Database Structure for auth service

 


 

``` SQL query
 

CREATE TABLE users (
 

    id SERIAL PRIMARY KEY,
 

    email VARCHAR(255) UNIQUE NOT NULL,
 

    password_hash VARCHAR(255),
 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 

    last_login_at TIMESTAMP,
 

	CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
 

);
 


 

CREATE TABLE auth_providers (
 

    id SERIAL PRIMARY KEY,
 

    user_id INTEGER NOT NULL REFERENCES users(id),
 

    provider VARCHAR(50) CHECK (provider IN ('google','phone', 'email')),
 

    provider_user_id VARCHAR(255),
 

    access_token TEXT,
 

    refresh_token TEXT,
 

    expires_at TIMESTAMP,
 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 

    UNIQUE (provider, provider_user_id),
 

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
 

);
 

```
 


 
