//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// List of services
var services = []string{
	"auth-service",
	"user-service",
	"quiz-service",
	"question-bank-service",
	"live-quiz-service",
	"leaderboard-service",
	"gamification-service",
	"results-report-service",
	"notification-service",
	"payment-service",
	"ai-service",
	"grpc-gateway",
}

// Create the entire directory structure
func InitProject() error {
	fmt.Println("Initializing project structure...")

	dirs := []string{
		"apps",
		"services",
		"shared-libs/proto",
		"shared-libs/common",
		"shared-libs/database",
		"shared-libs/logger",
		"deployment/k8s",
		"deployment/terraform",
		"docs",
	}

	// Create directories
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// Create microservices directories
	for _, service := range services {
		servicePath := fmt.Sprintf("services/%s", service)
		subDirs := []string{"models", "controllers", "routes", "utils"}
		for _, subDir := range subDirs {
			if err := os.MkdirAll(fmt.Sprintf("%s/%s", servicePath, subDir), os.ModePerm); err != nil {
				return err
			}
		}
		files := []string{"main.go", ".env", "Dockerfile"}
		for _, file := range files {
			f, err := os.Create(fmt.Sprintf("%s/%s", servicePath, file))
			if err != nil {
				return err
			}
			f.Close()
		}
	}

	// Create necessary files
	files := []string{
		".gitignore", "package.json", "go.mod", "README.md",
		"deployment/docker-compose.yml",
		"docs/README.md", "docs/API.md",
		"shared-libs/proto/auth.proto",
		"shared-libs/proto/quiz.proto",
	}

	for _, file := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		f.Close()
	}

	fmt.Println("Project structure initialized successfully!")
	return nil
}

// Initialize Go modules for each service
func InitGoModules() error {
	fmt.Println("Initializing Go modules for all microservices...")

	for _, service := range services {
		servicePath := fmt.Sprintf("services/%s", service)
		cmd := exec.Command("go", "mod", "init", fmt.Sprintf("neet-quiz-monorepo/%s", service))
		cmd.Dir = servicePath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("Go modules initialized successfully!")
	return nil
}

// Build all services
func Build() error {
	fmt.Println("Building all microservices...")

	for _, service := range services {
		fmt.Printf("🚀 Building %s...\n", service)
		cmd := exec.Command("go", "build", "-o", fmt.Sprintf("bin/%s", service), "./main.go")
		cmd.Dir = fmt.Sprintf("services/%s", service)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("Build completed successfully!")
	return nil
}

// Run all services using Docker Compose
func Run() error {
	fmt.Println("Starting all services using Docker Compose...")
	return exec.Command("docker-compose", "up", "--build").Run()
}

// Run tests for all services
func Test() error {
	fmt.Println("Running tests for all microservices...")

	for _, service := range services {
		fmt.Printf("🔍 Testing %s...\n", service)
		cmd := exec.Command("go", "test", "./...")
		cmd.Dir = fmt.Sprintf("services/%s", service)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("All tests completed successfully!")
	return nil
}

// Deploy to Kubernetes
func Deploy() error {
	fmt.Println("🚀 Deploying to Kubernetes...")
	return exec.Command("kubectl", "apply", "-f", "deployment/k8s/").Run()
}
