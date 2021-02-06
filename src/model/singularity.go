package model

import (
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/model/account"
	"io/ioutil"
)

type Singularity struct {
	Balances map[account.Account]uint `json:"balances"`
	Symbol   string                   `json:"symbol"`
}

func LoadSingularity(path string) (*Singularity, error) {
	logger.Info(fmt.Sprintf("Start reading singularity from disk with path %s!", path), "Layer:Model", "Func:LoadSingularity", "Status:Start")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("File couldn't be read!", err, "Layer:model", "Func:LoadSingularity", "Status:error")
		return nil, err
	}

	var loadedSingularity Singularity
	err = json.Unmarshal(content, &loadedSingularity)
	if err != nil {
		logger.Error("Content of file couldn't be marshalled!", err, "Layer:model", "Func:LoadSingularity", "Status:error")
		return nil, err
	}

	logger.Info(fmt.Sprintf("Start reading singularity from disk with path %s!", path), "Layer:Model", "Func:LoadSingularity", "Status:End")
	return &loadedSingularity, nil
}
