package mysql

import (
	"database/sql"
)

type MySQLUnitOfWork struct {
	Tx *sql.Tx
}

func (uow *MySQLUnitOfWork) GetRepository() IRepository {
	return &MySQLRepository{Tx: uow.Tx}
}

func (uow *MySQLUnitOfWork) Commit() error {
	return uow.Tx.Commit()
}

func (uow *MySQLUnitOfWork) Rollback() error {
	return uow.Tx.Rollback()
}
