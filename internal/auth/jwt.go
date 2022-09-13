package auth

import (
	"strconv"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

type JwtHandler struct {
	expireIn time.Duration
	secret   []byte
	method   jwt.SigningMethod
	// middleware.JWTConfig
}

func NewJwtHandler(cfg Config) *JwtHandler {
	return &JwtHandler{
		expireIn: cfg.ExpireAfter,
		secret:   []byte(cfg.SigningKey),
		method:   jwt.SigningMethodHS256,
	}
}

func (h *JwtHandler) GenerateFromUser(user *model.User) (string, error) {
	token := jwt.NewWithClaims(h.method, jwt.RegisteredClaims{
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.expireIn)),
		ID:        strconv.FormatUint(uint64(user.Id), 10),
	})

	return token.SignedString(h.secret)
}
