package repository

import (
	"bpos001/backend/internal/entity"
	"bpos001/backend/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }
func normalizePage(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func (r *Repository) listScope(resource, merchantID, branchID string) *gorm.DB {
	return r.db.Model(modelFor(resource)).Where("merchant_id = ? AND branch_id = ?", merchantID, branchID)
}
func modelFor(resource string) any {
	switch resource {
	case "posclients":
		return &entity.POSClient{}
	case "products":
		return &entity.POSProduct{}
	case "prices":
		return &entity.Price{}
	case "shifts":
		return &entity.ShiftSession{}
	case "sale-orders":
		return &entity.SaleOrder{}
	case "sale-items":
		return &entity.SaleItem{}
	case "payment-methods":
		return &entity.PaymentMethod{}
	case "payments":
		return &entity.Payment{}
	default:
		return nil
	}
}
func sliceFor(resource string) any {
	switch resource {
	case "posclients":
		return &[]entity.POSClient{}
	case "products":
		return &[]entity.POSProduct{}
	case "prices":
		return &[]entity.Price{}
	case "shifts":
		return &[]entity.ShiftSession{}
	case "sale-orders":
		return &[]entity.SaleOrder{}
	case "sale-items":
		return &[]entity.SaleItem{}
	case "payment-methods":
		return &[]entity.PaymentMethod{}
	case "payments":
		return &[]entity.Payment{}
	default:
		return nil
	}
}
func searchColumns(resource string) []string {
	switch resource {
	case "posclients":
		return []string{"pos_client_id", "client_name", "user_id", "status"}
	case "products":
		return []string{"sku_id", "sku_code", "barcode", "product_name"}
	case "prices":
		return []string{"price_id", "sku_id", "sku_code", "price_tier"}
	case "shifts":
		return []string{"shift_number", "pos_client_id", "status"}
	case "sale-orders":
		return []string{"receipt_number", "pos_client_id", "customer_id", "cashier_id", "status"}
	case "sale-items":
		return []string{"sku_id", "sku_code", "item_name"}
	case "payment-methods":
		return []string{"method_code", "method_name", "method_type"}
	case "payments":
		return []string{"payment_number", "status"}
	default:
		return nil
	}
}
func (r *Repository) List(resource, merchantID, branchID string, query model.ListQuery) (any, int64, int, int, error) {
	if modelFor(resource) == nil {
		return nil, 0, 0, 0, errors.New("unknown resource")
	}
	page, limit := normalizePage(query.Page, query.Limit)
	db := r.listScope(resource, merchantID, branchID)
	if query.Search != "" {
		like := "%" + strings.ToLower(query.Search) + "%"
		parts := []string{}
		vals := []any{}
		for _, col := range searchColumns(resource) {
			parts = append(parts, "LOWER("+col+") LIKE ?")
			vals = append(vals, like)
		}
		db = db.Where(strings.Join(parts, " OR "), vals...)
	}
	if query.Status != "" {
		db = db.Where("status = ?", strings.ToUpper(query.Status))
	}
	if query.POSClientID != "" {
		db = db.Where("pos_client_id = ?", query.POSClientID)
	}
	var total int64
	db.Count(&total)
	items := sliceFor(resource)
	err := db.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(items).Error
	return items, total, page, limit, err
}
func (r *Repository) Get(resource, merchantID, branchID, idOrCode string) (any, error) {
	m := modelFor(resource)
	if m == nil {
		return nil, errors.New("unknown resource")
	}
	cond := "merchant_id = ? AND branch_id = ? AND id::text = ?"
	args := []any{merchantID, branchID, idOrCode}
	switch resource {
	case "posclients":
		cond += " OR (merchant_id = ? AND branch_id = ? AND pos_client_id = ?)"
		args = append(args, merchantID, branchID, idOrCode)
	case "prices":
		cond += " OR (merchant_id = ? AND branch_id = ? AND price_id = ?)"
		args = append(args, merchantID, branchID, idOrCode)
	case "shifts":
		cond += " OR (merchant_id = ? AND branch_id = ? AND shift_number = ?)"
		args = append(args, merchantID, branchID, idOrCode)
	case "sale-orders":
		cond += " OR (merchant_id = ? AND branch_id = ? AND receipt_number = ?)"
		args = append(args, merchantID, branchID, idOrCode)
	case "payments":
		cond += " OR (merchant_id = ? AND branch_id = ? AND payment_number = ?)"
		args = append(args, merchantID, branchID, idOrCode)
	}
	db := r.db.Where(cond, args...)
	if resource == "sale-orders" {
		db = db.Preload("Items").Preload("Payments")
	}
	return m, db.First(m).Error
}
func (r *Repository) Create(resource string, item any) (any, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
func (r *Repository) Save(resource string, merchantID, branchID, id string, apply func(any) error) (any, error) {
	item, err := r.Get(resource, merchantID, branchID, id)
	if err != nil {
		return nil, err
	}
	if err := apply(item); err != nil {
		return nil, err
	}
	return item, r.db.Save(item).Error
}
func (r *Repository) Delete(resource, merchantID, branchID, id string) error {
	item, err := r.Get(resource, merchantID, branchID, id)
	if err != nil {
		return err
	}
	return r.db.Delete(item).Error
}
func (r *Repository) SelectPOS(merchantID, branchID, userID, idOrCode string) (entity.POSClient, error) {
	var item entity.POSClient
	err := r.db.Where("merchant_id = ? AND branch_id = ? AND is_active = ? AND (id::text = ? OR pos_client_id = ?)", merchantID, branchID, true, idOrCode, idOrCode).First(&item).Error
	return item, err
}
func (r *Repository) CurrentShift(merchantID, branchID, posClientID string) (entity.ShiftSession, error) {
	var item entity.ShiftSession
	err := r.db.Where("merchant_id = ? AND branch_id = ? AND pos_client_id = ? AND status = ?", merchantID, branchID, posClientID, "OPEN").Order("opened_at DESC").First(&item).Error
	return item, err
}
func (r *Repository) IncrementReceipt(merchantID, branchID, posClientID string) (string, error) {
	var pc entity.POSClient
	if err := r.db.Where("merchant_id = ? AND branch_id = ? AND pos_client_id = ?", merchantID, branchID, posClientID).First(&pc).Error; err != nil {
		return "", err
	}
	pc.RunningNumber++
	if err := r.db.Save(&pc).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%05d", pc.RunningPrefix, pc.RunningNumber), nil
}
func (r *Repository) EnsureShiftFresh(shift entity.ShiftSession) error {
	if shift.Status == "OPEN" && time.Since(shift.OpenedAt) > 24*time.Hour {
		return errors.New("shift is open longer than 24 hours and must be closed")
	}
	return nil
}

func (r *Repository) CloseExpiredShift(shift *entity.ShiftSession, userID string) error {
	now := time.Now()
	shift.Status = "CLOSED"
	shift.ClosedBy = userID
	shift.ClosedAt = &now
	shift.ClosingCashDifference = shift.ClosingCashActual - shift.ClosingCashExpected
	return r.db.Save(shift).Error
}
