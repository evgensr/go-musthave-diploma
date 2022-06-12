package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evgensr/go-musthave-diploma/internal/model"
	"github.com/evgensr/go-musthave-diploma/internal/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (login, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Login,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByLogin(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, login, encrypted_password FROM users WHERE login = $1",
		email,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Find(id int64) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, login, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// InsertOrder new order
func (r *UserRepository) InsertOrder(ctx context.Context, order model.Order) error {

	// log.Println(order)

	var id int64

	err := r.store.db.QueryRow(
		"INSERT INTO bonuses (user_id, order_id, change, type, status) VALUES($1,$2,$3,$4,$5) RETURNING id",
		order.UserID,
		order.ID,
		order.Amount,
		order.Type,
		order.Status,
	).Scan(&id)

	// log.Println(id)

	if err != nil {
		return err
	}

	return nil
}

// SelectUserForOrder
func (r *UserRepository) SelectUserForOrder(ctx context.Context, order model.Order) (int64, error) {

	var id int64
	err := r.store.db.QueryRow(
		"SELECT users.id FROM users JOIN bonuses ON users.id=bonuses.user_id WHERE order_id=$1 LIMIT 1",
		order.ID,
	).Scan(&id)

	if err != nil {
		return 0, errors.New("err sql")
	}

	return id, nil
}

// SelectAllOrders select all orders
func (r *UserRepository) SelectAllOrders(ctx context.Context, u int64) ([]*model.Order, error) {
	var listOrders []*model.Order

	row, err := r.store.db.Query(`SELECT order_id, status, change, change_date 
														FROM bonuses WHERE user_id=$1 ORDER BY change_date`, u)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var o model.Order
		err := row.Scan(&o.ID, &o.Status, &o.Amount, &o.Date)
		if err != nil {
			return nil, err
		}
		listOrders = append(listOrders, &o)
	}

	return listOrders, nil
}

// SelectBalance select balance
func (r *UserRepository) SelectBalance(ctx context.Context, u int64) (*model.Balance, error) {

	var val model.Balance
	row := r.store.db.QueryRow("SELECT COALESCE(SUM(change), 0), COALESCE(SUM(nullif(LEAST(change, 0),0)),0) FROM bonuses WHERE user_id=$1 AND status='PROCESSED'", u)
	err := row.Scan(&val.Current, &val.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &val, nil

}

// SelectAllWithdrawals select balance
func (r *UserRepository) SelectAllWithdrawals(ctx context.Context, u int64) ([]model.Withdrawal, error) {

	var listOrders []model.Withdrawal

	row, err := r.store.db.Query("SELECT order_id, change, change_date "+
		"FROM bonuses WHERE user_id=$1 AND type='withdraw' ORDER BY change_date", u)
	if err != nil {
		return nil, fmt.Errorf("sql err: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("init select from orders failed: %v", err)
	}

	for row.Next() {
		var o model.Withdrawal
		err := row.Scan(&o.ID, &o.Amount, &o.Date)
		if err != nil {
			return nil, fmt.Errorf("select orders failed: %v", err)
		}
		listOrders = append(listOrders, o)
	}

	return listOrders, nil

}
