package auth

import (
	"fmt"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserId *model.ID `json:"user_id"`
	jwt.RegisteredClaims
}

func (c UserClaims) Valid() error {
	err := c.RegisteredClaims.Valid()
	if c.UserId == nil {
		return fmt.Errorf("user id is required. %w", err)
	}

	return err
}
