package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	// Account represents a budget, tracking or category account - in
	// general something holding money through the sum of transactions
	Account struct {
		BaseModel
		Name   string      `json:"name"`
		Type   AccountType `json:"type"`
		Hidden bool        `json:"hidden"`
	}

	// AccountBalance wraps an Account and adds the balance
	AccountBalance struct {
		Account
		Balance float64 `json:"balance"`
	}

	// AccountType represents the type of an account
	AccountType string

	// Transaction represents some money movement between, from
	// or to accounts
	Transaction struct {
		BaseModel
		Time        time.Time     `json:"time"`
		Payee       string        `json:"payee"`
		Description string        `json:"description"`
		Amount      float64       `json:"amount"`
		Account     uuid.NullUUID `gorm:"type:uuid" json:"account"`
		Category    uuid.NullUUID `gorm:"type:uuid" json:"category"`
		Cleared     bool          `json:"cleared"`

		PairKey uuid.NullUUID `gorm:"type:uuid" json:"-"`
	}

	// BaseModel is used internally in all other models for common fields
	BaseModel struct {
		ID        uuid.UUID      `gorm:"type:uuid" json:"id"`
		CreatedAt time.Time      `json:"-"`
		UpdatedAt time.Time      `json:"-"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	}
)

// Known values of the AccountType enum
const (
	AccountTypeBudget   AccountType = "budget"
	AccountTypeCategory AccountType = "category"
	AccountTypeTracking AccountType = "tracking"
)

// IsValid checks whether the given AccountType belongs to the known
// types
func (a AccountType) IsValid() bool {
	for _, kat := range []AccountType{
		AccountTypeBudget,
		AccountTypeCategory,
		AccountTypeTracking,
	} {
		if kat == a {
			return true
		}
	}

	return false
}

// BeforeCreate ensures the object UUID is filled
func (b *BaseModel) BeforeCreate(*gorm.DB) (err error) {
	b.ID = uuid.New()
	return nil
}

// Validate executes some basic checks on the transaction
//
//nolint:gocyclo
func (t Transaction) Validate(c *Client) (err error) {
	var errs []error

	if t.Time.IsZero() {
		errs = append(errs, fmt.Errorf("time is zero"))
	}

	if !t.Account.Valid && !t.Category.Valid {
		errs = append(errs, fmt.Errorf("account and category are null"))
	}

	if t.Amount == 0 {
		errs = append(errs, fmt.Errorf("amount is zero"))
	}

	var acc, cat Account
	if t.Account.Valid {
		if acc, err = c.GetAccount(t.Account.UUID); err != nil {
			return fmt.Errorf("fetching account: %w", err)
		}
	}

	if t.Category.Valid {
		if cat, err = c.GetAccount(t.Category.UUID); err != nil {
			return fmt.Errorf("fetching category: %w", err)
		}
	}

	if acc.Type == AccountTypeBudget && !t.Category.Valid && !t.PairKey.Valid {
		errs = append(errs, fmt.Errorf("budget account transactions need a category"))
	}

	if acc.Type == AccountTypeTracking && t.Category.Valid {
		errs = append(errs, fmt.Errorf("tracking account transactions must not have a category"))
	}

	if t.Category.Valid && cat.Type != AccountTypeCategory {
		errs = append(errs, fmt.Errorf("category is not of type category"))
	}

	return errors.Join(errs...)
}
