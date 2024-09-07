```bash
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/yourmodule/mysql" // Replace with your actual module path
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func main() {
	// Setup MySQL database connection
	dsn := "user:password@tcp(localhost:3306)/your_database" // Update with your actual DSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Create MySQL repository
	repo := mysql.NewMySQLRepository(db)

	// Begin a transaction
	ctx := context.Background()
	uow, err := repo.BeginTransaction(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Perform transactional operations
	err = performTransactionalOperations(ctx, repo, uow)
	if err != nil {
		// Abort the transaction if an error occurs
		uow.tx.Rollback()
		log.Fatalf("Transaction aborted due to error: %v", err)
	}

	// Commit the transaction
	if err := uow.tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Transaction committed successfully")
}

// performTransactionalOperations performs database operations within a transaction
func performTransactionalOperations(ctx context.Context, repo *mysql.MySQLRepository, uow *mysql.MySQLUnitOfWork) error {
	// Example: Insert a new record
	qbInsert := mysql.NewQueryBuilder().
		SetInsertFields(map[string]interface{}{
			"name": "Jane Doe",
			"age":  29,
		})

	if err := repo.InsertOne(ctx, "users", qbInsert); err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}

	// Example: Update an existing record
	qbUpdate := mysql.NewQueryBuilder().
		SetFilter(map[string]interface{}{"name": "Jane Doe"}).
		SetUpdateFields(map[string]interface{}{
			"age": 30,
		})

	if err := repo.UpdateOne(ctx, "users", qbUpdate); err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	return nil
}
```
