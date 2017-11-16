// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"

	"github.com/btcsuite/btcwallet/chain"
	"github.com/btcsuite/btcwallet/rpc/legacyrpc"
	"github.com/btcsuite/btcwallet/wallet"
)

var (
	cfg *config
)

func main() {
	if err := walletMain(); err != nil {
		os.Exit(1)
	}
}

func walletMain() error {
	tcfg, _, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg
	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	// Show version at startup.
	log.Infof("Version %s", version())

	dbDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	loader := wallet.NewLoader(activeNet.Params, dbDir)

	if !cfg.NoInitialLoad {
		// Load the wallet database.  It must have been created already
		// or this will return an appropriate error.
		_, err = loader.OpenExistingWallet([]byte(cfg.WalletPass), true)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	log.Info("Shutdown complete")
	return nil
}
