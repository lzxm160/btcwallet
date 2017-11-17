// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	_ "encoding/hex"
	"fmt"
	"os"
	_ "reflect"
	_ "testing"

	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/btcsuite/btcd/chaincfg/chainhash"
	_ "github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	_ "github.com/btcsuite/btcwallet/wallet"
	"github.com/btcsuite/btcwallet/walletdb"
)

const (
// maxRecentHashes is the maximum number of hashes to keep in history
// for the purposes of rollbacks.
// maxRecentHashes = 20
// latestMgrVersion = 4
)

var MainNetParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "8333",
	DNSSeeds: []DNSSeed{
		{"seed.bitcoin.sipa.be", true},
		{"dnsseed.bluematt.me", true},
		{"dnsseed.bitcoin.dashjr.org", false},
		{"seed.bitcoinstats.com", true},
		{"seed.bitnodes.io", false},
		{"seed.bitcoin.jonasschnelli.ch", true},
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1d00ffff,
	BIP0034Height:            227931, // 000000000000024b89b42a942fe0d9fea3bb44ab7bd1b19115dd6a759c0808b8
	BIP0065Height:            388381, // 000000000000000004c2b624ed5d7756c508d90fd0da2c7c679febfa6c4735f0
	BIP0066Height:            363725, // 00000000000000000379eaa19dce8c9b722d46ae6a57c2f1a988119488b50931
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{11111, newHashFromStr("0000000069e244f73d78e8fd29ba2fd2ed618bd6fa2ee92559f542fdb26e7c1d")},
		{33333, newHashFromStr("000000002dd5588a74784eaa7ab0507a18ad16a236e7b1ce69f00d7ddfb5d0a6")},
		{74000, newHashFromStr("0000000000573993a3c9e41ce34471c079dcf5f52a0e824a81e7f953b8661a20")},
		{105000, newHashFromStr("00000000000291ce28027faea320c8d2b054b2e0fe44a773f3eefb151d6bdc97")},
		{134444, newHashFromStr("00000000000005b12ffd4cd315cd34ffd4a594f430ac814c91184a0d42d2b0fe")},
		{168000, newHashFromStr("000000000000099e61ea72015e79632f216fe6cb33d7899acb35b75c8303b763")},
		{193000, newHashFromStr("000000000000059f452a5f7340de6682a977387c17010ff6e6c3bd83ca8b1317")},
		{210000, newHashFromStr("000000000000048b95347e83192f69cf0366076336c639f9b7228e9ba171342e")},
		{216116, newHashFromStr("00000000000001b4f4b433e81ee46494af945cf96014816a4e2370f11b23df4e")},
		{225430, newHashFromStr("00000000000001c108384350f74090433e7fcf79a606b8e797f065b130575932")},
		{250000, newHashFromStr("000000000000003887df1f29024b06fc2200b55f8af8f35453d7be294df2d214")},
		{267300, newHashFromStr("000000000000000a83fbd660e918f218bf37edd92b748ad940483c7c116179ac")},
		{279000, newHashFromStr("0000000000000001ae8c72a0b0c301f67e3afca10e819efa9041e458e9bd7e40")},
		{300255, newHashFromStr("0000000000000000162804527c6e9b9f0563a280525f9d08c12041def0a0f3b2")},
		{319400, newHashFromStr("000000000000000021c6052e9becade189495d1c539aa37c58917305fd15f13b")},
		{343185, newHashFromStr("0000000000000000072b8bf361d01a6ba7d445dd024203fafc78768ed4368554")},
		{352940, newHashFromStr("000000000000000010755df42dba556bb72be6a32f3ce0b6941ce4430152c9ff")},
		{382320, newHashFromStr("00000000000000000a8dc6ed5b133d0eb2fd6af56203e4159789b092defd8ab2")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1479168000, // November 15, 2016 UTC
			ExpireTime: 1510704000, // November 15, 2017 UTC.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "bc", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x00, // starts with 1
	ScriptHashAddrID:        0x05, // starts with 3
	PrivateKeyID:            0x80, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 60,
}
var (
	cfg  *config
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
	// if tc.watchingOnly {
	// 	// Creating new accounts in watching-only mode should return ErrWatchingOnly
	// 	_, err := tc.manager.NewAccount("aaaa")
	// 	if !checkManagerError("Create account in watching-only mode", err,
	// 		waddrmgr.ErrWatchingOnly) {
	// 		tc.manager.Close()
	// 		return false
	// 	}
	// 	return true
	// }
	// Creating new accounts when wallet is locked should return ErrLocked
	// _, err := tc.manager.NewAccount("aaaa")
	// if !checkManagerError("Create account when wallet is locked", err,
	// 	waddrmgr.ErrLocked) {
	// 	tc.manager.Close()
	// 	return false
	// }
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
	fmt.Printf("NewAccount:%s,%d\n", testName, account)
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
			fmt.Printf("prefix:%s,account addr:%s\n", prefix, gotAddr.Address().String())
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

	account, err := tc.manager.LookupAccount("acct-create")
	if err != nil {
		fmt.Printf("LookupAccount: unexpected error: %v", err)
		return false
	}
	fmt.Printf("test account:%d\n", account)
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
	fmt.Println("=====testNewAccount================================")
	testNewAccount(tc)
	fmt.Println("=====testLookupAccount================================")
	// tc.account = 1
	testLookupAccount(tc)
	// testForEachAccount(tc)
	fmt.Println("=========testForEachAccountAddress0============================")
	fmt.Printf("tc.account:%d\n", tc.account)
	testForEachAccountAddress(tc)
	fmt.Println("=========testForEachAccountAddress1============================")
	tc.account = 1
	fmt.Printf("tc.account:%d\n", tc.account)
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
	err = waddrmgr.Create(mgrNamespace, seed, pubPassphrase,
		privPassphrase, &MainNetParams, fastScrypt)
	if err != nil {
		fmt.Printf("Create: unexpected error: %v", err)
		//return
	}
	mgr, err := waddrmgr.Open(mgrNamespace, pubPassphrase,
		&MainNetParams, nil)
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
	//测试创建模式
	tc := &testContext{
		db:           db,
		manager:      mgr,
		account:      0,
		create:       true,
		watchingOnly: false,
	}
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
	//测试open模式
	// tc := &testContext{
	// 	db:           db,
	// 	manager:      mgr,
	// 	account:      0,
	// 	create:       false,
	// 	watchingOnly: false,
	// }
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
