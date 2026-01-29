package event

import (
	"context"
	"econode-cloud/internal/model"

	"gorm.io/gorm"
)

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) Repo {
	return &eventRepo{db: db}
}

func (r *eventRepo) WithDB(db *gorm.DB) Repo {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(ctx context.Context, e *model.Event) error {
	return r.db.WithContext(ctx).Create(e).Error
}
