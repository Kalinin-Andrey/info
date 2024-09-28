package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"strings"
	"time"
)

type PriceAndCapRepository struct {
	*Repository
}

const ()

var _ price_and_cap.WriteRepository = (*PriceAndCapRepository)(nil)
var _ price_and_cap.ReadRepository = (*PriceAndCapRepository)(nil)

func NewPriceAndCapRepository(repository *Repository) *PriceAndCapRepository {
	return &PriceAndCapRepository{
		Repository: repository,
	}
}

const (
	price_and_cap_sql_Get    = "SELECT currency_id, price, daily_volume, cap, ts FROM cmc.price_and_cap WHERE currency_id = $1;"
	price_and_cap_sql_MGet   = "SELECT currency_id, price, daily_volume, cap, ts FROM cmc.price_and_cap FROM blog.blog WHERE currency_id = any($1);"
	price_and_cap_sql_GetAll = "SELECT currency_id, price, daily_volume, cap, ts FROM cmc.price_and_cap;"
	price_and_cap_sql_Create = "INSERT INTO cmc.price_and_cap(currency_id, price, daily_volume, cap, ts) VALUES ($1, $2, $3, $4, $5) RETURNING currency_id;"
	price_and_cap_sql_Update = "UPDATE cmc.price_and_cap SET price = $2, daily_volume = $3, cap = $4, ts = $5 WHERE currency_id = $1;"
	price_and_cap_sql_Delete = "DELETE FROM cmc.price_and_cap WHERE currency_id = $1;"
)

func (r *PriceAndCapRepository) Get(ctx context.Context, currencyID uint) (*price_and_cap.PriceAndCap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.Get"
	start := time.Now().UTC()

	entity := &price_and_cap.PriceAndCap{}
	if err := r.db.QueryRow(ctx, price_and_cap_sql_Get, currencyID).Scan(&entity.CurrencyID, &entity.Price, &entity.DailyVolume, &entity.Cap, &entity.Ts); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *PriceAndCapRepository) MGet(ctx context.Context, currencyIDs *[]uint) (*[]price_and_cap.PriceAndCap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.MGet"

	var entity price_and_cap.PriceAndCap
	res := make([]price_and_cap.PriceAndCap, 0, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, price_and_cap_sql_MGet, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Price, &entity.DailyVolume, &entity.Cap, &entity.Ts); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_MGet, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *PriceAndCapRepository) GetAll(ctx context.Context) (*[]price_and_cap.PriceAndCap, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.GetAll"

	var entity price_and_cap.PriceAndCap
	res := make([]price_and_cap.PriceAndCap, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, price_and_cap_sql_GetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_GetAll, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Price, &entity.DailyVolume, &entity.Cap, &entity.Ts); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_GetAll, err)
		}
		res = append(res, entity)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (r *PriceAndCapRepository) Create(ctx context.Context, entity *price_and_cap.PriceAndCap) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.Create"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, price_and_cap_sql_Create, entity.CurrencyID, entity.Price, entity.DailyVolume, entity.Cap, entity.Ts).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *PriceAndCapRepository) Update(ctx context.Context, entity *price_and_cap.PriceAndCap) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, price_and_cap_sql_Update, entity.CurrencyID, entity.Price, entity.DailyVolume, entity.Cap, entity.Ts)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *PriceAndCapRepository) Delete(ctx context.Context, CurrencyID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "PriceAndCapRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, price_and_cap_sql_Delete, CurrencyID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, price_and_cap_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
