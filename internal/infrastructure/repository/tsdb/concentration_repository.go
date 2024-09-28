package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"info/internal/domain/concentration"
	"info/internal/pkg/apperror"
	"strings"
	"time"
)

type ConcentrationRepository struct {
	*Repository
}

const ()

var _ concentration.WriteRepository = (*ConcentrationRepository)(nil)
var _ concentration.ReadRepository = (*ConcentrationRepository)(nil)

func NewConcentrationRepository(repository *Repository) *ConcentrationRepository {
	return &ConcentrationRepository{
		Repository: repository,
	}
}

const (
	concentration_sql_Get    = "SELECT currency_id, whales, investors, retail, others, d FROM cmc.concentration WHERE currency_id = $1;"
	concentration_sql_MGet   = "SELECT currency_id, whales, investors, retail, others, d FROM cmc.concentration FROM blog.blog WHERE currency_id = any($1);"
	concentration_sql_GetAll = "SELECT currency_id, whales, investors, retail, others, d FROM cmc.concentration;"
	concentration_sql_Create = "INSERT INTO cmc.concentration(currency_id, whales, investors, retail, others, d) VALUES ($1, $2, $3, $4, $5, $6) RETURNING currency_id;"
	concentration_sql_Update = "UPDATE cmc.concentration SET whales = $2, investors = $3, retail = $4, others = $5, d = $6 WHERE currency_id = $1;"
	concentration_sql_Delete = "DELETE FROM cmc.concentration WHERE currency_id = $1;"
)

func (r *ConcentrationRepository) Get(ctx context.Context, currencyID uint) (*concentration.Concentration, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.Get"
	start := time.Now().UTC()

	entity := &concentration.Concentration{}
	if err := r.db.QueryRow(ctx, concentration_sql_Get, currencyID).Scan(&entity.CurrencyID, &entity.Whales, &entity.Investors, &entity.Retail, &entity.Others, &entity.D); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *ConcentrationRepository) MGet(ctx context.Context, currencyIDs *[]uint) (*[]concentration.Concentration, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.MGet"

	var entity concentration.Concentration
	res := make([]concentration.Concentration, 0, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, concentration_sql_MGet, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Whales, &entity.Investors, &entity.Retail, &entity.Others, &entity.D); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_MGet, err)
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

func (r *ConcentrationRepository) GetAll(ctx context.Context) (*[]concentration.Concentration, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.GetAll"

	var entity concentration.Concentration
	res := make([]concentration.Concentration, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, concentration_sql_GetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_GetAll, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.Whales, &entity.Investors, &entity.Retail, &entity.Others, &entity.D); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_GetAll, err)
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

func (r *ConcentrationRepository) Create(ctx context.Context, entity *concentration.Concentration) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "ConcentrationRepository.Create"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, concentration_sql_Create, entity.CurrencyID, entity.Whales, entity.Investors, entity.Retail, entity.Others, entity.D).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *ConcentrationRepository) Update(ctx context.Context, entity *concentration.Concentration) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, concentration_sql_Update, entity.CurrencyID, entity.Whales, entity.Investors, entity.Retail, entity.Others, entity.D)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *ConcentrationRepository) Delete(ctx context.Context, CurrencyID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "ConcentrationRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, concentration_sql_Delete, CurrencyID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, concentration_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
