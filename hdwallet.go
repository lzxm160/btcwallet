// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	_"encoding/hex"
	"fmt"
	"os"
	_"reflect"
	_"testing"

	"github.com/btcsuite/btcd/chaincfg"
	_"github.com/btcsuite/btcd/chaincfg/chainhash"
	_"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/walletdb"
	_"github.com/btcsuite/btcwallet/wallet"
)
const (
	// maxRecentHashes is the maximum number of hashes to keep in history
	// for the purposes of rollbacks.
	// maxRecentHashes = 20
	// latestMgrVersion = 4
)
var (
	cfg *config
	seed = []byte{
		0x2a, 0x64, 0xdf, 0x08, 0x5e, 0xef, 0xed, 0xd8, 0xbf,
		0xdb, 0xb3, 0x31, 0x76, 0xb5, 0xba, 0x2e, 0x62, 0xe8,
		0xbe, 0x8b, 0x56, 0xc8, 0x83, 0x77, 0x95, 0x59, 0x8b,
		0xb6, 0xc4, 0x40, 0xc0, 0x64,
	}

	pubPassphrase   = []byte("_DJr{fL4H0O}*-0\n:V1izc)(6BomK")
	privPassphrase  = []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")
	pubPassphrase2  = []byte("-0NV4P~VSJBWbunw}%<Z]fuGpbN[ZI")
	privPassphrase2 = []byte("~{<]08%6!-?2s<$(8$8:f(5[4/!/{Y")

	// fastScrypt are parameters used throughout the tests to speed up the
	// scrypt operations.
	fastScrypt = &waddrmgr.ScryptOptions{
		N: 16,
		R: 8,
		P: 1,
	}

	// waddrmgrNamespaceKey is the namespace key for the waddrmgr package.
	waddrmgrNamespaceKey = []byte("waddrmgrNamespace")
)
func checkManagerError(testName string, gotErr error, wantErrCode waddrmgr.ErrorCode) bool {
	merr, ok := gotErr.(waddrmgr.ManagerError)
	if !ok {
		fmt.Printf("%s: unexpected error type - got %T, want %T",
			testName, gotErr, waddrmgr.ManagerError{})
		return false
	}
	if merr.ErrorCode != wantErrCode {
		fmt.Printf("%s: unexpected error code - got %s (%s), want %s",
			testName, merr.ErrorCode, merr.Description, wantErrCode)
		return false
	}

	return true
}
type testContext struct {
	db           walletdb.DB
	manager      *waddrmgr.Manager
	account      uint32
	create       bool
	unlocked     bool
	watchingOnly bool
}
func createDbNamespace(dbPath string) (walletdb.DB, walletdb.Namespace, error) {
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		return nil, nil, err
	}

	namespace, err := db.Namespace(waddrmgrNamespaceKey)
	if err != nil {
		db.Close()
		return nil, nil, err
	}

	return db, namespace, nil
}
func testNewAccount(tc *testContext) bool {
	if tc.watchingOnly {
		// Creating new accounts in watching-only mode should return ErrWatchingOnly
		_, err := tc.manager.NewAccount("aaaa")
		if !checkManagerError("Create account in watching-only mode", err,
			waddrmgr.ErrWatchingOnly) {
			tc.manager.Close()
			return false
		}
		return true
	}
	// Creating new accounts when wallet is locked should return ErrLocked
	_, err := tc.manager.NewAccount("aaaa")
	if !checkManagerError("Create account when wallet is locked", err,
		waddrmgr.ErrLocked) {
		tc.manager.Close()
		return false
	}
	// Unlock the wallet to decrypt cointype keys required
	// to derive account keys
	if err := tc.manager.Unlock(privPassphrase); err != nil {
		fmt.Printf("Unlock: unexpected error: %v", err)
		return false
	}
	tc.unlocked = true

	testName := "acct-create"
	expectedAccount := tc.account + 1
	if !tc.create {
		// Create a new account in open mode
		testName = "acct-open"
		expectedAccount++
	}
	account, err := tc.manager.NewAccount(testName)
	fmt.Println(testName)
	if err != nil {
		fmt.Printf("NewAccount: unexpected error: %v", err)
		return false
	}
	if account != expectedAccount {
		fmt.Printf("NewAccount "+
			"account mismatch -- got %d, "+
			"want %d\n", account, expectedAccount)
		return false
	}

	// // Test duplicate account name error
	// _, err = tc.manager.NewAccount(testName)
	// wantErrCode := waddrmgr.ErrDuplicateAccount
	// if !checkManagerError(testName, err, wantErrCode) {
	// 	return false
	// }
	// // Test account name validation
	// testName = "" // Empty account names are not allowed
	// _, err = tc.manager.NewAccount(testName)
	// wantErrCode = waddrmgr.ErrInvalidAccount
	// if !checkManagerError(testName, err, wantErrCode) {
	// 	return false
	// }
	// testName = "imported" // A reserved account name
	// _, err = tc.manager.NewAccount(testName)
	// wantErrCode = waddrmgr.ErrInvalidAccount
	// if !checkManagerError(testName, err, wantErrCode) {
	// 	return false
	// }
	return true
}
func testNamePrefix(tc *testContext) string {
	prefix := "Open "
	if tc.create {
		prefix = "Create "
	}

	return prefix + fmt.Sprintf("account #%d", tc.account)
}
func testForEachAccountAddress(tc *testContext) bool {
	prefix := testNamePrefix(tc) + " testForEachAccountAddress"
	{
		// addrs, _ := tc.manager.NextInternalAddresses(tc.account, 5)
		// for i := 0; i < len(addrs); i++ {
		// 	gotAddr := addrs[i]
		// 	fmt.Printf("prefix:%s,internal addr:%s\n",prefix,gotAddr.Address().String())
		// }
	}
	var addrs []waddrmgr.ManagedAddress
	err := tc.manager.ForEachAccountAddress(tc.account,
		func(maddr waddrmgr.ManagedAddress) error {
			addrs = append(addrs, maddr)
			return nil
		})
	if err != nil {
		fmt.Printf("%s: unexpected error: %v", prefix, err)
		return false
	}

	{
		// addrs, _ := tc.manager.NextInternalAddresses(tc.account, 5)
		for i := 0; i < len(addrs); i++ {
			gotAddr := addrs[i]
			fmt.Printf("prefix:%s,account addr:%s\n",prefix,gotAddr.Address().String())
		}
	}
	
	fmt.Println("-------------------------------------")
	{
		// addrs, _ := tc.manager.NextExternalAddresses(tc.account, 5)
		// for i := 0; i < len(addrs); i++ {
		// 	gotAddr := addrs[i]
		// 	fmt.Printf("prefix:%s,external addr:%s\n",prefix,gotAddr.Address().String())
		// }
	}



	return true
}
func testLookupAccount(tc *testContext) bool {
	
		account, err := tc.manager.LookupAccount("aaaa")
		if err != nil {
			fmt.Printf("LookupAccount: unexpected error: %v", err)
			return false
		}
		fmt.Printf("test account:%d\n",account)
	// Test last account
	lastAccount, err := tc.manager.LastAccount()
	var expectedLastAccount uint32
	expectedLastAccount = 1
	if !tc.create {
		// Existing wallet manager will have 3 accounts
		expectedLastAccount = 2
	}
	if lastAccount != expectedLastAccount {
		fmt.Printf("LookupAccount "+
			"account mismatch -- got %d, "+
			"want %d", lastAccount, expectedLastAccount)
		return false
	}

	// // Test account lookup for default account adddress
	// var expectedAccount uint32
	// for i, addr := range expectedAddrs {
	// 	addr, err := btcutil.NewAddressPubKeyHash(addr.addressHash,
	// 		tc.manager.ChainParams())
	// 	if err != nil {
	// 		tc.t.Errorf("AddrAccount #%d: unexpected error: %v", i, err)
	// 		return false
	// 	}
	// 	account, err := tc.manager.AddrAccount(addr)
	// 	if err != nil {
	// 		tc.t.Errorf("AddrAccount #%d: unexpected error: %v", i, err)
	// 		return false
	// 	}
	// 	if account != expectedAccount {
	// 		tc.t.Errorf("AddrAccount "+
	// 			"account mismatch -- got %d, "+
	// 			"want %d", account, expectedAccount)
	// 		return false
	// 	}
	// }
	return true
}
func testManagerAPI(tc *testContext) {
	// testLocking(tc)
	// testExternalAddresses(tc)
	// testInternalAddresses(tc)
	// testImportPrivateKey(tc)
	// testImportScript(tc)
	// testMarkUsed(tc)
	// testChangePassphrase(tc)

	// Reset default account
	// tc.account = 0
	testNewAccount(tc)
	// tc.account = 1
	testLookupAccount(tc)
	// testForEachAccount(tc)
	fmt.Println("=====================================")
	tc.account = 1
	testForEachAccountAddress(tc)

	// Rename account 1 "acct-create"
	
	// testRenameAccount(tc)
}
func test() {
	dbName := "mgrtest.bin"
	// _ = os.Remove(dbName)
	db, mgrNamespace, err := createDbNamespace(dbName)
	if err != nil {
		fmt.Printf("createDbNamespace: unexpected error: %v", err)
		return
	}
	// defer os.Remove(dbName)
	defer db.Close()

	// Open manager that does not exist to ensure the expected error is
	// returned.
	// _, err = waddrmgr.Open(mgrNamespace, pubPassphrase,
	// 	&chaincfg.MainNetParams, nil)
	// if !checkManagerError("Open non-existant", err, waddrmgr.ErrNoExist) {
	// 	return
	// }

	// Create a new manager.
	// err = waddrmgr.Create(mgrNamespace, seed, pubPassphrase,
	// 	privPassphrase, &chaincfg.MainNetParams, fastScrypt)
	// if err != nil {
	// 	fmt.Printf("Create: unexpected error: %v", err)
	// 	return
	// }
	mgr, err := waddrmgr.Open(mgrNamespace, pubPassphrase,
		&chaincfg.MainNetParams, nil)
	if err != nil {
		fmt.Printf("Open: unexpected error: %v", err)
		return
	}

	// NOTE: Not using deferred close here since part of the tests is
	// explicitly closing the manager and then opening the existing one.

	// Attempt to create the manager again to ensure the expected error is
	// returned.
	// err = waddrmgr.Create(mgrNamespace, seed, pubPassphrase,
	// 	privPassphrase, &chaincfg.MainNetParams, fastScrypt)
	// if !checkManagerError("Create existing", err, waddrmgr.ErrAlreadyExists) {
	// 	mgr.Close()
	// 	return
	// }

	// Run all of the manager API tests in create mode and close the
	// manager after they've completed
	// testManagerAPI(&testContext{
	// 	db:           db,
	// 	manager:      mgr,
	// 	account:      0,
	// 	create:       true,
	// 	watchingOnly: false,
	// })
	// mgr.Close()

	// Ensure the expected error is returned if the latest manager version
	// constant is bumped without writing code to actually do the upgrade.
	// *waddrmgr.TstLatestMgrVersion++
	// _, err = waddrmgr.Open(mgrNamespace, pubPassphrase,
	// 	&chaincfg.MainNetParams, nil)
	// if !checkManagerError("Upgrade needed", err, waddrmgr.ErrUpgrade) {
	// 	return
	// }
	// *waddrmgr.TstLatestMgrVersion--

	// Open the manager and run all the tests again in open mode which
	// avoids reinserting new addresses like the create mode tests do.
	// mgr, err = waddrmgr.Open(mgrNamespace, pubPassphrase,
	// 	&chaincfg.MainNetParams, nil)
	// if err != nil {
	// 	fmt.Printf("Open: unexpected error: %v", err)
	// 	return
	// }
	defer mgr.Close()

	tc := &testContext{
		db:           db,
		manager:      mgr,
		account:      0,
		create:       false,
		watchingOnly: false,
	}
	testManagerAPI(tc)

	// Unlock the manager so it can be closed with it unlocked to ensure
	// it works without issue.
	if err := mgr.Unlock(privPassphrase); err != nil {
		fmt.Printf("Unlock: unexpected error: %v", err)
	}
}
func main() {
	if err := walletMainx(); err != nil {
		os.Exit(1)
	}
}

func walletMainx() error {
	test()
	// tcfg, _, err := loadConfig()
	// if err != nil {
	// 	return err
	// }
	// cfg = tcfg
	// defer func() {
	// 	if logRotator != nil {
	// 		logRotator.Close()
	// 	}
	// }()

	// // Show version at startup.
	// log.Infof("Version %s", version())

	// dbDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	// loader := wallet.NewLoader(activeNet.Params, dbDir)

	// if !cfg.NoInitialLoad {
	// 	// Load the wallet database.  It must have been created already
	// 	// or this will return an appropriate error.
	// 	_, err = loader.OpenExistingWallet([]byte(cfg.WalletPass), true)
	// 	if err != nil {
	// 		log.Error(err)
	// 		return err
	// 	}
	// }

	// log.Info("Shutdown complete")
	return nil
}
