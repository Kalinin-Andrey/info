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
	"strconv"
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
	currency_sql_GetImportMaxTimeForUpdate = "SELECT currency_id, price_and_cap, concentration FROM cmc.import_max_time WHERE currency_id = ANY($1) FOR UPDATE;"
	currency_sql_MGet                      = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency WHERE id = any($1);"
	currency_sql_MGetBySlug                = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency WHERE slug = any($1);"
	currency_sql_GetAll                    = "SELECT id, symbol, slug, name, is_for_observing FROM cmc.currency WHERE is_for_observing = TRUE;"
	currency_sql_Create                    = "INSERT INTO cmc.currency(id, symbol, slug, name, is_for_observing) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id;"
	currency_sql_Update                    = "UPDATE cmc.currency SET symbol = $2, slug = $3, name = $4, is_for_observing = $5 WHERE id = $1;"
	currency_sql_Delete                    = "DELETE FROM cmc.currency WHERE id = $1;"

	import_max_time_sql_MCreate                    = "INSERT INTO cmc.import_max_time(currency_id, price_and_cap, concentration) VALUES "
	import_max_time_sql_MCreate_OnConflictDoUpdate = " ON CONFLICT (currency_id) DO UPDATE SET price_and_cap = EXCLUDED.price_and_cap, concentration = EXCLUDED.concentration;"
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

func (r *CurrencyRepository) GetImportMaxTimeForUpdateTx(ctx context.Context, tx domain.Tx, currencyIDs *[]uint) (map[uint]currency.ImportMaxTime, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetImportMaxTimeForUpdateTx"

	var entity currency.ImportMaxTime
	res := make(map[uint]currency.ImportMaxTime, len(*currencyIDs))

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
		res[entity.CurrencyID] = entity
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return res, nil
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

func (r *CurrencyRepository) GetAll(ctx context.Context) (*currency.CurrencyList, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "CurrencyRepository.GetAll"

	var entity currency.Currency
	res := make(currency.CurrencyList, 0, defaultCapacityForResult)

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

func (r CurrencyRepository) MCreateImportMaxTime(ctx context.Context, entities *[]currency.ImportMaxTime) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.MCreateImportMaxTime"
	const fields_nb = 3 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(import_max_time_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ")")
		params = append(params, entity.CurrencyID, entity.PriceAndCap, entity.Concentration)
	}
	b.WriteString(sql_OnConflictDoNothing)
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, b.String(), params...)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, b.String(), err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r CurrencyRepository) MUpsertImportMaxTimeMapTx(ctx context.Context, tx domain.Tx, entities map[uint]currency.ImportMaxTime) error {
	list := make([]currency.ImportMaxTime, 0, len(entities))
	var item currency.ImportMaxTime
	for _, item = range entities {
		list = append(list, item)
	}
	return r.MUpsertImportMaxTimeTx(ctx, tx, &list)
}

func (r CurrencyRepository) MUpsertImportMaxTimeTx(ctx context.Context, tx domain.Tx, entities *[]currency.ImportMaxTime) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	const metricName = "PriceAndCapRepository.MUpsertImportMaxTime"
	const fields_nb = 3 // при изменении количества полей нужно изменить MUpsertNmDimensions_Limit, чтобы, в результате, кол-во пар-ов не превышало 65т
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*fields_nb)
	b.WriteString(import_max_time_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*fields_nb+1) + ", $" + strconv.Itoa(i*fields_nb+2) + ", $" + strconv.Itoa(i*fields_nb+3) + ")")
		params = append(params, entity.CurrencyID, entity.PriceAndCap, entity.Concentration)
	}
	b.WriteString(import_max_time_sql_MCreate_OnConflictDoUpdate)
	start := time.Now().UTC()

	_, err := tx.Exec(ctx, b.String(), params...)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] %s query error; query: %s; error: %w", apperror.ErrInternal, metricName, b.String(), err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}
