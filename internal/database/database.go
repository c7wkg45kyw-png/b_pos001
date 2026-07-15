package database

import (
	"bpos001/backend/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) { return gorm.Open(postgres.Open(dsn), &gorm.Config{}) }
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.POSClient{}, &entity.POSProduct{}, &entity.Price{}, &entity.ShiftSession{}, &entity.SaleOrder{}, &entity.SaleItem{}, &entity.PaymentMethod{}, &entity.Payment{})
}
