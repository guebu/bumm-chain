package state

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mod/config"
	"go.mod/model/account"
	"go.mod/model/trx"
	"testing"
)

/*
func TestMain(m *testing.M) {
	fmt.Println("about to start the tests for state!!!!!!!!!")
	os.Exit(m.Run())
}
*/

const account1 account.Account = "guebu"
const account2 account.Account = "ferdl"
const trxAmount 			uint 			= 	uint(10)
const rewAmount				uint 			= 	uint(10)

const tooBigAmount			uint 			= 	uint(2000)
const suitableAmount		uint 			=   uint(500)
const initialBalAmount1 	uint 			= 	uint(1000)
const initialBalAmount2 	uint 			=	uint(1000)

func TestNewStateFromDisk(t *testing.T) {

	fmt.Println("about to start the tests for state!!!!!!!!!")

	state, err := NewStateFromDisk()

	assert.Nil(t, err, "State should be readable!")
	assert.NotNil(t, state, "State object should be available!")


	var trx1 = trx.Trx{
		From: account1,
		To: account2,
		Value: trxAmount,
		Data: "trx1",
	}

	var trx2 = trx.Trx{
		From: account1,
		To: account2,
		Value: trxAmount,
		Data: "trx2",
	}

	err1 := state.Add(trx1)
	err2 := state.Add(trx2)
	assert.Nil(t, err1, "After reading the state from disk, adding a trx shouldn't be a problem!" )
	assert.Nil(t, err2, "After reading the state from disk, adding a trx shouldn't be a problem!" )

	myHash, errPersist := state.Persist()
	assert.Nil(t, errPersist, "Persisting db file shouldn't be a problem!" )
	assert.NotNil(t, myHash, "Hash as result of persisting data should be available!" )
}

func TestNotSuccessfullAddingOfTrxToMemPool(t *testing.T){
	var trx = trx.Trx{
		From: account1,
		To: account2,
		Value: tooBigAmount,
		Data: "",
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
		account2: initialBalAmount2,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.Add(trx)

	assert.NotNil(t, err, "Based on a proper state test trx should be addedd to mem pool!")
}

func TestSuccessfullAddingOfTrxToMemPool(t *testing.T){

	var trx = trx.Trx{
		From: account1,
		To: account2,
		Value: trxAmount,
		Data: "",
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
		account2: initialBalAmount2,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.Add(trx)

	assert.Nil(t, err, "Based on a proper state test trx should be addedd to mem pool!")

}

func TestApplyReward(t *testing.T) {

	var rewardTrx = trx.Trx{
		From: account1,
		To: account1,
		Value: rewAmount,
		Data: config.RewardTrx,
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.apply(rewardTrx)
	assert.Nil(t, err, "Reward transaction should be processed successfully!")
	assert.Equal(t, initialBalAmount1 + rewAmount, state.Balances[account1], "Reward amount should be transfered completly to account!")
}

func TestProperBookingWithPositiveValue(t *testing.T) {

	var trx = trx.Trx{
		From: account1,
		To: account2,
		Value: trxAmount,
		Data: "",
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
		account2: initialBalAmount2,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.apply(trx)
	assert.Nil(t, err, "Transaction should be processed successfully!")
	assert.Equal(t, initialBalAmount1 - trxAmount, state.Balances[account1], "Trx amount should be subtracted completly from account1!")
	assert.Equal(t, initialBalAmount2 + trxAmount, state.Balances[account2], "Trx amount should be transfered completly to account2!")
}

func TestEnoughBalance(t *testing.T) {

	var suitableTrx = trx.Trx{
		From: account1,
		To: account2,
		Value: suitableAmount,
		Data: "",
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
		account2: initialBalAmount2,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.apply(suitableTrx)
	assert.Nil(t, err, "No error should be raised when balance is sufficient!")
	assert.Equal(t, initialBalAmount1 - suitableAmount, state.Balances[account1], "Trx amount should be subtracted completly from account1 when trx was successfull!")
	assert.Equal(t, initialBalAmount2 + suitableAmount, state.Balances[account2], "Trx amount should be transfered completly to account2 when trx was successfull!")
}

func TestNotEnoughBalance(t *testing.T) {

	var tooBigTrx = trx.Trx{
		From: account1,
		To: account2,
		Value: tooBigAmount,
		Data: "",
	}

	var initialBalances = map[account.Account]uint{
		account1: initialBalAmount1,
		account2: initialBalAmount2,
	}

	var state = State{
		Balances: initialBalances,
	}

	err := state.apply(tooBigTrx)
	assert.NotNil(t, err, "Error should be raised when not enough balance!")
	assert.Equal(t, initialBalAmount1, state.Balances[account1], "Initial balance shouldn't be changed for account1 after rejected trx!")
	assert.Equal(t, initialBalAmount2, state.Balances[account2], "Initial balance shouldn't be changed for account 2 after rejected trx!")
}
