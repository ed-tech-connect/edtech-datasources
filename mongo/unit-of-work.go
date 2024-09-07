package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	session mongo.Session
	db      *mongo.Database
}

func (uow *MongoUnitOfWork) GetRepository() IRepository {
	return &MongoRepository{db: uow.db, session: uow.session}
}

func (uow *MongoUnitOfWork) Commit() error {
	ctx := mongo.NewSessionContext(context.Background(), uow.session)
	err := uow.session.CommitTransaction(ctx)
	uow.session.EndSession(ctx)
	return err
}

func (uow *MongoUnitOfWork) Rollback() error {
	ctx := mongo.NewSessionContext(context.Background(), uow.session)
	err := uow.session.AbortTransaction(ctx)
	uow.session.EndSession(ctx)
	return err
}
