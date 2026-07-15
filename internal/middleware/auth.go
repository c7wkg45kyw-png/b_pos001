package middleware

import (
	"bpos001/backend/internal/config"
	"bpos001/backend/internal/model"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const AuthContextKey = "auth_context"

type Claims struct {
	MerchantID string   `json:"merchant_id"`
	BranchID   string   `json:"branch_id"`
	ClientID   string   `json:"client_id"`
	Type       string   `json:"type"`
	Scopes     []string `json:"scopes"`
	jwt.RegisteredClaims
}

func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			abort(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token")
			return
		}
		tokenText := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenText, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		}, jwt.WithIssuer(cfg.JWTIssuer))
		if err != nil || !token.Valid {
			abort(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			return
		}
		merchantID := claims.MerchantID
		if merchantID == "" {
			merchantID = c.GetHeader("X-Merchant-ID")
		}
		if merchantID == "" {
			abort(c, http.StatusUnauthorized, "UNAUTHORIZED", "merchant_id is required")
			return
		}
		branchID := claims.BranchID
		if branchID == "" {
			branchID = c.GetHeader("X-Branch-ID")
		}
		c.Set(AuthContextKey, model.AuthContext{UserID: claims.Subject, MerchantID: merchantID, BranchID: branchID, ClientID: claims.ClientID, Scopes: claims.Scopes})
		c.Next()
	}
}
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, _ := c.Get(AuthContextKey)
		ctx, ok := v.(model.AuthContext)
		if !ok {
			abort(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing auth context")
			return
		}
		for _, current := range ctx.Scopes {
			if current == scope || current == "pos:*" {
				c.Next()
				return
			}
		}
		abort(c, http.StatusForbidden, "FORBIDDEN", "missing scope: "+scope)
	}
}
func CurrentAuth(c *gin.Context) model.AuthContext {
	value, _ := c.Get(AuthContextKey)
	ctx, _ := value.(model.AuthContext)
	return ctx
}
func abort(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, model.ErrorResponse{Success: false, Code: code, Message: message})
}
