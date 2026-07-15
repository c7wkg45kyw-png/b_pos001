package usecase

import (
	"bpos001/backend/internal/entity"
	"bpos001/backend/internal/model"
	"bpos001/backend/internal/repository"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Usecase struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Usecase { return &Usecase{repo: repo} }
func (u *Usecase) List(auth model.AuthContext, resource string, query model.ListQuery) (model.PageResponse[any], error) {
	items, total, page, limit, err := u.repo.List(resource, auth.MerchantID, auth.BranchID, query)
	if err != nil {
		return model.PageResponse[any]{}, err
	}
	return model.PageResponse[any]{Items: []any{items}, Page: page, Limit: limit, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(limit)))}, nil
}
func (u *Usecase) Get(auth model.AuthContext, resource, id string) (any, error) {
	return u.repo.Get(resource, auth.MerchantID, auth.BranchID, id)
}
func (u *Usecase) Delete(auth model.AuthContext, resource, id string) error {
	return u.repo.Delete(resource, auth.MerchantID, auth.BranchID, id)
}
func base(auth model.AuthContext) entity.Base {
	return entity.Base{MerchantID: auth.MerchantID, BranchID: auth.BranchID, CreatedBy: auth.UserID, UpdatedBy: auth.UserID}
}
func (u *Usecase) Create(auth model.AuthContext, resource string, payload map[string]any) (any, error) {
	item, err := mapPayload(resource, payload)
	if err != nil {
		return nil, err
	}
	setBase(item, base(auth))
	if resource == "shifts" {
		sh := item.(*entity.ShiftSession)
		if sh.OpenedBy == "" {
			sh.OpenedBy = auth.UserID
		}
	}
	if resource == "sale-orders" {
		so := item.(*entity.SaleOrder)
		if so.POSClientID == "" {
			return nil, errors.New("posclient_id is required")
		}
		if so.ReceiptNumber == "" {
			rn, err := u.repo.IncrementReceipt(auth.MerchantID, auth.BranchID, so.POSClientID)
			if err == nil {
				so.ReceiptNumber = rn
			}
		}
		if so.CashierID == "" {
			so.CashierID = auth.UserID
		}
		calculateSO(so)
	}
	return u.repo.Create(resource, item)
}
func (u *Usecase) Update(auth model.AuthContext, resource, id string, payload map[string]any) (any, error) {
	return u.repo.Save(resource, auth.MerchantID, auth.BranchID, id, func(item any) error {
		if err := applyPayload(item, payload); err != nil {
			return err
		}
		setUpdated(item, auth.UserID)
		if so, ok := item.(*entity.SaleOrder); ok {
			calculateSO(so)
		}
		return nil
	})
}
func (u *Usecase) SelectPOS(auth model.AuthContext, idOrCode string) (entity.POSClient, error) {
	return u.repo.SelectPOS(auth.MerchantID, auth.BranchID, auth.UserID, idOrCode)
}
func (u *Usecase) GetShiftDetails(auth model.AuthContext, posClientID string) (map[string]any, error) {
	shift, err := u.repo.CurrentShift(auth.MerchantID, auth.BranchID, posClientID)
	if err != nil {
		return map[string]any{"status": "NO_OPEN_SHIFT", "posclient_id": posClientID}, nil
	}
	expired := time.Since(shift.OpenedAt) > 24*time.Hour
	if expired {
		if err := u.repo.CloseExpiredShift(&shift, auth.UserID); err != nil {
			return nil, err
		}
		return map[string]any{"status": "CLOSED", "must_close": true, "reason": "shift exceeded 24 hours", "shift": shift}, nil
	}
	return map[string]any{"status": shift.Status, "must_close": false, "shift": shift}, nil
}
func mapPayload(resource string, p map[string]any) (any, error) {
	switch resource {
	case "posclients":
		item := &entity.POSClient{}
		applyPayload(item, p)
		return item, nil
	case "products":
		item := &entity.POSProduct{}
		applyPayload(item, p)
		return item, nil
	case "prices":
		item := &entity.Price{}
		applyPayload(item, p)
		return item, nil
	case "shifts":
		item := &entity.ShiftSession{OpenedAt: time.Now(), Status: "OPEN"}
		applyPayload(item, p)
		if item.ShiftNumber == "" {
			item.ShiftNumber = fmt.Sprintf("SHIFT-%d", time.Now().Unix())
		}
		return item, nil
	case "sale-orders":
		item := &entity.SaleOrder{VATRate: 7, Status: "SUCCESS"}
		applyPayload(item, p)
		return item, nil
	case "sale-items":
		item := &entity.SaleItem{}
		applyPayload(item, p)
		return item, nil
	case "payment-methods":
		item := &entity.PaymentMethod{IsActive: true}
		applyPayload(item, p)
		return item, nil
	case "payments":
		item := &entity.Payment{Status: "SUCCESS"}
		applyPayload(item, p)
		if item.PaidAt == nil {
			now := time.Now()
			item.PaidAt = &now
		}
		return item, nil
	default:
		return nil, errors.New("unknown resource")
	}
}
func setBase(item any, b entity.Base) {
	switch v := item.(type) {
	case *entity.POSClient:
		v.Base = b
	case *entity.POSProduct:
		v.Base = b
	case *entity.Price:
		v.Base = b
	case *entity.ShiftSession:
		v.Base = b
	case *entity.SaleOrder:
		v.Base = b
	case *entity.SaleItem:
		v.Base = b
	case *entity.PaymentMethod:
		v.Base = b
	case *entity.Payment:
		v.Base = b
	}
}
func setUpdated(item any, userID string) {
	switch v := item.(type) {
	case *entity.POSClient:
		v.UpdatedBy = userID
	case *entity.POSProduct:
		v.UpdatedBy = userID
	case *entity.Price:
		v.UpdatedBy = userID
	case *entity.ShiftSession:
		v.UpdatedBy = userID
	case *entity.SaleOrder:
		v.UpdatedBy = userID
	case *entity.SaleItem:
		v.UpdatedBy = userID
	case *entity.PaymentMethod:
		v.UpdatedBy = userID
	case *entity.Payment:
		v.UpdatedBy = userID
	}
}
func str(p map[string]any, k string) string {
	if v, ok := p[k]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}
func flt(p map[string]any, k string) float64 {
	if v, ok := p[k]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		case string:
			var f float64
			fmt.Sscan(n, &f)
			return f
		}
	}
	return 0
}
func boo(p map[string]any, k string, fallback bool) bool {
	if v, ok := p[k]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
		return strings.EqualFold(fmt.Sprint(v), "true")
	}
	return fallback
}
func applyPayload(item any, p map[string]any) error {
	switch v := item.(type) {
	case *entity.POSClient:
		if s := str(p, "posclient_id"); s != "" {
			v.POSClientID = s
		}
		if s := str(p, "user_id"); s != "" {
			v.UserID = s
		}
		if s := str(p, "client_name"); s != "" {
			v.ClientName = s
		}
		if s := str(p, "running_prefix"); s != "" {
			v.RunningPrefix = s
		}
		if s := str(p, "status"); s != "" {
			v.Status = strings.ToUpper(s)
		}
		v.IsActive = boo(p, "is_active", true)
	case *entity.POSProduct:
		if s := str(p, "sku_id"); s != "" {
			v.SKUID = s
		}
		if s := str(p, "sku_code"); s != "" {
			v.SKUCode = s
		}
		if s := str(p, "barcode"); s != "" {
			v.Barcode = s
		}
		if s := str(p, "product_name"); s != "" {
			v.ProductName = s
		}
		if s := str(p, "unit"); s != "" {
			v.Unit = s
		}
		v.IsActive = boo(p, "is_active", true)
	case *entity.Price:
		if s := str(p, "price_id"); s != "" {
			v.PriceID = s
		}
		if s := str(p, "sku_id"); s != "" {
			v.SKUID = s
		}
		if s := str(p, "sku_code"); s != "" {
			v.SKUCode = s
		}
		if s := str(p, "price_tier"); s != "" {
			v.PriceTier = s
		}
		if s := str(p, "currency"); s != "" {
			v.Currency = s
		}
		if _, ok := p["unit_price"]; ok {
			v.UnitPrice = flt(p, "unit_price")
		}
		v.IsActive = boo(p, "is_active", true)
	case *entity.ShiftSession:
		if s := str(p, "shift_number"); s != "" {
			v.ShiftNumber = s
		}
		if s := str(p, "posclient_id"); s != "" {
			v.POSClientID = s
		}
		if s := str(p, "opened_by"); s != "" {
			v.OpenedBy = s
		}
		if s := str(p, "closed_by"); s != "" {
			v.ClosedBy = s
		}
		if _, ok := p["opening_cash"]; ok {
			v.OpeningCash = flt(p, "opening_cash")
		}
		if _, ok := p["closing_cash_actual"]; ok {
			v.ClosingCashActual = flt(p, "closing_cash_actual")
		}
		if s := str(p, "status"); s != "" {
			v.Status = strings.ToUpper(s)
		}
	case *entity.SaleOrder:
		if s := str(p, "receipt_number"); s != "" {
			v.ReceiptNumber = s
		}
		if s := str(p, "posclient_id"); s != "" {
			v.POSClientID = s
		}
		if s := str(p, "shift_session_id"); s != "" {
			v.ShiftSessionID = s
		}
		if s := str(p, "customer_id"); s != "" {
			v.CustomerID = s
		}
		if s := str(p, "cashier_id"); s != "" {
			v.CashierID = s
		}
		if _, ok := p["subtotal"]; ok {
			v.Subtotal = flt(p, "subtotal")
		}
		if _, ok := p["discount"]; ok {
			v.Discount = flt(p, "discount")
		}
		if _, ok := p["vat_rate"]; ok {
			v.VATRate = flt(p, "vat_rate")
		}
		if s := str(p, "status"); s != "" {
			v.Status = strings.ToUpper(s)
		}
	case *entity.SaleItem:
		if s := str(p, "sale_order_id"); s != "" {
			v.SaleOrderID = s
		}
		if s := str(p, "sku_id"); s != "" {
			v.SKUID = s
		}
		if s := str(p, "sku_code"); s != "" {
			v.SKUCode = s
		}
		if s := str(p, "item_name"); s != "" {
			v.ItemName = s
		}
		v.Quantity = flt(p, "quantity")
		v.UnitPrice = flt(p, "unit_price")
		v.Discount = flt(p, "discount")
		v.TotalPrice = (v.Quantity * v.UnitPrice) - v.Discount
	case *entity.PaymentMethod:
		if s := str(p, "method_code"); s != "" {
			v.MethodCode = s
		}
		if s := str(p, "method_name"); s != "" {
			v.MethodName = s
		}
		if s := str(p, "method_type"); s != "" {
			v.MethodType = strings.ToUpper(s)
		}
		if s := str(p, "promptpay_target"); s != "" {
			v.PromptPayTarget = s
		}
		v.IsActive = boo(p, "is_active", true)
	case *entity.Payment:
		if s := str(p, "payment_number"); s != "" {
			v.PaymentNumber = s
		}
		if s := str(p, "sale_order_id"); s != "" {
			v.SaleOrderID = s
		}
		if s := str(p, "payment_method_id"); s != "" {
			v.PaymentMethodID = s
		}
		if _, ok := p["amount"]; ok {
			v.Amount = flt(p, "amount")
		}
		if s := str(p, "status"); s != "" {
			v.Status = strings.ToUpper(s)
		}
		if s := str(p, "evidence_url"); s != "" {
			v.EvidenceURL = s
		}
	default:
		return errors.New("unknown payload target")
	}
	return nil
}
func calculateSO(v *entity.SaleOrder) {
	if v.VATRate == 0 {
		v.VATRate = 7
	}
	net := v.Subtotal - v.Discount
	if net < 0 {
		net = 0
	}
	v.VAT = net * v.VATRate / 100
	v.GrandTotal = net + v.VAT
}
