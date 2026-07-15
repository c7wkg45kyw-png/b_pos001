package handler

import (
	"encoding/json"

	"bpos001/backend/internal/middleware"
	"bpos001/backend/internal/model"
	"bpos001/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type POSHandler struct{ usecase *usecase.Usecase }

func NewPOSHandler(u *usecase.Usecase) *POSHandler { return &POSHandler{usecase: u} }
func (h *POSHandler) List(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q model.ListQuery
		_ = c.ShouldBindQuery(&q)
		result, err := h.usecase.List(middleware.CurrentAuth(c), resource, q)
		if err != nil {
			handleError(c, err)
			return
		}
		ok(c, resource, result)
	}
}
func (h *POSHandler) Get(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.usecase.Get(middleware.CurrentAuth(c), resource, c.Param("id_or_code"))
		if err != nil {
			handleError(c, err)
			return
		}
		ok(c, resource, result)
	}
}
func (h *POSHandler) Create(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := bindResourcePayload(c, resource)
		if err != nil {
			fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		result, err := h.usecase.Create(middleware.CurrentAuth(c), resource, payload)
		if err != nil {
			handleError(c, err)
			return
		}
		created(c, resource+" created", result)
	}
}
func (h *POSHandler) Update(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := bindResourcePayload(c, resource)
		if err != nil {
			fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		result, err := h.usecase.Update(middleware.CurrentAuth(c), resource, c.Param("id"), payload)
		if err != nil {
			handleError(c, err)
			return
		}
		ok(c, resource+" updated", result)
	}
}
func (h *POSHandler) Delete(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.usecase.Delete(middleware.CurrentAuth(c), resource, c.Param("id")); err != nil {
			handleError(c, err)
			return
		}
		noContent(c)
	}
}
func (h *POSHandler) SelectPOS(c *gin.Context) {
	result, err := h.usecase.SelectPOS(middleware.CurrentAuth(c), c.Param("id_or_code"))
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, "pos selected", result)
}
func (h *POSHandler) ShiftDetails(c *gin.Context) {
	posClientID := c.Param("id_or_code")
	if posClientID == "" {
		posClientID = c.Param("posclient_id")
	}
	result, err := h.usecase.GetShiftDetails(middleware.CurrentAuth(c), posClientID)
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, "shift details", result)
}

func (h *POSHandler) CloseShift(c *gin.Context) {
	payload := map[string]any{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	result, err := h.usecase.CloseShift(middleware.CurrentAuth(c), c.Param("id_or_code"), payload)
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, "shift closed", result)
}

func (h *POSHandler) AdjustPrice(c *gin.Context) {
	var req model.PriceAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	raw, _ := json.Marshal(req)
	payload := map[string]any{}
	_ = json.Unmarshal(raw, &payload)
	result, err := h.usecase.AdjustPrice(middleware.CurrentAuth(c), payload)
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, "price adjusted", result)
}

func bindResourcePayload(c *gin.Context, resource string) (map[string]any, error) {
	var req any
	switch resource {
	case "posclients":
		req = &model.POSClientRequest{}
	case "products":
		req = &model.POSProductRequest{}
	case "prices":
		req = &model.PriceRequest{}
	case "shifts":
		req = &model.ShiftSessionRequest{}
	case "sale-orders":
		req = &model.SaleOrderRequest{}
	case "sale-items":
		req = &model.SaleItemRequest{}
	case "payment-methods":
		req = &model.PaymentMethodRequest{}
	case "payments":
		req = &model.PaymentRequest{}
	default:
		req = map[string]any{}
	}
	if err := c.ShouldBindJSON(req); err != nil {
		return nil, err
	}
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (h *POSHandler) Summary(c *gin.Context) {
	result, err := h.usecase.Summary(middleware.CurrentAuth(c))
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, "summary", result)
}
