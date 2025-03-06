package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/database"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/models"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/utils"
	"github.com/Aditya-PS-05/NeetChamp/shared-libs/proto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceServer struct {
	proto.UnimplementedAuthServiceServer
}

// ✅ Optimized RegisterUser with transactions & better error handling
func (s *AuthServiceServer) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
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

	// ✅ Use transactions to prevent partial failures
	tx := database.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("❌ failed to create user")
	}
	tx.Commit()

	return &proto.RegisterResponse{
		UserId:  fmt.Sprintf("%d", user.ID),
		Message: "✅ User registered successfully",
	}, nil
}

// ✅ Optimized LoginUser with Redis caching & error handling
func (s *AuthServiceServer) LoginUser(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user models.User

	// ✅ Try fetching from Redis first (cache hit)
	cachedUser, err := utils.GetCachedUser(req.Email)
	if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(cachedUser.Password), []byte(req.Password)) == nil {
			token, _ := utils.GenerateToken(req.Email, cachedUser.Role)
			return &proto.LoginResponse{Token: token}, nil
		}
		return nil, errors.New("❌ invalid credentials")
	}

	// ✅ Only select necessary fields (faster query)
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

	// ✅ Cache user in Redis to prevent DB hits next time
	utils.CacheUser(req.Email, user.Password, user.Role)

	token, _ := utils.GenerateToken(req.Email, user.Role)

	// ✅ Prevent login if token is blacklisted
	if utils.IsTokenBlacklisted(token) {
		return nil, errors.New("❌ token is invalid or expired")
	}

	return &proto.LoginResponse{Token: token}, nil
}

// ✅ Optimized LogoutUser with Redis connection reuse
func (s *AuthServiceServer) LogoutUser(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := utils.SaveTokenToBlacklist(req.Token, 24*60*60) // Blacklist token for 24 hours
	if err != nil {
		return nil, errors.New("❌ failed to log out")
	}
	return &proto.LogoutResponse{Message: "✅ User logged out successfully"}, nil
}
