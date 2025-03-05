package controllers

import (
	"context"
	"errors"

	"github.com/Aditya-PS-05/NeetChamp/shared-libs/proto"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/database"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/models"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceServer struct {
	proto.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return nil, errors.New("failed to create user")
	}

	return &proto.RegisterResponse{
		UserId:  string(user.ID),
		Message: "User registered successfully",
	}, nil
}

func (s *AuthServiceServer) LoginUser(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, _ := utils.GenerateToken(user.Email, user.Role)

	// 🔹 Prevent login with a blacklisted token
	if utils.IsTokenBlacklisted(token) {
		return nil, errors.New("token is invalid or expired")
	}

	return &proto.LoginResponse{Token: token}, nil
}

func (s *AuthServiceServer) LogoutUser(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	err := utils.SaveTokenToBlacklist(req.Token, 24*60*60) // 🔹 Blacklist token for 24 hours
	if err != nil {
		return nil, errors.New("failed to log out")
	}
	return &proto.LogoutResponse{Message: "User logged out successfully"}, nil
}
