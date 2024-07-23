package database

import (
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
)

const constAcctIDNamespace = "00000000-0000-0000-0000-%012s"

var (
	// UnallocatedMoney is a category UUID which is automatically created
	// during database migration phase and therefore always available
	UnallocatedMoney = makeConstAcctID(1)
	// StartingBalance is a category UUID which is automatically created
	// and hidden during database migration and used in frontend as constant
	StartingBalance = makeConstAcctID(2) //nolint:mnd

	invalidAcc = makeConstAcctID(math.MaxUint32)

	migrateCreateAccounts = []Account{
		{
			BaseModel: BaseModel{ID: UnallocatedMoney},
			Hidden:    false,
			Name:      "Unallocated Money",
			Type:      AccountTypeCategory,
		},
		{
			BaseModel: BaseModel{ID: StartingBalance},
			Hidden:    true,
			Name:      "Starting Balance",
			Type:      AccountTypeCategory,
		},
	}
)

func makeConstAcctID(fixedNumber uint32) uuid.UUID {
	return uuid.MustParse(fmt.Sprintf(constAcctIDNamespace, strconv.FormatUint(uint64(fixedNumber), 16)))
}
