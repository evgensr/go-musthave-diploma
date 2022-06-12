package store

import (
	"context"
	"github.com/evgensr/go-musthave-diploma/internal/model"
)

type UserRepository interface {
	Create(*model.User) error
	Find(int64) (*model.User, error)
	FindByLogin(string) (*model.User, error)
	SelectUserForOrder(ctx context.Context, order model.Order) (int64, error)
	InsertOrder(ctx context.Context, order model.Order) error
	SelectAllOrders(ctx context.Context, u int64) ([]*model.Order, error)
	SelectBalance(ctx context.Context, user int64) (*model.Balance, error)
	SelectAllWithdrawals(context.Context, int64) ([]model.Withdrawal, error)
}
