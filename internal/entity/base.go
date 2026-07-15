package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	MerchantID string         `gorm:"size:80;not null;index" json:"merchant_id"`
	BranchID   string         `gorm:"size:80;not null;index" json:"branch_id"`
	CreatedBy  string         `gorm:"size:120" json:"-"`
	UpdatedBy  string         `gorm:"size:120" json:"-"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
