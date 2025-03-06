package models

import (
	authModels "github.com/Aditya-PS-05/NeetChamp/auth-service/models"
	"gorm.io/gorm"
)

// Alias User from auth-service
type User = authModels.User

type Mentee struct {
	UserID         string `gorm:"primaryKey"`
	Experience     int    `gorm:"default:0"` // XP points
	Level          int    `gorm:"default:1"` // User's level
	CurrentStreak  int    `gorm:"default:0"` // Consecutive active days
	LongestStreak  int    `gorm:"default:0"` // Highest streak ever
	QuizzesPlayed  int    `gorm:"default:0"` // Total quizzes attempted
	CorrectAnswers int    `gorm:"default:0"` // Number of correct answers
	Badges         string `gorm:"type:text"` // JSON string of earned badge IDs
	User           User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Admin struct {
	UserID       string `gorm:"primaryKey"`
	Permissions  string `gorm:"type:text"` // JSON string defining admin permissions
	ManagedUsers int    `gorm:"default:0"` // Number of users managed
	User         User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Mentee{}, &Admin{})
}
