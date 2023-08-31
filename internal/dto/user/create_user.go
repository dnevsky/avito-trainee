package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type CreateUserDTO struct {
	Segments []string `json:"segments" form:"segments"`
}

func (dto *CreateUserDTO) Validate() error {
	validate := validator.New()
	if err := conform.Strings(dto); err != nil {
		return err
	}
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}
