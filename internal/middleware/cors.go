package middleware

import (
	"bpos001/backend/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CORS(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := ""
		for _, item := range cfg.CORSAllowOrigins {
			if item == "*" || item == origin {
				allowed = item
				if item == "*" && origin != "" {
					allowed = origin
				}
				break
			}
		}
		if allowed == "" && cfg.AppEnv != "production" && isLocalOrigin(origin) {
			allowed = origin
		}
		if allowed != "" {
			c.Header("Access-Control-Allow-Origin", allowed)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept-Language, X-Merchant-ID, X-Branch-ID")
		c.Header("Access-Control-Allow-Methods", strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, ", "))
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
func isLocalOrigin(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") || strings.HasPrefix(origin, "http://[::1]:")
}
