package trx

import (
	"fmt"
	"github.com/guebu/common-utils/logger"
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
	logger.Info("Create new transaction", "Layer:Model", "Func:NewTrx", "Status:Start")
	trx := Trx{
		From: fromAcc,
		To: toAcc,
		Value: value,
		Data: trxTxt,
	}
	logger.Info(trx.ToString(), "Layer:Model", "Func:NewTrx", "Status:Pending")
	logger.Info("Create new transaction", "Layer:Model", "Func:NewTrx", "Status:End")
	return &trx
}

func (t Trx) ToString() string {
	sep := " -- "
	s := "From: " + string(t.From) + sep + "To: " + string(t.To) + sep + "Value: " + fmt.Sprintf("%d", t.Value) + sep + "Data: " + t.Data
	return s
}

