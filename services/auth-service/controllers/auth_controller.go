package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/database"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/models"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/utils"
	"github.com/Aditya-PS-05/NeetChamp/shared-libs/proto/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
}

// ✅ Register user with transaction & duplicate email check
func (s *AuthServiceServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("❌ email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("❌ failed to hash password")
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	// ✅ Use transaction to prevent partial failures
	tx := database.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("❌ failed to create user")
	}
	tx.Commit()

	return &auth.RegisterResponse{
		UserId:  fmt.Sprintf("%d", user.ID),
		Message: "✅ User registered successfully",
	}, nil
}

// ✅ Login user with Redis caching & optimized error handling
func (s *AuthServiceServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user models.User

	// ✅ Check Redis cache first
	cachedUser, err := utils.GetCachedUser(req.Email)
	if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(cachedUser.Password), []byte(req.Password)) == nil {
			token, _ := utils.GenerateToken(req.Email, cachedUser.Role)
			return &auth.LoginResponse{Token: token}, nil
		}
		return nil, errors.New("❌ invalid credentials")
	}

	// ✅ Fetch from DB if cache miss
	if err := database.DB.Select("id, password, role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("❌ user not found")
		}
		return nil, errors.New("❌ database error")
	}

	// ✅ Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("❌ invalid credentials")
	}

	// ✅ Cache user in Redis
	utils.CacheUser(req.Email, user.Password, user.Role)

	token, _ := utils.GenerateToken(req.Email, user.Role)

	// ✅ Prevent login if token is blacklisted
	if utils.IsTokenBlacklisted(token) {
		return nil, errors.New("❌ token is invalid or expired")
	}

	return &auth.LoginResponse{Token: token}, nil
}

// ✅ Logout user & blacklist token
func (s *AuthServiceServer) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	err := utils.SaveTokenToBlacklist(req.Token, 24*60*60) // Blacklist token for 24 hours
	if err != nil {
		return nil, errors.New("❌ failed to log out")
	}
	return &auth.LogoutResponse{Message: "✅ User logged out successfully"}, nil
}

// ✅ Fetch authenticated user details
func (s *AuthServiceServer) GetAuthUser(ctx context.Context, req *auth.GetAuthUserRequest) (*auth.GetAuthUserResponse, error) {
	var user models.User

	if err := database.DB.Select("id, name, email, role").Where("id = ?", req.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("❌ user not found")
		}
		return nil, errors.New("❌ database error")
	}

	return &auth.GetAuthUserResponse{
		UserId: fmt.Sprintf("%d", user.ID),
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}
