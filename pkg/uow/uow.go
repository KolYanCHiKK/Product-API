package uow

import (
	"app/product-api/pkg/db"
	"context"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	Db *db.Db
}

func NewUnitOfWork(db *db.Db) *UnitOfWork {
	return &UnitOfWork{
		Db: db,
	}
}

type TransactionsManager struct {
	UoW *UnitOfWork
}

func NewTransactionsManager(uow *UnitOfWork) *TransactionsManager {
	return &TransactionsManager{
		UoW: uow,
	}
}

func (m *TransactionsManager) Execute(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.UoW.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
