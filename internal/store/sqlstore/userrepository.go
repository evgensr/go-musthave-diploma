package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evgensr/go-musthave-diploma/internal/model"
	"github.com/evgensr/go-musthave-diploma/internal/store"
	"log"
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
	var id int64

	err := r.store.db.QueryRow(
		"INSERT INTO bonuses (user_id, order_id, change, type, status) VALUES($1,$2,$3,$4,$5) RETURNING id",
		order.UserID,
		order.ID,
		order.Amount,
		order.Type,
		order.Status,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

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
func (r *UserRepository) SelectAllWithdrawals(ctx context.Context, u int64) (*[]model.Withdrawal, error) {

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

	return &listOrders, nil

}

func (r *UserRepository) doTransaction(fu ...func() error) error {

	tx, err := r.store.db.Begin()
	if err != nil {
		return fmt.Errorf("starting connection failed: %v", err)
	}
	defer tx.Rollback()

	for _, f := range fu {
		err := f()
		if err != nil {
			return fmt.Errorf("transaction failed: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

func (r *UserRepository) SelectOrdersForUpdate(ctx context.Context, oin chan []model.Order, oout chan model.Order) {
	var listOrders []model.Order
	err := r.doTransaction(
		func() error {

			row, err := r.store.db.Query(`SELECT order_id, status FROM bonuses 
										WHERE status not in ('PROCESSED', 'INVALID') LIMIT 1 FOR UPDATE SKIP LOCKED`)
			if err != nil {
				return fmt.Errorf("init select from bonuses failed: %v", err)
			}
			defer row.Close()

			for row.Next() {
				var o model.Order
				err := row.Scan(&o.ID, &o.Status)
				if err != nil {
					return fmt.Errorf("select bonuses for update failed: %v", err)
				}
				listOrders = append(listOrders, o)
			}
			oin <- listOrders
			return nil
		},
		func() error {
			stmtBonuses, err := r.store.db.Prepare("UPDATE bonuses SET change=$1, status=$2 where order_id=$3")

			if err != nil {
				return fmt.Errorf("init update users failed: %v", err)
			}

			stmtUsers, err := r.store.db.Prepare("UPDATE users SET balance=balance+$1 where id=$2;")

			if err != nil {
				return fmt.Errorf("init update users failed: %v", err)
			}

			for {
			insertUpdates:
				select {
				case bonus, ok := <-oout:
					if !ok {
						break insertUpdates
					}

					if _, err = stmtUsers.Exec(bonus.Amount, bonus.UserID); err != nil {
						log.Println(bonus.Amount)
						return fmt.Errorf("update user amount failed: %v", err)
					}

					if _, err = stmtBonuses.Exec(bonus.Amount, bonus.Status, bonus.ID); err != nil {
						return fmt.Errorf("update amount failed: %v", err)
					}

				case <-ctx.Done():
					log.Println("context canceled")
					return nil
				}
				return nil
			}
		})

	if err != nil {
		log.Fatalf("transaction failed: %v", err)
	}
}
