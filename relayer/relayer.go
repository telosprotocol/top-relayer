package relayer

import (
	"sync"

	"toprelayer/config"
	"toprelayer/relayer/eth2top"
	"toprelayer/relayer/top2eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

type IChainRelayer interface {
	Init(fromUrl, toUrl, keypath, pass string, chainid uint64, contract common.Address, batch int) error
	StartRelayer(*sync.WaitGroup) error
	ChainId() uint64
}

func StartRelayer(wg *sync.WaitGroup, handlercfg *config.HeaderSyncConfig, chainpass map[uint64]string) (err error) {
	handler := NewHeaderSyncHandler(handlercfg)
	err = handler.Init(wg, chainpass)
	if err != nil {
		return err
	}
	return handler.StartRelayer()
}

func GetRelayer(chain string) (relayer IChainRelayer) {
	switch chain {
	case config.ETH_CHAIN:
		relayer = new(top2eth.Top2EthRelayer)
	case config.TOP_CHAIN:
		relayer = new(eth2top.Eth2TopRelayer)
	default:
		logger.Error("Unsupport chain:", chain)
	}
	return
}
