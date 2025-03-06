package controllers

import (
	"context"
	"errors"
	"strconv"

	pb "github.com/Aditya-PS-05/NeetChamp/shared-libs/proto"
	"github.com/Aditya-PS-05/NeetChamp/user-service/models"
	"gorm.io/gorm"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	DB *gorm.DB
}

func (c *UserController) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user models.User
	if err := c.DB.First(&user, "id = ?", req.UserId).Error; err != nil {
		return nil, errors.New("user not found")
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)

	response := &pb.GetUserResponse{
		UserId: userIDStr, // ✅ Convert uint to string
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}

	// Fetch Mentee or Admin details
	if user.Role == "mentee" {
		var mentee models.Mentee
		if err := c.DB.First(&mentee, "user_id = ?", req.UserId).Error; err == nil {
			response.UserDetails = &pb.GetUserResponse_Mentee{
				Mentee: &pb.Mentee{
					Experience:     int32(mentee.Experience),
					Level:          int32(mentee.Level),
					CurrentStreak:  int32(mentee.CurrentStreak),
					LongestStreak:  int32(mentee.LongestStreak),
					QuizzesPlayed:  int32(mentee.QuizzesPlayed),
					CorrectAnswers: int32(mentee.CorrectAnswers),
					Badges:         mentee.Badges,
				},
			}
		}
	} else if user.Role == "admin" {
		var admin models.Admin
		if err := c.DB.First(&admin, "user_id = ?", req.UserId).Error; err == nil {
			response.UserDetails = &pb.GetUserResponse_Admin{
				Admin: &pb.Admin{
					Permissions:  admin.Permissions,
					ManagedUsers: int32(admin.ManagedUsers),
				},
			}
		}
	}

	return response, nil
}

// UpdateUser handles updating both mentee & admin data
func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var user models.User
	if err := c.DB.First(&user, "id = ?", req.UserId).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	c.DB.Save(&user)

	// Update Mentee
	if req.GetMentee() != nil {
		var mentee models.Mentee
		if err := c.DB.First(&mentee, "user_id = ?", req.UserId).Error; err == nil {
			mentee.Experience = int(req.GetMentee().Experience)
			mentee.Level = int(req.GetMentee().Level)
			mentee.CurrentStreak = int(req.GetMentee().CurrentStreak)
			mentee.LongestStreak = int(req.GetMentee().LongestStreak)
			mentee.QuizzesPlayed = int(req.GetMentee().QuizzesPlayed)
			mentee.CorrectAnswers = int(req.GetMentee().CorrectAnswers)
			mentee.Badges = req.GetMentee().Badges
			c.DB.Save(&mentee)
		}
	}

	// Update Admin
	if req.GetAdmin() != nil {
		var admin models.Admin
		if err := c.DB.First(&admin, "user_id = ?", req.UserId).Error; err == nil {
			admin.Permissions = req.GetAdmin().Permissions
			admin.ManagedUsers = int(req.GetAdmin().ManagedUsers)
			c.DB.Save(&admin)
		}
	}

	return &pb.UpdateUserResponse{Message: "User updated successfully"}, nil
}

// DeleteUser handles deleting both mentee & admin
func (c *UserController) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := c.DB.Delete(&models.User{}, "id = ?", req.UserId).Error; err != nil {
		return nil, errors.New("failed to delete user")
	}
	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
