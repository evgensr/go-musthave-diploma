package model

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"math"
	"time"
)

type User struct {
	ID                int64   `json:"id"`
	Login             string  `json:"login"`
	Password          string  `json:"password,omitempty"`
	Balance           float64 `json:"balance"`
	EncryptedPassword string  `json:"-"`
}

type Order struct {
	ID     string    `json:"number,omitempty"`
	Status string    `json:"status,omitempty"`
	Amount int64     `json:"accrual,omitempty"`
	Date   time.Time `json:"uploaded_at,omitempty"`
	Type   string    `json:"type,omitempty"`
	UserID int64     `json:"user_id,omitempty"`
}

type AccrualOrder struct {
	ID     string `json:"order,omitempty"`
	Status string `json:"status,omitempty"`
	Amount int64  `json:"accrual,omitempty"`
}

type Balance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}

type Withdrawal struct {
	ID     string    `json:"order,omitempty"`
	Amount int64     `json:"sum,omitempty"`
	Date   time.Time `json:"processed_at,omitempty"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Required, validation.Length(2, 100)),
		validation.Field(&u.Password, validation.By(requiredif(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return nil
		}
		u.EncryptedPassword = enc
	}
	return nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil

}

func (o *Order) MarshalJSON() ([]byte, error) {

	type newOrder struct {
		ID     string  `json:"number,omitempty"`
		Status string  `json:"status"`
		Amount float64 `json:"accrual,omitempty"`
		Date   string  `json:"uploaded_at,omitempty"`
	}

	nb := newOrder{
		Amount: math.Abs(float64(o.Amount)) / 100,
		Date:   o.Date.Format(time.RFC3339),
		ID:     fmt.Sprint(o.ID),
		Status: o.Status,
	}

	return json.Marshal(nb)
}

func (w *Withdrawal) MarshalJSON() ([]byte, error) {

	type newWithdrawal struct {
		ID     string  `json:"order,omitempty"`
		Amount float64 `json:"sum,omitempty"`
		Date   string  `json:"processed_at,omitempty"`
	}

	nb := newWithdrawal{
		ID:     fmt.Sprint(w.ID),
		Amount: math.Abs(float64(w.Amount)),
		Date:   w.Date.Format(time.RFC3339),
	}

	return json.Marshal(nb)
}
