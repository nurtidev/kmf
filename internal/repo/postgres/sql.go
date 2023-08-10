package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurtidev/kmf/internal/model"
	"time"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) Save(ctx context.Context, currency *model.Currency) error {
	query := `
		INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.DB.Exec(ctx, query, currency.Title, currency.Code, currency.Value, currency.Date)
	return err
}

func (r *Repository) FindByDateAndCode(ctx context.Context, date time.Time, code string) ([]*model.Currency, error) {
	var currencies []*model.Currency
	query := `
		SELECT * FROM R_CURRENCY
		WHERE A_DATE = $1
	`
	args := []interface{}{date}
	if code != "" {
		query += " AND CODE = $2"
		args = append(args, code)
	}
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var currency model.Currency
		err = rows.Scan(&currency.ID, &currency.Title, &currency.Code, &currency.Value, &currency.Date)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, &currency)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}
