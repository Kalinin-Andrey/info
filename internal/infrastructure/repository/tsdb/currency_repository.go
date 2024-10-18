package tsdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"info/internal/domain"
	"info/internal/domain/currency"
	"info/internal/pkg/apperror"
	"strings"
	"time"
)

type CurrencyRepository struct {
	*Repository
}

const ()

var _ currency.WriteRepository = (*CurrencyRepository)(nil)
var _ currency.ReadRepository = (*CurrencyRepository)(nil)

func NewCurrencyRepository(repository *Repository) *CurrencyRepository {
	return &CurrencyRepository{
		Repository: repository,
	}
}

const (
	currency_sql_Get                       = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency WHERE id = $1;"
	currency_sql_GetBySlug                 = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency WHERE slug = $1;"
	currency_sql_GetImportMaxTime          = "SELECT currency_id, price_and_cap, concentration FROM cmc.import_max_time WHERE currency_id = ANY($1);"
	currency_sql_GetImportMaxTimeForUpdate = "SELECT currency_id, price_and_cap, concentration FROM cmc.import_max_time WHERE currency_id = ANY($1) FOR UPDATE;"
	currency_sql_MGet                      = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency FROM blog.blog WHERE id = any($1);"
	currency_sql_MGetBySlug                = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency FROM blog.blog WHERE slug = any($1);"
	currency_sql_GetAll                    = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency;"
	currency_sql_Create                    = "INSERT INTO cmc.currency(id, symbol, slug, name, is_for_observing) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id;"
	currency_sql_Update                    = "UPDATE cmc.currency SET symbol = $2, slug = $3, name = $4, is_for_observing = $5 WHERE id = $1;"
	currency_sql_Delete                    = "DELETE FROM cmc.currency WHERE id = $1;"
)

func (r *CurrencyRepository) Get(ctx context.Context, ID uint) (*currency.Currency, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Get"
	start := time.Now().UTC()

	entity := &currency.Currency{}
	if err := r.db.QueryRow(ctx, currency_sql_Get, ID).Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Get, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *CurrencyRepository) GetBySlug(ctx context.Context, slug string) (*currency.Currency, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetBySlug"
	start := time.Now().UTC()

	entity := &currency.Currency{}
	if err := r.db.QueryRow(ctx, currency_sql_GetBySlug, slug).Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetBySlug, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *CurrencyRepository) GetImportMaxTime(ctx context.Context, currencyIDs *[]uint) (*[]currency.ImportMaxTime, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetImportMaxTime"

	var entity currency.ImportMaxTime
	res := make([]currency.ImportMaxTime, 0, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_GetImportMaxTime, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTime, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.PriceAndCap, &entity.Concentration); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTime, err)
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

func (r *CurrencyRepository) GetImportMaxTimeForUpdate(ctx context.Context, tx domain.Tx, currencyIDs *[]uint) (*[]currency.ImportMaxTime, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetImportMaxTimeForUpdate"

	var entity currency.ImportMaxTime
	res := make([]currency.ImportMaxTime, 0, len(*currencyIDs))

	start := time.Now().UTC()
	rows, err := tx.Query(ctx, currency_sql_GetImportMaxTimeForUpdate, *currencyIDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTimeForUpdate, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.CurrencyID, &entity.PriceAndCap, &entity.Concentration); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetImportMaxTimeForUpdate, err)
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

func (r *CurrencyRepository) MGet(ctx context.Context, IDs *[]uint) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.MGet"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, len(*IDs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_MGet, *IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGet, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGet, err)
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

func (r *CurrencyRepository) MGetBySlug(ctx context.Context, slugs *[]string) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.MGetBySlug"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, len(*slugs))

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_MGetBySlug, *slugs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetBySlug, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_MGetBySlug, err)
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

func (r *CurrencyRepository) GetAll(ctx context.Context) (*[]currency.Currency, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetAll"

	var entity currency.Currency
	res := make([]currency.Currency, 0, defaultCapacityForResult)

	start := time.Now().UTC()
	rows, err := r.db.Query(ctx, currency_sql_GetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetAll, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entity.ID, &entity.Symbol, &entity.Slug, &entity.Name, &entity.IsForObserving); err != nil {
			r.metrics.SqlMetrics.Inc(metricName, metricsFail)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
			return nil, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_GetAll, err)
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

func (r *CurrencyRepository) Create(ctx context.Context, entity *currency.Currency) (ID uint, err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "CurrencyRepository.Upsert"
	start := time.Now().UTC()

	if err := r.db.QueryRow(ctx, currency_sql_Create, entity.ID, entity.Symbol, entity.Slug, entity.Name, entity.IsForObserving).Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return 0, apperror.ErrNotFound
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return 0, fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Create, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return ID, nil
}

func (r *CurrencyRepository) Update(ctx context.Context, entity *currency.Currency) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Update"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, currency_sql_Update, entity.ID, entity.Symbol, entity.Slug, entity.Name, entity.IsForObserving)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Update, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *CurrencyRepository) Delete(ctx context.Context, ID uint) error {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.Delete"
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, currency_sql_Delete, ID)
	if err != nil {
		if strings.Contains(err.Error(), errMsg_duplicateKey) {
			r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
			r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
			return apperror.ErrBadRequest
		}
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, currency_sql_Delete, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
