package model

import "time"

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
type PageResponse[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
type ListQuery struct {
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	Search      string `form:"search"`
	Status      string `form:"status"`
	BranchID    string `form:"branch_id"`
	POSClientID string `form:"posclient_id"`
}
type AuditFields struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
type AuthContext struct {
	UserID     string   `json:"user_id"`
	MerchantID string   `json:"merchant_id"`
	BranchID   string   `json:"branch_id"`
	ClientID   string   `json:"client_id"`
	Scopes     []string `json:"scopes"`
}
