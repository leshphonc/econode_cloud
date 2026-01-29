package event

import (
	"context"
	"econode-cloud/internal/model"

	"gorm.io/gorm"
)

type Repo interface {
	WithDB(db *gorm.DB) Repo
	Create(ctx context.Context, event *model.Event) error
}
