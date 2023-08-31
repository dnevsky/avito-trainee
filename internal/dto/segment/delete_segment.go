package segment

import (
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type DeleteSegmentDTO struct {
	Name string `json:"name" form:"name" validate:"required,min=3,max=255" conform:"trim"`
}

func (dto *DeleteSegmentDTO) Validate() error {
	validate := validator.New()
	if err := conform.Strings(dto); err != nil {
		return err
	}
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}
