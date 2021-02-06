package helper

import (
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/config"
	"path/filepath"
)

func GetSingularityFilePath() string {
	singularityFilePath := filepath.Join(config.BaseFileDirPath, config.FilesDirName, config.FileNameSingularity)
	logger.Info(fmt.Sprintf("SingularityFilePath: %s", singularityFilePath), "Layer:Helper")
	return singularityFilePath
}

func GetTrxDBFilePath() string {
	trxDBFilePath := filepath.Join(config.BaseFileDirPath, config.FilesDirName, config.FileNameTrxDB)
	logger.Info(fmt.Sprintf("TrxDBFilePath: %s", trxDBFilePath), "Layer:Helper")
	return trxDBFilePath
}
