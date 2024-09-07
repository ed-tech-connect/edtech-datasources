package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQueryBuilder struct {
	columns []string
	update  bson.M
	filter  bson.M
	limit   int64
	skip    int64
	sort    bson.D
}

func NewMongoQueryBuilder() *MongoQueryBuilder {
	return &MongoQueryBuilder{}
}

func (b *MongoQueryBuilder) Select(columns []string) *MongoQueryBuilder {
	b.columns = columns
	return b
}

func (b *MongoQueryBuilder) Where(filter bson.M) *MongoQueryBuilder {
	b.filter = filter
	return b
}

func (b *MongoQueryBuilder) Limit(limit int64) *MongoQueryBuilder {
	b.limit = limit
	return b
}

func (b *MongoQueryBuilder) Skip(skip int64) *MongoQueryBuilder {
	b.skip = skip
	return b
}

func (b *MongoQueryBuilder) Sort(sort bson.D) *MongoQueryBuilder {
	b.sort = sort
	return b
}

func (b *MongoQueryBuilder) BuildProjection() bson.M {
	if len(b.columns) == 0 {
		return nil
	}

	projection := bson.M{}
	for _, column := range b.columns {
		projection[column] = 1
	}
	return projection
}

func (b *MongoQueryBuilder) BuildFindOptions() *options.FindOptions {
	opts := options.Find()
	if b.limit > 0 {
		opts.SetLimit(b.limit)
	}
	if b.skip > 0 {
		opts.SetSkip(b.skip)
	}
	if len(b.sort) > 0 {
		opts.SetSort(b.sort)
	}
	return opts
}

func (b *MongoQueryBuilder) BuildFindOneOptions() *options.FindOneOptions {
	opts := options.FindOne()
	if len(b.sort) > 0 {
		opts.SetSort(b.sort)
	}
	return opts
}
