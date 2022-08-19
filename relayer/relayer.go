package relayer

import (
	"fmt"
	"sync"

	"toprelayer/config"
	"toprelayer/relayer/crosschainrelayer"
	"toprelayer/relayer/monitor"
	"toprelayer/relayer/toprelayer"

	"github.com/wonderivan/logger"
)

var (
	topRelayers = map[string]IChainRelayer{
		config.ETH_CHAIN:  new(toprelayer.Eth2TopRelayer),
		config.BSC_CHAIN:  new(toprelayer.Bsc2TopRelayer),
		config.HECO_CHAIN: new(toprelayer.Heco2TopRelayer)}

	crossChainRelayer = new(crosschainrelayer.CrossChainRelayer)
)

type IChainRelayer interface {
	Init(chainName string, cfg *config.Relayer, listenUrl string, pass string) error
	StartRelayer(*sync.WaitGroup) error
}

func startOneRelayer(chainName string, relayer IChainRelayer, cfg *config.Relayer, listenUrl string, pass string, wg *sync.WaitGroup) error {
	err := relayer.Init(chainName, cfg, listenUrl, pass)
	if err != nil {
		logger.Error("startOneRelayer error:", err)
		return err
	}

	wg.Add(1)
	go func() {
		err = relayer.StartRelayer(wg)
	}()
	if err != nil {
		logger.Error("relayer.StartRelayer error:", err)
		return err
	}
	return nil
}

func StartRelayer(cfg *config.Config, pass string, wg *sync.WaitGroup) error {
	// start monitor
	err := monitor.MonitorMsgInit(cfg.RelayerToRun)
	if err != nil {
		logger.Error("MonitorMsgInit fail:", err)
		return err
	}

	// start relayer
	topConfig, exist := cfg.RelayerConfig[config.TOP_CHAIN]
	if !exist {
		return fmt.Errorf("not found TOP chain config")
	}
	RelayerConfig, exist := cfg.RelayerConfig[cfg.RelayerToRun]
	if !exist {
		return fmt.Errorf("not found config of RelayerToRun")
	}
	if cfg.RelayerToRun == config.TOP_CHAIN {
		for name, c := range cfg.RelayerConfig {
			logger.Info("name: ", name)
			if name == config.TOP_CHAIN {
				continue
			}
			if name != config.ETH_CHAIN && name != config.BSC_CHAIN && name != config.HECO_CHAIN {
				logger.Warn("TopRelayer not support:", name)
				continue
			}
			topRelayer, exist := topRelayers[name]
			if !exist {
				logger.Warn("unknown chain config:", name)
				continue
			}
			err := startOneRelayer(name, topRelayer, topConfig, c.Url, pass, wg)
			if err != nil {
				logger.Error("StartRelayer %v error: %v", name, err)
				continue
			}
		}
	} else {
		err := startOneRelayer(cfg.RelayerToRun, crossChainRelayer, RelayerConfig, topConfig.Url, pass, wg)
		if err != nil {
			logger.Error("StartRelayer error:", err)
			return err
		}
	}

	return nil
}
