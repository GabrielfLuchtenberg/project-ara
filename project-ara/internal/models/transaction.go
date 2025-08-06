package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionSource string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"

	TransactionSourceText  TransactionSource = "text"
	TransactionSourceVoice TransactionSource = "voice"
	TransactionSourceImage TransactionSource = "image"
)

type Transaction struct {
	ID              uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null" json:"user_id"`
	Amount          float64           `gorm:"type:decimal(10,2);not null" json:"amount"`
	Description     string            `gorm:"type:text" json:"description"`
	TransactionType TransactionType   `gorm:"type:varchar(10);not null" json:"transaction_type"`
	Source          TransactionSource `gorm:"type:varchar(20);not null" json:"source"`
	CreatedAt       time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CorrectedAt     *time.Time        `json:"corrected_at,omitempty"`
	CorrectionData  *json.RawMessage  `gorm:"type:jsonb" json:"correction_data,omitempty"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (t *Transaction) IsIncome() bool {
	return t.TransactionType == TransactionTypeIncome
}

func (t *Transaction) IsExpense() bool {
	return t.TransactionType == TransactionTypeExpense
}
