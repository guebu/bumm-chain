package state

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"go.mod/model/trx"
	"time"
)


type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

type Block struct {
	Header BlockHeader // metadata (parent block hash + time)
	TRXs []trx.Trx // new transactions only (payload)
}

type BlockHeader struct {
	Parent Hash // parent block reference
	Time int64
}


func (bh BlockHeader) ToString() string {
	sep := "--"
	t := time.Unix(bh.Time, 0)
	s := "### Blockheader START ###\n" + "Parent: " + fmt.Sprintf("%x", bh.Parent) + sep + "Time: " + t.String() + "\n### Blockheader END ###\n"
	return s
}

func (b Block) ToString() string {
	start := "### Block START ###\n" + "Blockheader: " + b.Header.ToString()
	s := "\n### Transactions START ###\n"
	end := "\n### Block END ###\n"
	for i, t := range b.TRXs {
		s = s + "Trx. No.: " + fmt.Sprintf("%d", i) + "  " +  t.ToString() + "\n"
	}
	return start + s + end
}


func (b BlockFS) ToString() string {

	sep := "--"

	s := "Hash: " + fmt.Sprintf("%x", b.Key) + sep + "Value: " + b.Value.ToString()
	return s
}

func (b Block) Hash() (*Hash, error) {
	logger.Info("Start hashing the block...", "Layer:Model", "Func:Hash", "Status:Start")
	logger.Info(fmt.Sprintf("Block-Data BEFORE HASHING - Hash from Header: %x", b.Header.Parent), "Layer:Model", "Func:Hash", "Status:Pending")
	blockJson, err := json.Marshal(b)

	if err != nil {
		logger.Error("Error in Marshalling the block!", err, "Layer:Model", "Func:Hash" )
		return nil, err
	}

	var byteArray [32]byte
	var blockHash Hash = byteArray
	blockHash = sha256.Sum256(blockJson)

	logger.Info( fmt.Sprintf("%x", blockHash), "Layer:Model", "Func:Hash")


	logger.Info(fmt.Sprintf("Block-Data AFTER HASHING - Hash from Header: %x", b.Header.Parent), "Layer:Model", "Func:Hash", "Status:Pending")
	logger.Info("End of hashing the block...", "Layer:Model", "Func:Hash", "Status:End")
	return &blockHash, nil
}

func NewBlock(parentHash Hash, creationTime int64, memPool []trx.Trx) *Block {
	logger.Info("Start of creating new block...", "Layer:Model", "Func:NewBlock", "Status:Start")
	logger.Info( fmt.Sprintf("Block data - Parent-Hash: %x", parentHash), "Layer:Model", "Func:NewBlock", "Status:Pending")

	newBlock := Block{ 	BlockHeader{ Time: creationTime,
									Parent: parentHash },
						memPool }

	logger.Info("#### New Block Info ####", "Layer:Model", "Func:NewBlock", "Status:Pending")
	logger.Info(newBlock.ToString(), "Layer:Model", "Func:NewBlock", "Status:Pending")
	logger.Info("#### New Block Info ####", "Layer:Model", "Func:NewBlock", "Status:Pending")
	logger.Info("Start of creating new block...", "Layer:Model", "Func:NewBlock", "Status:End")
	return &newBlock
}
