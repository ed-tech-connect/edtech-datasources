package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ed-tech-connect/edtech-datasources/sqlscan"
)

type MySQLRepository struct {
	db *sql.DB
	Tx *sql.Tx
}

func NewMySQLRepository(db *sql.DB) IRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) BeginTransaction(ctx context.Context) (IUnitOfWork, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &MySQLUnitOfWork{Tx: tx}, nil
}

func (r *MySQLRepository) getExecutor() (queryExecutor, error) {
	if r.Tx != nil {
		return r.Tx, nil
	}
	return r.db, nil
}

func (r *MySQLRepository) FindOne(ctx context.Context, tableName string, builder *QueryBuilder, result interface{}) error {
	query, args := builder.BuildSelectQuery(tableName)
	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	row, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error finding a record: %w", err)
	}
	if err := sqlscan.Row(result, row); err != nil {
		if strings.EqualFold(err.Error(), "sql: no rows in result set") {
			return nil
		}
		return fmt.Errorf("failed to scan row: %w", err)
	}
	return nil
}

func (r *MySQLRepository) FindMany(ctx context.Context, tableName string, builder *QueryBuilder, results interface{}) (int, error) {
	countQuery, countArgs := builder.BuildCountQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return 0, err
	}

	var totalCount int
	err = executor.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("error fetching total count: %w", err)
	}

	query, args := builder.BuildSelectManyQuery(tableName)
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("error finding many records: %w", err)
	}
	defer rows.Close()
	if err := sqlscan.Rows(results, rows); err != nil {
		return 0, fmt.Errorf("failed to scan rows: %w", err)
	}
	return totalCount, nil
}

func (r *MySQLRepository) UpdateOne(ctx context.Context, tableName string, qb *QueryBuilder) error {
	query, args := qb.BuildUpdateQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}
	return nil
}

func (r *MySQLRepository) UpdateMany(ctx context.Context, tableName string, qb *QueryBuilder) error {
	query, args := qb.BuildUpdateManyQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating records: %w", err)
	}
	return nil
}

func (r *MySQLRepository) InsertOne(ctx context.Context, tableName string, qb *QueryBuilder) error {
	query, args := qb.BuildInsertQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	return nil
}

func (r *MySQLRepository) DeleteOne(ctx context.Context, tableName string, qb *QueryBuilder) error {
	query, args := qb.BuildDeleteQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}
	return nil
}

func (r *MySQLRepository) DeleteMany(ctx context.Context, tableName string, qb *QueryBuilder) error {
	query, args := qb.BuildDeleteQuery(tableName)

	executor, err := r.getExecutor()
	if err != nil {
		return err
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting records: %w", err)
	}
	return nil
}
