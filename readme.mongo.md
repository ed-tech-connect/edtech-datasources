# MongoDB Repository

This package provides an implementation of a MongoDB repository with support for common database operations such as querying, updating, inserting, and deleting documents. It uses a `MongoQueryBuilder` to build queries dynamically.

## Package Overview

- **`MongoRepository`**: Implements repository methods for interacting with MongoDB.
- **`MongoUnitOfWork`**: Manages transactions within MongoDB.
- **`MongoQueryBuilder`**: Helps construct MongoDB queries dynamically.

## Installation

To use this package, ensure you have the `mongo-driver` package installed:

```bash
go get go.mongodb.org/mongo-driver/mongo
```

```bash
import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/yourmodule/mongo" // Replace with your actual module path
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		// Handle error
	}

	db := client.Database("your_database_name")
	repo := mongo.NewMongoRepository(db)
}

func findOneExample(repo *mongo.MongoRepository) {
	ctx := context.Background()

	// Create a new query builder
	qb := mongo.NewMongoQueryBuilder().SetFilter(bson.M{"name": "John"}).SetProjection(bson.M{"_id": 0, "name": 1})

	var result bson.M
	err := repo.FindOne(ctx, "collection_name", qb, &result)
	if err != nil {
		// Handle error
	}
	fmt.Println(result)
}


func findManyExample(repo *mongo.MongoRepository) {
	ctx := context.Background()

	// Create a new query builder
	qb := mongo.NewMongoQueryBuilder().
		SetFilter(bson.M{"age": bson.M{"$gte": 18}}).
		SetProjection(bson.M{"_id": 0, "name": 1, "age": 1})

	var results []bson.M
	err := repo.FindMany(ctx, "collection_name", qb, &results)
	if err != nil {
		// Handle error
	}
	fmt.Println(results)
}

func updateOneExample(repo *mongo.MongoRepository) {
	ctx := context.Background()

	// Create a new query builder
	qb := mongo.NewMongoQueryBuilder().
		SetFilter(bson.M{"name": "John Doe"}).
		SetUpdate(bson.M{"age": 30})

	err := repo.UpdateOne(ctx, "collection_name", qb)
	if err != nil {
		// Handle error
	}
}

func insertOneExample(repo *mongo.MongoRepository) {
	ctx := context.Background()

	document := bson.M{"name": "Jane Doe", "age": 25}
	err := repo.InsertOne(ctx, "collection_name", document)
	if err != nil {
		// Handle error
	}
}

func deleteOneExample(repo *mongo.MongoRepository) {
	ctx := context.Background()

	// Create a new query builder
	qb := mongo.NewMongoQueryBuilder().
		SetFilter(bson.M{"name": "Jane Doe"})

	err := repo.DeleteOne(ctx, "collection_name", qb)
	if err != nil {
		// Handle error
	}
}

```

### Transactional

```bash
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/yourmodule/mongo" // Replace with your actual module path
)

func main() {
	// Setup MongoDB client and repository
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	db := client.Database("your_database_name")
	repo := mongo.NewMongoRepository(db)

	// Begin a transaction
	ctx := context.Background()
	uow, err := repo.BeginTransaction(ctx)
	if err != nil {
		fmt.Println("Failed to begin transaction:", err)
		return
	}

	// Perform operations within the transaction
	err = performTransactionalOperations(ctx, repo, uow)
	if err != nil {
		// Abort the transaction in case of error
		uow.session.AbortTransaction(ctx)
		fmt.Println("Transaction aborted due to error:", err)
		return
	}

	// Commit the transaction if all operations are successful
	if err := uow.session.CommitTransaction(ctx); err != nil {
		fmt.Println("Failed to commit transaction:", err)
		return
	}

	fmt.Println("Transaction committed successfully")
}

// performTransactionalOperations performs database operations within a transaction
func performTransactionalOperations(ctx context.Context, repo *mongo.MongoRepository, uow *mongo.MongoUnitOfWork) error {
	// Example: Insert a new document
	document := bson.M{"name": "John Doe", "age": 30}
	if err := repo.InsertOne(ctx, "collection_name", document); err != nil {
		return fmt.Errorf("error inserting document: %w", err)
	}

	// Example: Update an existing document
	update := bson.M{"age": 31}
	qb := mongo.NewMongoQueryBuilder().
		SetFilter(bson.M{"name": "John Doe"}).
		SetUpdate(update)

	if err := repo.UpdateOne(ctx, "collection_name", qb); err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}
```
