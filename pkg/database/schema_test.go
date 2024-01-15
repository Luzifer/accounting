package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionValidateErrs(t *testing.T) {
	dbc, err := New("sqlite", testDSN)
	require.NoError(t, err)

	// We need one test account of each type
	actB, err := dbc.CreateAccount("test", AccountTypeBudget)
	require.NoError(t, err)
	actT, err := dbc.CreateAccount("test", AccountTypeTracking)
	require.NoError(t, err)

	assert.Error(t, Transaction{
		Time:        time.Time{}, // ERR: Zero time
		Payee:       "test",
		Description: "test",
		Amount:      10,
		Account:     uuid.NullUUID{UUID: actB.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true},
		Cleared:     false,
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      10,
		Account:     uuid.NullUUID{}, // ERR: Both null
		Category:    uuid.NullUUID{}, // ERR: Both null
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      0, // ERR: Zero
		Account:     uuid.NullUUID{UUID: actB.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true},
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      50,
		Account:     uuid.NullUUID{UUID: actB.ID, Valid: true},
		Category:    uuid.NullUUID{}, // ERR: Budget without cat
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      50,
		Account:     uuid.NullUUID{UUID: actT.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true}, // ERR: Tracking with cat
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      50,
		Account:     uuid.NullUUID{UUID: actB.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: actT.ID, Valid: true}, // ERR: Cat is not cat
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      50,
		Account:     uuid.NullUUID{UUID: invalidAcc, Valid: true}, // ERR: Account does not exist
		Category:    uuid.NullUUID{UUID: UnallocatedMoney, Valid: true},
	}.Validate(dbc))

	assert.Error(t, Transaction{
		Time:        time.Now(),
		Payee:       "test",
		Description: "test",
		Amount:      50,
		Account:     uuid.NullUUID{UUID: actB.ID, Valid: true},
		Category:    uuid.NullUUID{UUID: invalidAcc, Valid: true}, // ERR: Cat does not exist
	}.Validate(dbc))
}
