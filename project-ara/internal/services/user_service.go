package services

import (
	"fmt"

	"gorm.io/gorm"

	"project-ara/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetOrCreateUser(phoneNumber string) (*models.User, error) {
	var user models.User

	// Try to find existing user
	if err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user
			user = models.User{
				PhoneNumber:            phoneNumber,
				TrialTransactionsCount: 0,
				SubscriptionStatus:     "trial",
			}

			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	return &user, nil
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserService) UpdateSubscriptionStatus(userID string, status string) error {
	return s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("subscription_status", status).
		Error
}

func (s *UserService) CanUserCreateTransaction(userID string) (bool, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	return user.CanCreateTransaction(), nil
}
