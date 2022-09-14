package auth

import (
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	oldJwt "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	contextKey = "token"
)

type JwtHandler struct {
	expireIn time.Duration
	secret   []byte
	method   jwt.SigningMethod
	config   middleware.JWTConfig
}

func NewJwtHandler(cfg Config) *JwtHandler {
	return &JwtHandler{
		expireIn: cfg.ExpireAfter,
		secret:   []byte(cfg.SigningKey),
		method:   jwt.SigningMethodHS256,
		config: middleware.JWTConfig{
			SigningKey:    []byte(cfg.SigningKey),
			Claims:        &UserClaims{},
			SigningMethod: jwt.SigningMethodHS256.Name,
			ContextKey:    contextKey,
		},
	}
}

func (h *JwtHandler) GenerateFromUser(user *model.User) (string, error) {
	token := jwt.NewWithClaims(h.method, UserClaims{
		UserId: &user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.expireIn)),
		},
	})

	return token.SignedString(h.secret)
}

func (h *JwtHandler) ParseToUserClaims(c echo.Context) *UserClaims {
	token := c.Get(contextKey).(*oldJwt.Token)
	return token.Claims.(*UserClaims)
}

func (h *JwtHandler) Config() middleware.JWTConfig {
	return h.config
}
