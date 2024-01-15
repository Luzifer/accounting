// Package database has a database client to access the transactions
// and account database together with helpers to interact with those
// tables
package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/Luzifer/go_helpers/v2/backoff"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dbMaxRetries = 5

type (
	// Client is the database client
	Client struct {
		db *gorm.DB
	}
)

// New creates a new database client for the given DSN
func New(dbtype, dsn string) (*Client, error) {
	var conn gorm.Dialector
	switch dbtype {
	case "cockroach", "crdb", "postgres", "postgresql":
		conn = postgres.Open(dsn)

	case "sqlite", "sqlite3":
		conn = sqlite.Open(dsn)

	default:
		return nil, fmt.Errorf("unknown db-type %s", dbtype)
	}

	db, err := gorm.Open(conn, &gorm.Config{
		Logger: logger.New(loggerWriter{logrus.StandardLogger().WriterLevel(logrus.TraceLevel)}, logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err = db.AutoMigrate(
		&Account{},
		&Transaction{},
	); err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}

	if err = db.Save(&Account{
		BaseModel: BaseModel{
			ID: UnallocatedMoney,
		},
		Name:   "Unallocated Money",
		Type:   AccountTypeCategory,
		Hidden: false,
	}).Error; err != nil {
		return nil, fmt.Errorf("ensuring unallocated money category: %w", err)
	}

	return &Client{
		db: db,
	}, nil
}

// CreateAccount creates and returns a new account of the given type
func (c *Client) CreateAccount(name string, accType AccountType) (a Account, err error) {
	a = Account{
		Name: name,
		Type: accType,
	}

	if !accType.IsValid() {
		return a, fmt.Errorf("invalid account type %s", accType)
	}

	if err = c.retryTx(func(db *gorm.DB) error {
		return db.Save(&a).Error
	}); err != nil {
		return a, fmt.Errorf("creating account: %w", err)
	}

	return a, nil
}

// CreateTransaction takes a prepared transaction and stores it
func (c *Client) CreateTransaction(tx Transaction) (ntx Transaction, err error) {
	if err = tx.Validate(c); err != nil {
		return tx, fmt.Errorf("validating transaction: %w", err)
	}

	if err = c.retryTx(func(db *gorm.DB) error {
		return db.Save(&tx).Error
	}); err != nil {
		return tx, fmt.Errorf("creating transaction: %w", err)
	}

	return tx, nil
}

// DeleteTransaction deletes a transaction
func (c *Client) DeleteTransaction(id uuid.UUID) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		return db.Delete(&Transaction{}, "id = ?", id).Error
	}); err != nil {
		return fmt.Errorf("deleting transaction: %w", err)
	}

	return nil
}

// GetAccount retrieves an Account using its ID
func (c *Client) GetAccount(id uuid.UUID) (a Account, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.First(&a, "id = ?", id).Error
	}); err != nil {
		return a, fmt.Errorf("fetching account: %w", err)
	}

	return a, nil
}

// GetTransactionByID returns a single transaction by its ID
func (c *Client) GetTransactionByID(id uuid.UUID) (tx Transaction, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.First(&tx, "id = ?", id).Error
	}); err != nil {
		return tx, fmt.Errorf("getting transaction: %w", err)
	}

	return tx, nil
}

// ListAccountBalances returns a list of accounts with their
// corresponding balance
func (c *Client) ListAccountBalances() (a []AccountBalance, err error) {
	accs, err := c.ListAccounts()
	if err != nil {
		return nil, fmt.Errorf("listing accounts: %w", err)
	}

	for _, acc := range accs {
		if err = c.retryRead(func(db *gorm.DB) error {
			q := db.
				Model(&Transaction{})

			if acc.Type == AccountTypeCategory {
				q = q.Where("category = ?", acc.ID)
			} else {
				q = q.Where("account = ?", acc.ID)
			}

			ab := AccountBalance{
				Account: acc,
				Balance: 0,
			}

			var v *float64
			if err = q.
				Select("sum(amount)").
				Scan(&v).
				Error; err != nil {
				return fmt.Errorf("getting sum: %w", err)
			}

			if v != nil {
				ab.Balance = *v
			}

			a = append(a, ab)
			return nil
		}); err != nil {
			return nil, fmt.Errorf("getting account balance for %s: %w", acc.ID, err)
		}
	}

	return a, nil
}

// ListAccounts returns a list of all accounts
func (c *Client) ListAccounts() (a []Account, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.Find(&a, "hidden = ?", false).Error
	}); err != nil {
		return a, fmt.Errorf("listing accounts: %w", err)
	}

	return a, nil
}

// ListAccountsByType returns a list of all accounts of the given type
func (c *Client) ListAccountsByType(at AccountType) (a []Account, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.Find(&a, "type = ?", at).Error
	}); err != nil {
		return a, fmt.Errorf("listing accounts: %w", err)
	}

	return a, nil
}

// ListTransactionsByAccount retrieves all transactions for an account
// or category
func (c *Client) ListTransactionsByAccount(acc uuid.UUID, since time.Time) (txs []Transaction, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.
			Where("time >= ?", since).
			Find(&txs, "account = ? OR category = ?", acc, acc).
			Error
	}); err != nil {
		return txs, fmt.Errorf("listing transactions: %w", err)
	}

	return txs, nil
}

// TransferMoney creates new Transactions for the given account
// transfer. The account type of the from and to account must match
// for this to work.
func (c *Client) TransferMoney(from, to uuid.UUID, amount float64) (err error) {
	var fromAcc, toAcc Account

	if fromAcc, err = c.GetAccount(from); err != nil {
		return fmt.Errorf("getting source account: %w", err)
	}

	if toAcc, err = c.GetAccount(to); err != nil {
		return fmt.Errorf("getting target account: %w", err)
	}

	if fromAcc.Type != toAcc.Type {
		return fmt.Errorf("account type mismatch: %s != %s", fromAcc.Type, toAcc.Type)
	}

	var txs []*Transaction
	switch fromAcc.Type {
	case AccountTypeBudget, AccountTypeTracking:
		// Create TX with null-category
		txs = []*Transaction{
			{
				Time:        time.Now().UTC(),
				Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Amount:      -amount,
				Account:     uuid.NullUUID{UUID: from, Valid: true},
				Category:    uuid.NullUUID{},
				Cleared:     false,
			},
			{
				Time:        time.Now().UTC(),
				Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Amount:      amount,
				Account:     uuid.NullUUID{UUID: to, Valid: true},
				Category:    uuid.NullUUID{},
				Cleared:     false,
			},
		}

	case AccountTypeCategory:
		// Create TX with null-account
		txs = []*Transaction{
			{
				Time:        time.Now().UTC(),
				Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Amount:      -amount,
				Account:     uuid.NullUUID{},
				Category:    uuid.NullUUID{UUID: from, Valid: true},
				Cleared:     false,
			},
			{
				Time:        time.Now().UTC(),
				Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Amount:      amount,
				Account:     uuid.NullUUID{},
				Category:    uuid.NullUUID{UUID: to, Valid: true},
				Cleared:     false,
			},
		}
	}

	if err = c.retryTx(func(tx *gorm.DB) (err error) {
		for _, t := range txs {
			if err = tx.Save(t).Error; err != nil {
				return fmt.Errorf("saving transaction: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("creating transactions: %w", err)
	}

	return nil
}

// TransferMoneyWithCategory creates new Transactions for the given
// account transfer. This is not possible for category type accounts.
func (c *Client) TransferMoneyWithCategory(from, to uuid.UUID, amount float64, category uuid.UUID) (err error) {
	var fromAcc, toAcc Account

	if fromAcc, err = c.GetAccount(from); err != nil {
		return fmt.Errorf("getting source account: %w", err)
	}

	if toAcc, err = c.GetAccount(to); err != nil {
		return fmt.Errorf("getting target account: %w", err)
	}

	if fromAcc.Type == AccountTypeCategory || toAcc.Type == AccountTypeCategory {
		return fmt.Errorf("transfer contained category-type account")
	}

	if err = c.retryTx(func(tx *gorm.DB) (err error) {
		fromTx := Transaction{
			Time:        time.Now().UTC(),
			Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
			Amount:      -amount,
			Account:     uuid.NullUUID{UUID: from, Valid: true},
			Category:    uuid.NullUUID{},
			Cleared:     false,
		}

		if fromAcc.Type == AccountTypeBudget {
			fromTx.Category = uuid.NullUUID{UUID: category, Valid: true}
		}

		toTx := Transaction{
			Time:        time.Now().UTC(),
			Description: fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
			Amount:      amount,
			Account:     uuid.NullUUID{UUID: to, Valid: true},
			Category:    uuid.NullUUID{},
			Cleared:     false,
		}

		if toAcc.Type == AccountTypeBudget {
			toTx.Category = uuid.NullUUID{UUID: category, Valid: true}
		}

		for _, t := range []*Transaction{&fromTx, &toTx} {
			if err = tx.Save(t).Error; err != nil {
				return fmt.Errorf("saving transaction: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("creating transactions: %w", err)
	}

	return nil
}

// UpdateAccountHidden updates the hidden flag for the given Account
func (c *Client) UpdateAccountHidden(id uuid.UUID, hidden bool) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		return db.
			Model(&Account{}).
			Where("id = ?", id).
			Update("hidden", hidden).
			Error
	}); err != nil {
		return fmt.Errorf("updating account: %w", err)
	}

	return nil
}

// UpdateAccountName sets a new name for the given account ID
func (c *Client) UpdateAccountName(id uuid.UUID, name string) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		return db.
			Model(&Account{}).
			Where("id = ?", id).
			Update("name", name).
			Error
	}); err != nil {
		return fmt.Errorf("updating account: %w", err)
	}

	return nil
}

// UpdateTransactionCategory modifies the category of the given
// transaction. (It is not possible to remove a category with this)
func (c *Client) UpdateTransactionCategory(id uuid.UUID, cat uuid.UUID) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		var tx Transaction
		if err = db.First(&tx, "id = ?", id).Error; err != nil {
			return fmt.Errorf("fetching transaction: %w", err)
		}

		tx.Category = uuid.NullUUID{UUID: cat, Valid: true}
		if err = tx.Validate(c); err != nil {
			return fmt.Errorf("validating transaction: %w", err)
		}

		if err = db.
			Save(&tx).
			Error; err != nil {
			return fmt.Errorf("saving transaction: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("updating transaction: %w", err)
	}

	return nil
}

// UpdateTransactionCleared modifies the "cleared" flag for the given
// transaction
func (c *Client) UpdateTransactionCleared(id uuid.UUID, cleared bool) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		return db.
			Model(&Transaction{}).
			Where("id = ?", id).
			Update("cleared", cleared).
			Error
	}); err != nil {
		return fmt.Errorf("updating transaction: %w", err)
	}

	return nil
}

func (c *Client) retryRead(fn func(db *gorm.DB) error) error {
	//nolint:wrapcheck
	return backoff.NewBackoff().
		WithMaxIterations(dbMaxRetries).
		Retry(func() error {
			err := fn(c.db)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return backoff.NewErrCannotRetry(err)
			}
			return err
		})
}

func (c *Client) retryTx(fn func(db *gorm.DB) error) error {
	//nolint:wrapcheck
	return backoff.NewBackoff().
		WithMaxIterations(dbMaxRetries).
		Retry(func() error {
			return c.db.Transaction(fn)
		})
}
