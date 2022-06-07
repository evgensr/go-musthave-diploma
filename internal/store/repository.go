package store

import "github.com/evgensr/go-musthave-diploma/internal/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByLogin(string) (*model.User, error)
}
