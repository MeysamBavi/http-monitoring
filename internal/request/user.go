package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(3, 50), is.Alphanumeric),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 50), is.ASCII))
}
