package trx

import (
	"go.mod/config"
	"go.mod/model/account"
)

type Trx struct {
	From  account.Account `json:"from"`
	To    account.Account `json:"to"`
	Value uint            `json:"value"`
	Data  string          `json:"data"`
}

func (t Trx) IsReward() bool {
	return t.Data == config.RewardTrx
}

func NewTrx(fromAcc account.Account, toAcc account.Account, value uint, trxTxt string) *Trx {
	trx := Trx{
		From: fromAcc,
		To: toAcc,
		Value: value,
		Data: trxTxt,
	}

	return &trx

}

