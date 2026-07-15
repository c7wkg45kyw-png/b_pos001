package entity

import "time"

type POSClient struct {
	Base
	POSClientID   string `gorm:"size:80;not null;uniqueIndex:idx_posclient_code" json:"posclient_id"`
	UserID        string `gorm:"size:120;index" json:"user_id"`
	ClientName    string `gorm:"size:160;not null" json:"client_name"`
	RunningPrefix string `gorm:"size:30;not null;default:'IV'" json:"running_prefix"`
	RunningNumber int64  `gorm:"not null;default:0" json:"running_number"`
	Status        string `gorm:"size:40;not null;default:'ACTIVE';index" json:"status"`
	IsActive      bool   `gorm:"not null;default:true;index" json:"is_active"`
}
type POSProduct struct {
	Base
	SKUID       string `gorm:"size:120;not null;index" json:"sku_id"`
	SKUCode     string `gorm:"size:120;index" json:"sku_code"`
	Barcode     string `gorm:"size:120;index" json:"barcode"`
	ProductName string `gorm:"size:255;not null" json:"product_name"`
	Unit        string `gorm:"size:40;not null;default:'unit'" json:"unit"`
	IsActive    bool   `gorm:"not null;default:true;index" json:"is_active"`
}
type Price struct {
	Base
	PriceID       string     `gorm:"size:80;not null;uniqueIndex:idx_price_code" json:"price_id"`
	SKUID         string     `gorm:"size:120;not null;index" json:"sku_id"`
	SKUCode       string     `gorm:"size:120;index" json:"sku_code"`
	PriceTier     string     `gorm:"size:80;not null;default:'BASIC'" json:"price_tier"`
	Currency      string     `gorm:"size:3;not null;default:'THB'" json:"currency"`
	UnitPrice     float64    `gorm:"type:decimal(10,4);not null;default:0" json:"unit_price"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	IsActive      bool       `gorm:"not null;default:true;index" json:"is_active"`
}
type ShiftSession struct {
	Base
	ShiftNumber           string     `gorm:"size:80;not null;uniqueIndex" json:"shift_number"`
	POSClientID           string     `gorm:"size:80;not null;index" json:"posclient_id"`
	OpenedBy              string     `gorm:"size:120;not null" json:"opened_by"`
	ClosedBy              string     `gorm:"size:120" json:"closed_by"`
	OpeningCash           float64    `gorm:"type:decimal(10,4);not null;default:0" json:"opening_cash"`
	ClosingCashExpected   float64    `gorm:"type:decimal(10,4);not null;default:0" json:"closing_cash_expected"`
	ClosingCashActual     float64    `gorm:"type:decimal(10,4);not null;default:0" json:"closing_cash_actual"`
	ClosingCashDifference float64    `gorm:"type:decimal(10,4);not null;default:0" json:"closing_cash_difference"`
	Status                string     `gorm:"size:40;not null;default:'OPEN';index" json:"status"`
	OpenedAt              time.Time  `gorm:"not null" json:"opened_at"`
	ClosedAt              *time.Time `json:"closed_at"`
}
type SaleOrder struct {
	Base
	ReceiptNumber  string     `gorm:"size:80;not null;uniqueIndex" json:"receipt_number"`
	POSClientID    string     `gorm:"size:80;not null;index" json:"posclient_id"`
	ShiftSessionID string     `gorm:"type:uuid;index" json:"shift_session_id"`
	CustomerID     string     `gorm:"size:120" json:"customer_id"`
	CashierID      string     `gorm:"size:120;not null" json:"cashier_id"`
	Subtotal       float64    `gorm:"type:decimal(10,4);not null;default:0" json:"subtotal"`
	Discount       float64    `gorm:"type:decimal(10,4);not null;default:0" json:"discount"`
	VATRate        float64    `gorm:"type:decimal(10,4);not null;default:7" json:"vat_rate"`
	VAT            float64    `gorm:"type:decimal(10,4);not null;default:0" json:"vat"`
	GrandTotal     float64    `gorm:"type:decimal(10,4);not null;default:0" json:"grand_total"`
	Status         string     `gorm:"size:40;not null;default:'SUCCESS';index" json:"status"`
	Items          []SaleItem `gorm:"foreignKey:SaleOrderID" json:"items,omitempty"`
	Payments       []Payment  `gorm:"foreignKey:SaleOrderID" json:"payments,omitempty"`
}
type SaleItem struct {
	Base
	SaleOrderID string  `gorm:"type:uuid;not null;index" json:"sale_order_id"`
	SKUID       string  `gorm:"size:120;not null;index" json:"sku_id"`
	SKUCode     string  `gorm:"size:120;index" json:"sku_code"`
	ItemName    string  `gorm:"size:255;not null" json:"item_name"`
	Quantity    float64 `gorm:"type:decimal(10,4);not null;default:0" json:"quantity"`
	UnitPrice   float64 `gorm:"type:decimal(10,4);not null;default:0" json:"unit_price"`
	Discount    float64 `gorm:"type:decimal(10,4);not null;default:0" json:"discount"`
	TotalPrice  float64 `gorm:"type:decimal(10,4);not null;default:0" json:"total_price"`
}
type PaymentMethod struct {
	Base
	MethodCode      string `gorm:"size:80;not null;index" json:"method_code"`
	MethodName      string `gorm:"size:160;not null" json:"method_name"`
	MethodType      string `gorm:"size:40;not null;default:'CASH'" json:"method_type"`
	PromptPayTarget string `gorm:"size:120" json:"promptpay_target"`
	IsActive        bool   `gorm:"not null;default:true;index" json:"is_active"`
}
type Payment struct {
	Base
	PaymentNumber   string     `gorm:"size:80;not null;uniqueIndex" json:"payment_number"`
	SaleOrderID     string     `gorm:"type:uuid;not null;index" json:"sale_order_id"`
	PaymentMethodID string     `gorm:"type:uuid;not null;index" json:"payment_method_id"`
	Amount          float64    `gorm:"type:decimal(10,4);not null;default:0" json:"amount"`
	Status          string     `gorm:"size:40;not null;default:'SUCCESS';index" json:"status"`
	PaidAt          *time.Time `json:"paid_at"`
	EvidenceURL     string     `gorm:"size:500" json:"evidence_url"`
}
