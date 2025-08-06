package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PhoneNumber            string     `gorm:"type:varchar(20);unique;not null" json:"phone_number"`
	CreatedAt              time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	TrialTransactionsCount int        `gorm:"default:0" json:"trial_transactions_count"`
	SubscriptionStatus     string     `gorm:"type:varchar(20);default:'trial'" json:"subscription_status"`
	SubscriptionExpiresAt  *time.Time `json:"subscription_expires_at,omitempty"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *User) IsTrialExpired() bool {
	return u.TrialTransactionsCount >= 50
}

func (u *User) CanCreateTransaction() bool {
	if u.SubscriptionStatus == "active" {
		return true
	}
	return u.TrialTransactionsCount < 50
}
