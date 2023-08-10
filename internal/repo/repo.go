package repo

import (
	"context"
	"github.com/nurtidev/kmf/internal/model"
	"time"
)

type Repository interface {
	Save(ctx context.Context, currency *model.Currency) error
	FindByDateAndCode(ctx context.Context, date time.Time, code string) ([]*model.Currency, error)
}
