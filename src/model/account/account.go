package account

type Account string

func NewAccount(account string) *Account{
	returnAccount := Account(account)
	return &returnAccount
}

