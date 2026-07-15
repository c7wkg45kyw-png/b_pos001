package model

import "time"

type POSClientRequest struct {
	POSClientID   string `json:"posclient_id" binding:"required"`
	UserID        string `json:"user_id"`
	ClientName    string `json:"client_name" binding:"required"`
	RunningPrefix string `json:"running_prefix"`
	Status        string `json:"status"`
	IsActive      *bool  `json:"is_active,omitempty"`
}

type POSProductRequest struct {
	SKUID       string `json:"sku_id" binding:"required"`
	SKUCode     string `json:"sku_code"`
	Barcode     string `json:"barcode"`
	ProductName string `json:"product_name" binding:"required"`
	Unit        string `json:"unit"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type PriceRequest struct {
	PriceID       string     `json:"price_id" binding:"required"`
	SKUID         string     `json:"sku_id" binding:"required"`
	SKUCode       string     `json:"sku_code"`
	PriceTier     string     `json:"price_tier"`
	Currency      string     `json:"currency"`
	UnitPrice     float64    `json:"unit_price" binding:"required"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type PriceAdjustRequest struct {
	SKUID     string  `json:"sku_id"`
	SKUCode   string  `json:"sku_code"`
	Barcode   string  `json:"barcode"`
	PriceTier string  `json:"price_tier"`
	Currency  string  `json:"currency"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
}

type ShiftSessionRequest struct {
	ShiftNumber           string  `json:"shift_number"`
	POSClientID           string  `json:"posclient_id" binding:"required"`
	OpenedBy              string  `json:"opened_by"`
	ClosedBy              string  `json:"closed_by"`
	OpeningCash           float64 `json:"opening_cash"`
	ClosingCashExpected   float64 `json:"closing_cash_expected"`
	ClosingCashActual     float64 `json:"closing_cash_actual"`
	ClosingCashDifference float64 `json:"closing_cash_difference"`
	Status                string  `json:"status"`
}

type SaleOrderRequest struct {
	ReceiptNumber  string  `json:"receipt_number"`
	POSClientID    string  `json:"posclient_id" binding:"required"`
	ShiftSessionID string  `json:"shift_session_id"`
	CustomerID     string  `json:"customer_id"`
	CashierID      string  `json:"cashier_id"`
	Subtotal       float64 `json:"subtotal"`
	Discount       float64 `json:"discount"`
	VATRate        float64 `json:"vat_rate"`
	Status         string  `json:"status"`
}

type SaleItemRequest struct {
	SaleOrderID string  `json:"sale_order_id" binding:"required"`
	SKUID       string  `json:"sku_id" binding:"required"`
	SKUCode     string  `json:"sku_code"`
	ItemName    string  `json:"item_name" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Discount    float64 `json:"discount"`
}

type PaymentMethodRequest struct {
	MethodCode      string `json:"method_code" binding:"required"`
	MethodName      string `json:"method_name" binding:"required"`
	MethodType      string `json:"method_type"`
	PromptPayTarget string `json:"promptpay_target"`
	IsActive        *bool  `json:"is_active,omitempty"`
}

type PaymentRequest struct {
	PaymentNumber   string     `json:"payment_number" binding:"required"`
	SaleOrderID     string     `json:"sale_order_id" binding:"required"`
	PaymentMethodID string     `json:"payment_method_id" binding:"required"`
	Amount          float64    `json:"amount" binding:"required"`
	Status          string     `json:"status"`
	PaidAt          *time.Time `json:"paid_at"`
	EvidenceURL     string     `json:"evidence_url"`
}
