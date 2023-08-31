package segment

import (
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type CreateSegmentDTO struct {
	Name              string  `json:"name" form:"name" validate:"required,min=3,max=255" conform:"trim"`
	TTL               *string `json:"ttl" form:"ttl"`
	AutoAttachPercent *int    `json:"auto_attach_percent" form:"auto_attach_percent"`
}

func (dto *CreateSegmentDTO) Validate() error {
	validate := validator.New()
	if err := conform.Strings(dto); err != nil {
		return err
	}
	if err := validate.Struct(dto); err != nil {
		return err
	}
	return nil
}
