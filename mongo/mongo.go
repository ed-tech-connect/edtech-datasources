package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db      *mongo.Database
	session mongo.Session
}

func NewMongoRepository(db *mongo.Database) IRepository {
	return &MongoRepository{db: db}
}

func (r *MongoRepository) BeginTransaction(ctx context.Context) (IUnitOfWork, error) {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return &MongoUnitOfWork{session: session, db: r.db}, nil
}

func (r *MongoRepository) FindOne(ctx context.Context, collectionName string, qb *MongoQueryBuilder, result interface{}) error {
	projection := qb.BuildProjection()
	findOptions := qb.BuildFindOneOptions()

	err := r.db.Collection(collectionName).FindOne(ctx, qb.filter, findOptions.SetProjection(projection)).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return fmt.Errorf("error finding one record: %w", err)
	}
	return nil
}

func (r *MongoRepository) FindMany(ctx context.Context, collectionName string, qb *MongoQueryBuilder, results interface{}) (int, error) {
	totalCount, err := r.db.Collection(collectionName).CountDocuments(ctx, qb.filter)
	if err != nil {
		return 0, fmt.Errorf("error fetching total count: %w", err)
	}

	projection := qb.BuildProjection()
	findOptions := qb.BuildFindOptions()

	cursor, err := r.db.Collection(collectionName).Find(ctx, qb.filter, findOptions.SetProjection(projection))
	if err != nil {
		return 0, fmt.Errorf("error finding many records: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return 0, fmt.Errorf("error decoding many records: %w", err)
	}
	return int(totalCount), nil
}

func (r *MongoRepository) UpdateOne(ctx context.Context, collectionName string, qb *MongoQueryBuilder) (map[string]interface{}, error) {
	if qb.update == nil {
		return nil, fmt.Errorf("update document must be specified")
	}

	result, err := r.db.Collection(collectionName).UpdateOne(ctx, qb.filter, bson.M{"$set": qb.update})
	if err != nil {
		return nil, fmt.Errorf("error updating record: %w", err)
	}

	response := map[string]interface{}{
		"acknowledged":  true,
		"matchedCount":  result.MatchedCount,
		"modifiedCount": result.ModifiedCount,
	}
	return response, nil
}

func (r *MongoRepository) InsertOne(ctx context.Context, collectionName string, document bson.M) (map[string]interface{}, error) {
	result, err := r.db.Collection(collectionName).InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("error inserting record: %w", err)
	}

	response := map[string]interface{}{
		"acknowledged": true,
		"insertedId":   result.InsertedID,
	}

	return response, nil
}

func (r *MongoRepository) DeleteOne(ctx context.Context, collectionName string, qb *MongoQueryBuilder) (map[string]interface{}, error) {
	result, err := r.db.Collection(collectionName).DeleteOne(ctx, qb.filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting record: %w", err)
	}

	response := map[string]interface{}{
		"acknowledged": true,
		"deletedCount": result.DeletedCount,
	}

	return response, nil
}
