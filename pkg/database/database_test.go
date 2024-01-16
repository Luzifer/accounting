package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDSN = "file::memory:?cache=shared"

// const testDSN = "/tmp/test.db"
func TestCreateDB(t *testing.T) {
	_, err := New("sqlite", testDSN)
	require.NoError(t, err)

	_, err = New("IDoNotExist", testDSN)
	require.Error(t, err)
}

func TestAccountManagement(t *testing.T) {
	dbc, err := New("sqlite", testDSN)
	require.NoError(t, err)

	// Try to create invalid account type
	_, err = dbc.CreateAccount("test", AccountType("foobar"))
	require.Error(t, err)

	// Create account for testing and validate ID
	act, err := dbc.CreateAccount("test", AccountTypeBudget)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, act.ID)

	// Store ID
	actID := act.ID

	// Fetch account by ID
	act, err = dbc.GetAccount(actID)
	assert.NoError(t, err)
	assert.Equal(t, actID, act.ID)
	assert.Equal(t, "test", act.Name)

	// List all accounts
	accs, err := dbc.ListAccounts(false)
	assert.NoError(t, err)
	assert.Len(t, accs, 2)

	// Hide account and list again
	assert.NoError(t, dbc.UpdateAccountHidden(actID, true))
	accs, err = dbc.ListAccounts(false)
	assert.NoError(t, err)
	assert.Len(t, accs, 1)

	// Unhide account and list again
	assert.NoError(t, dbc.UpdateAccountHidden(actID, false))
	accs, err = dbc.ListAccounts(false)
	assert.NoError(t, err)
	assert.Len(t, accs, 2)

	// List accounts from other type
	accs, err = dbc.ListAccountsByType(AccountTypeCategory, false)
	assert.NoError(t, err)
	assert.Len(t, accs, 1)

	// List accounts from existing type
	accs, err = dbc.ListAccountsByType(AccountTypeBudget, false)
	assert.NoError(t, err)
	assert.Len(t, accs, 1)

	// Rename account
	assert.NoError(t, dbc.UpdateAccountName(actID, "renamed"))
	act, err = dbc.GetAccount(actID)
	assert.NoError(t, err)
	assert.Equal(t, actID, act.ID)
	assert.Equal(t, "renamed", act.Name)
}

//nolint:funlen
func TestTransactions(t *testing.T) {
	dbc, err := New("sqlite", testDSN)
	require.NoError(t, err)

	checkAcctBal := func(bals []AccountBalance, act uuid.UUID, bal float64) {
		for _, b := range bals {
			if b.ID == act {
				assert.Equal(t, bal, b.Balance)
				return
			}
		}

		t.Errorf("account %s balance not found", act)
	}

	// Set up some accounts for testing
	tb1, err := dbc.CreateAccount("test1", AccountTypeBudget)
	require.NoError(t, err)
	tb2, err := dbc.CreateAccount("test2", AccountTypeBudget)
	require.NoError(t, err)
	tt, err := dbc.CreateAccount("test", AccountTypeTracking)
	require.NoError(t, err)
	tc, err := dbc.CreateAccount("test", AccountTypeCategory)
	require.NoError(t, err)

	// Try to enter an invalid tx
	_, err = dbc.CreateTransaction(Transaction{
		Payee:       "ACME Inc.",
		Description: "Monthly Income",
		Amount:      1000,
		Account:     uuid.NullUUID{UUID: tb1.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true},
		Cleared:     true,
	})
	require.Error(t, err)

	// Lets earn some money
	tx, err := dbc.CreateTransaction(Transaction{
		Time:        time.Now(),
		Payee:       "ACME Inc.",
		Description: "Monthly Income",
		Amount:      1000,
		Account:     uuid.NullUUID{UUID: tb1.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true},
		Cleared:     true,
	})
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, tx.ID)

	// Now we should have money…
	bals, err := dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 1000)
	checkAcctBal(bals, tb2.ID, 0)
	checkAcctBal(bals, tt.ID, 0)
	checkAcctBal(bals, tc.ID, 0)
	checkAcctBal(bals, UnallocatedMoney, 1000)

	// Lets redistribute the money
	require.NoError(t, dbc.TransferMoney(UnallocatedMoney, tc.ID, 500))
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 1000)
	checkAcctBal(bals, tb2.ID, 0)
	checkAcctBal(bals, tt.ID, 0)
	checkAcctBal(bals, tc.ID, 500)
	checkAcctBal(bals, UnallocatedMoney, 500)

	// Now transfer some money to another budget account
	require.NoError(t, dbc.TransferMoney(tb1.ID, tb2.ID, 100))
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 900)
	checkAcctBal(bals, tb2.ID, 100)
	checkAcctBal(bals, tt.ID, 0)
	checkAcctBal(bals, tc.ID, 500)
	checkAcctBal(bals, UnallocatedMoney, 500)

	// And some to a tracking account (needs category)
	require.NoError(t, dbc.TransferMoneyWithCategory(tb1.ID, tt.ID, 100, tc.ID))
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 800)
	checkAcctBal(bals, tb2.ID, 100)
	checkAcctBal(bals, tt.ID, 100)
	checkAcctBal(bals, tc.ID, 400)
	checkAcctBal(bals, UnallocatedMoney, 500)

	// We might also spend money
	lltx, err := dbc.CreateTransaction(Transaction{
		Time:        time.Now(),
		Payee:       "Landlord",
		Description: "Rent",
		Amount:      -100,
		Account:     uuid.NullUUID{UUID: tb1.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: tc.ID, Valid: true},
		Cleared:     false,
	})
	require.NoError(t, err)
	assert.False(t, lltx.Cleared)
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 700)
	checkAcctBal(bals, tb2.ID, 100)
	checkAcctBal(bals, tt.ID, 100)
	checkAcctBal(bals, tc.ID, 300)
	checkAcctBal(bals, UnallocatedMoney, 500)

	// List transactions
	txs, err := dbc.ListTransactionsByAccount(tb1.ID, time.Time{})
	require.NoError(t, err)
	assert.Len(t, txs, 4)

	txs, err = dbc.ListTransactionsByAccount(UnallocatedMoney, time.Time{})
	require.NoError(t, err)
	assert.Len(t, txs, 2)

	// Oh, wrong category
	require.NoError(t, dbc.UpdateTransactionCategory(lltx.ID, UnallocatedMoney))
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 700)
	checkAcctBal(bals, tb2.ID, 100)
	checkAcctBal(bals, tt.ID, 100)
	checkAcctBal(bals, tc.ID, 400)
	checkAcctBal(bals, UnallocatedMoney, 400)

	// Lets try to move it to a broken category
	require.Error(t, dbc.UpdateTransactionCategory(lltx.ID, tt.ID))

	// Lets try to move an account instead of a tx
	require.Error(t, dbc.UpdateTransactionCategory(tb1.ID, UnallocatedMoney))

	// Clear the tx
	require.NoError(t, dbc.UpdateTransactionCleared(lltx.ID, true))
	lltx, err = dbc.GetTransactionByID(lltx.ID)
	require.NoError(t, err)
	assert.True(t, lltx.Cleared)

	// We made an error and didn't pay the landlord
	require.NoError(t, dbc.DeleteTransaction(lltx.ID))
	bals, err = dbc.ListAccountBalances(false)
	require.NoError(t, err)
	checkAcctBal(bals, tb1.ID, 800)
	checkAcctBal(bals, tb2.ID, 100)
	checkAcctBal(bals, tt.ID, 100)
	checkAcctBal(bals, tc.ID, 400)
	checkAcctBal(bals, UnallocatedMoney, 500)

	// Get a deleted transaction
	_, err = dbc.GetTransactionByID(lltx.ID)
	require.Error(t, err)

	// List transactions
	txs, err = dbc.ListTransactionsByAccount(tb1.ID, time.Time{})
	require.NoError(t, err)
	assert.Len(t, txs, 3)
}
