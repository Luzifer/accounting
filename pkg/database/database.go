// Package database has a database client to access the transactions
// and account database together with helpers to interact with those
// tables
package database

import (
	"errors"
	"fmt"
	"math"
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

	for i := range migrateCreateAccounts {
		a := migrateCreateAccounts[i]
		if err = db.Save(&a).Error; err != nil {
			return nil, fmt.Errorf("ensuring default account %q: %w", a.Name, err)
		}
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
		tx, err := c.GetTransactionByID(id)
		if err != nil {
			return err
		}

		if tx.PairKey.Valid {
			// We got a paired transaction which would be out-of-sync if we
			// only delete one part of it so instead of doing a delete on the
			// ID of the transaction, we do a delete on the pair-key
			return db.Delete(&Transaction{}, "pair_key = ?", tx.PairKey.UUID).Error
		}

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
func (c *Client) ListAccountBalances(showHidden bool) (a []AccountBalance, err error) {
	accs, err := c.ListAccounts(showHidden)
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
				// Fix database doing e-15 stuff by rounding to full cents
				ab.Balance = math.Round(*v*100) / 100 //nolint:gomnd
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
//
//revive:disable-next-line:flag-parameter
func (c *Client) ListAccounts(showHidden bool) (a []Account, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		q := db.Model(&Account{})

		if !showHidden {
			q = q.Where("hidden = ?", false)
		}

		return q.Find(&a).Error
	}); err != nil {
		return a, fmt.Errorf("listing accounts: %w", err)
	}

	return a, nil
}

// ListAccountsByType returns a list of all accounts of the given type
//
//revive:disable-next-line:flag-parameter
func (c *Client) ListAccountsByType(at AccountType, showHidden bool) (a []Account, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		q := db.Where("type = ?", at)

		if !showHidden {
			q = q.Where("hidden = ?", false)
		}

		return q.Find(&a).Error
	}); err != nil {
		return a, fmt.Errorf("listing accounts: %w", err)
	}

	return a, nil
}

// ListTransactions retrieves all transactions
func (c *Client) ListTransactions(since, until time.Time) (txs []Transaction, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.
			Where("time >= ? and time <= ?", since, until).
			Find(&txs).
			Error
	}); err != nil {
		return txs, fmt.Errorf("listing transactions: %w", err)
	}

	return txs, nil
}

// ListTransactionsByAccount retrieves all transactions for an account
// or category
func (c *Client) ListTransactionsByAccount(acc uuid.UUID, since, until time.Time) (txs []Transaction, err error) {
	if err = c.retryRead(func(db *gorm.DB) error {
		return db.
			Where("time >= ? and time <= ?", since, until).
			Find(&txs, "account = ? OR category = ?", acc, acc).
			Error
	}); err != nil {
		return txs, fmt.Errorf("listing transactions: %w", err)
	}

	return txs, nil
}

// MarkAccountReconciled marks all cleared transactions as reconciled.
// The account balance is NOT checked in this method.
func (c *Client) MarkAccountReconciled(acc uuid.UUID) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		return db.
			Model(&Transaction{}).
			Where("account = ?", acc).
			Where("cleared = ?", true).
			Update("reconciled", true).
			Error
	}); err != nil {
		return fmt.Errorf("updating transactions: %w", err)
	}

	return nil
}

// TransferMoney creates new Transactions for the given account
// transfer. The account type of the from and to account must match
// for this to work.
func (c *Client) TransferMoney(from, to uuid.UUID, amount float64, description string) (err error) {
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

	pairKey := uuid.Must(uuid.NewRandom())

	var txs []*Transaction
	switch fromAcc.Type {
	case AccountTypeBudget, AccountTypeTracking:
		// Create TX with null-category
		txs = []*Transaction{
			{
				Time:        time.Now().UTC(),
				Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Description: description,
				Amount:      -amount,
				Account:     uuid.NullUUID{UUID: from, Valid: true},
				Category:    uuid.NullUUID{},
				Cleared:     false,
				PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
			},
			{
				Time:        time.Now().UTC(),
				Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Description: description,
				Amount:      amount,
				Account:     uuid.NullUUID{UUID: to, Valid: true},
				Category:    uuid.NullUUID{},
				Cleared:     false,
				PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
			},
		}

	case AccountTypeCategory:
		// Create TX with null-account
		txs = []*Transaction{
			{
				Time:        time.Now().UTC(),
				Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Description: description,
				Amount:      -amount,
				Account:     uuid.NullUUID{},
				Category:    uuid.NullUUID{UUID: from, Valid: true},
				Cleared:     false,
				PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
			},
			{
				Time:        time.Now().UTC(),
				Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
				Description: description,
				Amount:      amount,
				Account:     uuid.NullUUID{},
				Category:    uuid.NullUUID{UUID: to, Valid: true},
				Cleared:     false,
				PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
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
func (c *Client) TransferMoneyWithCategory(from, to uuid.UUID, amount float64, description string, category uuid.UUID) (err error) {
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

	pairKey := uuid.Must(uuid.NewRandom())

	if err = c.retryTx(func(tx *gorm.DB) (err error) {
		fromTx := Transaction{
			Time:        time.Now().UTC(),
			Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
			Description: description,
			Amount:      -amount,
			Account:     uuid.NullUUID{UUID: from, Valid: true},
			Category:    uuid.NullUUID{},
			Cleared:     false,
			PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
		}

		if fromAcc.Type == AccountTypeBudget {
			fromTx.Category = uuid.NullUUID{UUID: category, Valid: true}
		}

		toTx := Transaction{
			Time:        time.Now().UTC(),
			Payee:       fmt.Sprintf("Transfer: %s → %s", fromAcc.Name, toAcc.Name),
			Description: description,
			Amount:      amount,
			Account:     uuid.NullUUID{UUID: to, Valid: true},
			Category:    uuid.NullUUID{},
			Cleared:     false,
			PairKey:     uuid.NullUUID{UUID: pairKey, Valid: true},
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

// UpdateTransaction takes a transaction, fetches the stored transaction
// applies some sanity actions and stores it back to the database
func (c *Client) UpdateTransaction(txID uuid.UUID, tx Transaction) (err error) {
	if err = c.retryTx(func(db *gorm.DB) error {
		var oldTX Transaction
		if err := db.First(&oldTX, "id = ?", txID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return backoff.NewErrCannotRetry(fmt.Errorf("fetching old transaction: %w", err)) //nolint:wrapcheck
			}
			return fmt.Errorf("fetching old transaction: %w", err)
		}

		tx.ID = txID
		tx.Account = oldTX.Account // Changing that would create chaos
		tx.PairKey = oldTX.PairKey // Updating a paired tx should not decouple it

		if err = tx.Validate(c); err != nil {
			return fmt.Errorf("validating transaction: %w", err)
		}

		if err = db.Save(&tx).Error; err != nil {
			return fmt.Errorf("saving transaction: %w", err)
		}

		if !oldTX.PairKey.Valid || tx.Amount == oldTX.Amount {
			// is not a paired transaction or amount did not change: skip rest
			return nil
		}

		// transaction is paired and amount changed, we need to update the
		// paired transaction too or it will cause trouble

		if err = db.Model(&Transaction{}).
			Where("pair_key = ?", oldTX.PairKey.UUID).
			Where("id <> ?", oldTX.ID).
			Update("amount", -tx.Amount).
			Error; err != nil {
			return fmt.Errorf("updating amount for paired transaction: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("updating transaction: %w", err)
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
