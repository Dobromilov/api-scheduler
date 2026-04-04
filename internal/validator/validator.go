package validator

import (
	"Scheduler-api/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func ValidateConfig(config *models.SimulationConfig) error {
	return validate.Struct(config)
}
