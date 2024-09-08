package mysql

import (
	"context"
	"database/sql"
)

type IRepository interface {
	FindOne(context.Context, string, *QueryBuilder, interface{}) error
	FindMany(context.Context, string, *QueryBuilder, interface{}) (int, error)
	UpdateOne(context.Context, string, *QueryBuilder) (map[string]interface{}, error)
	UpdateMany(context.Context, string, *QueryBuilder) (map[string]interface{}, error)
	InsertOne(context.Context, string, *QueryBuilder) (map[string]interface{}, error)
	DeleteOne(context.Context, string, *QueryBuilder) (map[string]interface{}, error)
	DeleteMany(context.Context, string, *QueryBuilder) (map[string]interface{}, error)

	BeginTransaction(ctx context.Context) (IUnitOfWork, error)
}

type IUnitOfWork interface {
	Commit() error
	Rollback() error
	GetRepository() IRepository
}

type queryExecutor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
