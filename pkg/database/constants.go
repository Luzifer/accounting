package database

import "github.com/google/uuid"

const constAcctIDNamespace = "17de217e-94d7-4a9b-8833-ecca7f0eb6ca"

var (
	// UnallocatedMoney is a category UUID which is automatically created
	// during database migration phase and therefore always available
	UnallocatedMoney = uuid.NewSHA1(uuid.MustParse(constAcctIDNamespace), []byte("unallocated-money"))

	invalidAcc = uuid.NewSHA1(uuid.MustParse(constAcctIDNamespace), []byte("INVALID ACCOUNT"))
)
