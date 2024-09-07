package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type IRepository interface {
	FindOne(context.Context, string, *MongoQueryBuilder, interface{}) error
	FindMany(context.Context, string, *MongoQueryBuilder, interface{}) (int, error)
	UpdateOne(context.Context, string, *MongoQueryBuilder) error
	InsertOne(context.Context, string, bson.M) error
	DeleteOne(context.Context, string, *MongoQueryBuilder) error

	BeginTransaction(ctx context.Context) (IUnitOfWork, error)
}

type IUnitOfWork interface {
	Commit() error
	Rollback() error
	GetRepository() IRepository
}
