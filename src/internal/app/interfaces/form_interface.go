package interfaces

import (
	"context"

	database "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/generated"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/models"
)

type FormService interface {
	FindAll(ctx context.Context) ([]*database.FindAllFormsRow, error)
	FindAllView(ctx context.Context) ([]*database.FormView, error)
	Create(ctx context.Context, createFormRequest *models.CreateFormRequest) error
}
