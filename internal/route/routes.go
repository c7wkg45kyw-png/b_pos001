package route

import (
	"bpos001/backend/internal/config"
	"bpos001/backend/internal/handler"
	"bpos001/backend/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct{ POS *handler.POSHandler }

func Register(router *gin.Engine, cfg config.Config, h Handlers) {
	router.Use(middleware.CORS(cfg))
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"success": true, "message": "healthy"}) })
	router.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")
	router.GET("/docs/swagger.html", swaggerHTML)
	api := router.Group("/api/v1")
	api.Use(middleware.Auth(cfg))
	api.GET("/pos/select/:id_or_code", middleware.RequireScope("pos:read"), h.POS.SelectPOS)
	api.GET("/posclients/:id_or_code/shift-details", middleware.RequireScope("pos:read"), h.POS.ShiftDetails)
	for _, res := range []string{"posclients", "products", "prices", "shifts", "sale-orders", "sale-items", "payment-methods", "payments"} {
		registerCRUD(api, res, h.POS)
	}
}
func registerCRUD(api *gin.RouterGroup, resource string, h *handler.POSHandler) {
	base := "/" + resource
	api.GET(base, middleware.RequireScope("pos:read"), h.List(resource))
	api.POST(base, middleware.RequireScope("pos:create"), h.Create(resource))
	api.GET(base+"/:id_or_code", middleware.RequireScope("pos:read"), h.Get(resource))
	api.PUT(base+"/:id", middleware.RequireScope("pos:update"), h.Update(resource))
	api.PATCH(base+"/:id", middleware.RequireScope("pos:update"), h.Update(resource))
	api.DELETE(base+"/:id", middleware.RequireScope("pos:delete"), h.Delete(resource))
}
func swaggerHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html><html><head><title>BPOS001 Swagger</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>window.onload=()=>SwaggerUIBundle({url:'/docs/openapi.yaml',dom_id:'#swagger-ui'});</script></body></html>`)
}
