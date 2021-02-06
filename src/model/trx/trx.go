package trx

import (
	"go.mod/config"
	"go.mod/model"
)

type Trx struct {
	From  model.Account `json:"from"`
	To    model.Account `json:"to"`
	Value uint          `json:"value"`
	Data  string        `json:"data"`
}

func (t Trx) IsReward() bool {
	return t.Data == config.RewardTrx
}

