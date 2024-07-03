package main

//go:generate abigen --sol cpuHeavyTC/CPUHeavy.sol --pkg cpuHeavyTC --out cpuHeavyTC/CPUHeavy.go
//go:generate abigen --sol userStorageTC/UserStorage.sol --pkg userStorageTC --out userStorageTC/UserStorage.go

import (
	"context"
	"flag"
	"fmt"
	"github.com/klaytn/klaytn-load-tester/tasks"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/newEthereumAccessListTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/newEthereumDynamicFeeTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/newFeeDelegatedSmartContractExecutionTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/newFeeDelegatedSmartContractExecutionWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/newSmartContractExecutionTC"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/klaytn/klaytn-load-tester/account"
	"github.com/klaytn/klaytn-load-tester/testcases/erc20TransferTC"
	"github.com/klaytn/klaytn-load-tester/testcases/erc721TransferTC"
	"github.com/klaytn/klaytn-load-tester/testcases/storageTrieWriteTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/ethereumTxAccessListTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/ethereumTxDynamicFeeTC"
	"github.com/klaytn/klaytn-load-tester/testcases/txtypeTCs/ethereumTxLegacyTC"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	klay "github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/myzhan/boomer"
)

// sets build options from ldflags.
var (
	Version   = "1.0.0"
	Commit    string
	Branch    string
	Tag       string
	BuildDate string
	BuildUser string
)

type locustTestCase interface {
	init()
	run()
}

var (
	coinbasePrivatekey = ""
	gCli               *klay.Client
	gEndpoint          string

	coinbase    *account.Account
	newCoinbase *account.Account

	nUserForUnsigned    = 5 //number of virtual user account for unsigned tx
	accGrpForUnsignedTx []*account.Account

	nUserForSigned    = 5
	accGrpForSignedTx []*account.Account

	nUserForNewAccounts  = 5
	accGrpForNewAccounts []*account.Account

	activeUserPercent = 100

	SmartContractAccount *account.Account

	tcStr     string
	tcStrList []string

	chargeValue *big.Int

	gasPrice *big.Int
	baseFee  *big.Int
)

func Create(endpoint string) *klay.Client {
	c, err := klay.Dial(endpoint)
	if err != nil {
		log.Fatalf("Failed to connect RPC: %v", err)
	}
	return c
}

func inTheTCList(tcName string) bool {
	for _, tc := range tcStrList {
		if tcName == tc {
			return true
		}
	}
	return false
}

// Dedicated and fixed private key used to deploy a smart contract for ERC20 and ERC721 value transfer performance test.
var ERC20DeployPrivateKeyStr = "eb2c84d41c639178ff26a81f488c196584d678bb1390cc20a3aeb536f3969a98"
var ERC721DeployPrivateKeyStr = "45c40d95c9b7898a21e073b5bf952bcb05f2e70072e239a8bbd87bb74a53355e"

// prepareERC20Transfer sets up ERC20 transfer performance test.
func prepareERC20Transfer(accGrp []*account.Account) {
	if !inTheTCList(erc20TransferTC.Name) {
		return
	}
	erc20DeployAcc := account.GetAccountFromKey(0, ERC20DeployPrivateKeyStr)
	log.Printf("prepareERC20Transfer", "addr", erc20DeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts([]*account.Account{erc20DeployAcc})

	// A smart contract for ERC20 value transfer performance TC.
	erc20TransferTC.SmartContractAccount = deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.DeployERC20, "ERC20 Performance Test Contract")
	newCoinBaseAccountMap := map[common.Address]*account.Account{newCoinbase.GetAddress(): newCoinbase}
	firstChargeTokenToTestAccounts(newCoinBaseAccountMap, erc20TransferTC.SmartContractAccount.GetAddress(), erc20DeployAcc.TransferERC20, big.NewInt(1e11))

	chargeTokenToTestAccounts(accGrp, erc20TransferTC.SmartContractAccount.GetAddress(), newCoinbase.TransferERC20, big.NewInt(1e4))
}

// prepareERC721Transfer sets up ERC721 transfer performance test.
func prepareERC721Transfer(accGrp []*account.Account) {
	if !inTheTCList(erc721TransferTC.Name) {
		return
	}
	erc721DeployAcc := account.GetAccountFromKey(0, ERC721DeployPrivateKeyStr)
	log.Printf("prepareERC721Transfer", "addr", erc721DeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts([]*account.Account{erc721DeployAcc})

	// A smart contract for ERC721 value transfer performance TC.
	erc721TransferTC.SmartContractAccount = deploySingleSmartContract(erc721DeployAcc, erc721DeployAcc.DeployERC721, "ERC721 Performance Test Contract")

	// Wait for reward tester to get started
	time.Sleep(30 * time.Second)
	newCoinbase.MintERC721ToTestAccounts(gCli, accGrp, erc721TransferTC.SmartContractAccount.GetAddress(), 5)
	log.Println("MintERC721ToTestAccounts", "len(accGrp)", len(accGrp))
}

// Dedicated and fixed private key used to deploy a smart contract for storage trie write performance test.
var storageTrieDeployPrivateKeyStr = "3737c381633deaaa4c0bdbc64728f6ef7d381b17e1d30bbb74665839cec942b8"

// prepareStorageTrieWritePerformance sets up ERC20 storage trie write performance test.
func prepareStorageTrieWritePerformance() {
	if !inTheTCList(storageTrieWriteTC.Name) {
		return
	}
	storageTrieDeployAcc := account.GetAccountFromKey(0, storageTrieDeployPrivateKeyStr)
	log.Printf("prepareStorageTrieWritePerformance", "addr", storageTrieDeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts([]*account.Account{storageTrieDeployAcc})

	// A smart contract for storage trie store performance TC.
	storageTrieWriteTC.SmartContractAccount = deploySingleSmartContract(storageTrieDeployAcc, storageTrieDeployAcc.DeployStorageTrieWrite, "Storage Trie Performance Test Contract")
}

func prepareTestAccountsAndContracts(accGrp []*account.Account) {
	// First, charging KLAY to the test accounts.
	chargeKLAYToTestAccounts(accGrp)

	// Second, deploy contracts used for some TCs.
	// If the test case is not on the list, corresponding contract won't be deployed.
	prepareERC20Transfer(accGrp)
	prepareStorageTrieWritePerformance()

	// Third, deploy contracts for general tests.
	// A smart contract for general smart contract related TCs.
	GeneralSmartContract := deploySmartContract(newCoinbase.TransferNewSmartContractDeployTxHumanReadable, "General Purpose Test Smart Contract")
	newSmartContractExecutionTC.SmartContractAccount = GeneralSmartContract
	newFeeDelegatedSmartContractExecutionTC.SmartContractAccount = GeneralSmartContract
	newFeeDelegatedSmartContractExecutionWithRatioTC.SmartContractAccount = GeneralSmartContract
	ethereumTxLegacyTC.SmartContractAccount = GeneralSmartContract
	ethereumTxAccessListTC.SmartContractAccount = GeneralSmartContract
	ethereumTxDynamicFeeTC.SmartContractAccount = GeneralSmartContract
	newEthereumAccessListTC.SmartContractAccount = GeneralSmartContract
	newEthereumDynamicFeeTC.SmartContractAccount = GeneralSmartContract
}

func chargeKLAYToTestAccounts(accGrp []*account.Account) {
	log.Printf("Start charging KLAY to test accounts")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, acc := range accGrp {
		for {
			_, _, err := newCoinbase.TransferSignedTxReturnTx(true, gCli, acc, chargeValue)
			if err == nil {
				break // Success, move to next account.
			}
			numChargedAcc, lastFailedNum = estimateRemainingTime(accGrp, numChargedAcc, lastFailedNum)
		}
		numChargedAcc++
	}

	log.Printf("Finished charging KLAY to %d test account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

type tokenChargeFunc func(initialCharge bool, c *client.Client, tokenContractAddr common.Address, recipient *account.Account, value *big.Int) (*types.Transaction, *big.Int, error)

// firstChargeTokenToTestAccounts charges initially generated tokens to newCoinbase account for further testing.
// As this work is done simultaneously by different slaves, this should be done in "try and check" manner.
func firstChargeTokenToTestAccounts(accGrp map[common.Address]*account.Account, tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int) {
	log.Printf("Start initial token charging to new coinbase")

	numChargedAcc := 0
	for _, recipientAccount := range accGrp {
		for {
			tx, _, err := tokenChargeFn(true, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			for err != nil {
				log.Printf("Failed to execute %s: err %s", tx.Hash().String(), err.Error())
				time.Sleep(1 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
				tx, _, err = tokenChargeFn(true, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			}
			ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			receipt, err := bind.WaitMined(ctx, gCli, tx)
			cancelFn()
			if receipt != nil {
				break
			}
		}
		numChargedAcc++
	}

	log.Printf("Finished initial token charging to %d new coinbase account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

// chargeTokenToTestAccounts charges default token to the test accounts for testing.
// As it is done independently among the slaves, it has simpler logic than firstChargeTokenToTestAccounts.
func chargeTokenToTestAccounts(accGrp []*account.Account, tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int) {
	log.Printf("Start charging tokens to test accounts")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, recipientAccount := range accGrp {
		for {
			_, _, err := tokenChargeFn(false, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			if err == nil {
				break // Success, move to next account.
			}
			numChargedAcc, lastFailedNum = estimateRemainingTime(accGrp, numChargedAcc, lastFailedNum)
		}
		numChargedAcc++
	}

	log.Printf("Finished charging tokens to %d test account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

func estimateRemainingTime(accGrp []*account.Account, numChargedAcc, lastFailedNum int) (int, int) {
	if lastFailedNum > 0 {
		// Not 1st failed cases.
		TPS := (numChargedAcc - lastFailedNum) / 5 // TPS of only this slave during `txpool is full` situation.
		lastFailedNum = numChargedAcc

		if TPS <= 5 {
			log.Printf("Retry to charge test account #%d. But it is too slow. %d TPS\n", numChargedAcc, TPS)
		} else {
			remainTime := (len(accGrp) - numChargedAcc) / TPS
			remainHour := remainTime / 3600
			remainMinute := (remainTime % 3600) / 60

			log.Printf("Retry to charge test account #%d. Estimated remaining time: %d hours %d mins later\n", numChargedAcc, remainHour, remainMinute)
		}
	} else {
		// 1st failed case.
		lastFailedNum = numChargedAcc
		log.Printf("Retry to charge test account #%d.\n", numChargedAcc)
	}
	time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
	return numChargedAcc, lastFailedNum
}

type contractDeployFunc func(c *client.Client, to *account.Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error)

// deploySmartContract deploys smart contracts by the number of locust slaves.
// In other words, each slave owns its own contract for testing.
func deploySmartContract(contractDeployFn contractDeployFunc, contractName string) *account.Account {
	addr, lastTx, _, err := contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	for err != nil {
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), gCli, lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

// deploySingleSmartContract deploys only one smart contract among the slaves.
// It the contract is already deployed by other slave, it just calculates the address of the contract.
func deploySingleSmartContract(erc20DeployAcc *account.Account, contractDeployFn contractDeployFunc, contractName string) *account.Account {
	addr, lastTx, _, err := contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	for err != nil {
		if err == account.AlreadyDeployedErr {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		if strings.HasPrefix(err.Error(), "known transaction") {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), gCli, lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

func prepareAccounts() {
	totalChargeValue := new(big.Int)
	totalChargeValue.Mul(chargeValue, big.NewInt(int64(nUserForUnsigned+nUserForSigned+nUserForNewAccounts+1)))

	// Import coinbase Account
	coinbase = account.GetAccountFromKey(0, coinbasePrivatekey)
	newCoinbase = account.NewAccount(0)

	if len(chargeValue.Bits()) != 0 {
		for {
			coinbase.GetNonceFromBlock(gCli)
			hash, _, err := coinbase.TransferSignedTx(gCli, newCoinbase, totalChargeValue)
			if err != nil {
				log.Printf("%v: charge newCoinbase fail: %v\n", os.Getpid(), err)
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			log.Printf("%v : charge newCoinbase: %v, Txhash=%v\n", os.Getpid(), newCoinbase.GetAddress().String(), hash.String())

			getReceipt := false
			// After this loop waiting for 10 sec, It will retry to charge with new nonce.
			// it means another node stole the nonce.
			for i := 0; i < 5; i++ {
				time.Sleep(2000 * time.Millisecond)
				ctx := context.Background()

				//_, err := gCli.TransactionReceipt(ctx, hash)
				//if err != nil {
				//	getReceipt = true
				//	log.Printf("%v : charge newCoinbase success: %v\n", os.Getpid(), newCoinbase.GetAddress().String())
				//	break
				//}
				//log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), newCoinbase.GetAddress().String())

				val, err := gCli.BalanceAt(ctx, newCoinbase.GetAddress(), nil)
				if err == nil {
					if val.Cmp(big.NewInt(0)) == 1 {
						getReceipt = true
						log.Printf("%v : charge newCoinbase success: %v, balance=%v peb\n", os.Getpid(), newCoinbase.GetAddress().String(), val.String())
						break
					}
					log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), newCoinbase.GetAddress().String())
				} else {
					log.Printf("%v : check balance err: %v\n", os.Getpid(), err)
				}
			}

			if getReceipt {
				break
			}
		}
	}

	println("Unsigned Account Group Preparation...")
	//bar := pb.StartNew(nUserForUnsigned)

	// Create test account pool
	for i := 0; i < nUserForUnsigned; i++ {
		accGrpForUnsignedTx = append(accGrpForUnsignedTx, account.NewAccount(i))
		fmt.Printf("%v\n", accGrpForUnsignedTx[i].GetAddress().String())
		//bar.Increment()
	}
	//bar.Finish()	//bar.FinishPrint("Completed.")
	//
	println("Signed Account Group Preparation...")
	//bar = pb.StartNew(nUserForSigned)

	for i := 0; i < nUserForSigned; i++ {
		accGrpForSignedTx = append(accGrpForSignedTx, account.NewAccount(i))
		fmt.Printf("%v\n", accGrpForSignedTx[i].GetAddress().String())
		//bar.Increment()
	}

	println("New account group preparation...")
	for i := 0; i < nUserForNewAccounts; i++ {
		accGrpForNewAccounts = append(accGrpForNewAccounts, account.NewKlaytnAccount(i))
	}
}

func initArgs(tcNames string) {
	chargeKLAYAmount := 1000000000
	gEndpointPtr := flag.String("endpoint", "http://localhost:8545", "Target EndPoint")
	nUserForSignedPtr := flag.Int("vusigned", nUserForSigned, "num of test account for signed Tx TC")
	nUserForUnsignedPtr := flag.Int("vuunsigned", nUserForUnsigned, "num of test account for unsigned Tx TC")
	activeUserPercentPtr := flag.Int("activepercent", activeUserPercent, "percent of active accounts")
	keyPtr := flag.String("key", "", "privatekey of coinbase")
	chargeKLAYAmountPtr := flag.Int("charge", chargeKLAYAmount, "charging amount for each test account in KLAY")
	versionPtr := flag.Bool("version", false, "show version number")
	httpMaxIdleConnsPtr := flag.Int("http.maxidleconns", 100, "maximum number of idle connections in default http client")
	flag.StringVar(&tcStr, "tc", tcNames, "tasks which user want to run, multiple tasks are separated by comma.")

	flag.Parse()

	if *versionPtr || (len(os.Args) >= 2 && os.Args[1] == "version") {
		printVersion()
		os.Exit(0)
	}

	if *keyPtr == "" {
		log.Fatal("key argument is not defined. You should set the key for coinbase.\n example) klaytc -key='2ef07640fd8d3f568c23185799ee92e0154bf08ccfe5c509466d1d40baca3430'")
	}

	// setup default http client.
	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.MaxIdleConns = *httpMaxIdleConnsPtr
		tr.MaxIdleConnsPerHost = *httpMaxIdleConnsPtr
	}

	// for TC Selection
	if tcStr != "" {
		// Run tasks without connecting to the master.
		tcStrList = strings.Split(tcStr, ",")
	}

	gEndpoint = *gEndpointPtr

	nUserForSigned = *nUserForSignedPtr
	nUserForUnsigned = *nUserForUnsignedPtr
	activeUserPercent = *activeUserPercentPtr
	coinbasePrivatekey = *keyPtr
	chargeKLAYAmount = *chargeKLAYAmountPtr
	chargeValue = new(big.Int)
	chargeValue.Set(new(big.Int).Mul(big.NewInt(int64(chargeKLAYAmount)), big.NewInt(params.KLAY)))

	fmt.Println("Arguments are set like the following:")
	fmt.Printf("- Target EndPoint = %v\n", gEndpoint)
	fmt.Printf("- nUserForSigned = %v\n", nUserForSigned)
	fmt.Printf("- nUserForUnsigned = %v\n", nUserForUnsigned)
	fmt.Printf("- activeUserPercent = %v\n", activeUserPercent)
	fmt.Printf("- coinbasePrivatekey = %v\n", coinbasePrivatekey)
	fmt.Printf("- charging KLAY Amount = %v\n", chargeKLAYAmount)
	fmt.Printf("- tc = %v\n", tcStr)
}

func updateChainID() {
	fmt.Println("Updating ChainID from RPC")
	for {
		ctx := context.Background()
		chainID, err := gCli.ChainID(ctx)

		if err == nil {
			fmt.Println("chainID :", chainID)
			account.SetChainID(chainID)
			break
		}
		fmt.Println("Retrying updating chainID... ERR: ", err)

		time.Sleep(2 * time.Second)
	}
}

func updateGasPrice() {
	// TODO: refactor to updating gasPrice with goverance.magma.upperboundbasefee
	gasPrice = big.NewInt(750000000000)

	/* Deprecated because of KIP-71 hardfork
	fmt.Println("Updating GasPrice from RPC")
	for {
		ctx := context.Background()
		gp, err := gCli.SuggestGasPrice(ctx)

		if err == nil {
			gasPrice = gp
			fmt.Println("gas price :", gasPrice.String())
			break
		}
		fmt.Println("Retrying updating GasPrice... ERR: ", err)

		time.Sleep(2 * time.Second)
	}
	*/
	account.SetGasPrice(gasPrice)
}

func updateBaseFee() {
	baseFee = big.NewInt(0)
	// TODO: Uncomment below when klaytn 1.8.0 is released.
	//for {
	//	ctx := context.Background()
	//	h, err := gCli.HeaderByNumber(ctx, nil)
	//
	//	if err == nil {
	//		baseFee = h.BaseFee
	//		fmt.Println("base fee :", baseFee.String())
	//		break
	//	}
	//	fmt.Println("Retrying updating BaseFee... ERR: ", err)
	//
	//	time.Sleep(2 * time.Second)
	//}
	account.SetBaseFee(baseFee)
}

func setRLimit(resourceType int, val uint64) error {
	if runtime.GOOS == "darwin" {
		return nil
	}
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(resourceType, &rLimit)
	if err != nil {
		return err
	}
	rLimit.Cur = val
	err = syscall.Setrlimit(resourceType, &rLimit)
	if err != nil {
		return err
	}
	return nil
}

func printVersion() {
	version := Version
	if len(Commit) >= 7 {
		version += "-" + Commit[:7]
	}
	if Tag != "" && Tag != "undefined" {
		version = Tag
	}
	fmt.Printf("Version :\t%s\n", version)
	fmt.Printf("git.Branch :\t%s\n", Branch)
	fmt.Printf("git.Commit :\t%s\n", Commit)
	fmt.Printf("git.Tag :\t%s\n", Tag)
	fmt.Printf("build.Date :\t%s\n", BuildDate)
	fmt.Printf("build.User :\t%s\n", BuildUser)
}

func main() {
	// Call initTCList to get all TC names
	taskSet := initTCList()

	var tcNames string
	for i, task := range taskSet {
		if i != 0 {
			tcNames += ","
		}
		tcNames += task.Name
	}

	initArgs(tcNames)

	// Create Cli pool
	gCli = Create(gEndpoint)

	// Update chainID
	updateChainID()

	// Update gasPrice
	updateGasPrice()

	gasPrice = big.NewInt(750000000000)

	// Update baseFee
	updateBaseFee()

	// Set coinbase & Create Test Account
	prepareAccounts()

	filteredTask := tasks.NewFilteredTasks(tcStrList, accGrpForUnsignedTx, gEndpoint)

	println("Adding tasks")

	accGrp := append(accGrpForSignedTx, accGrpForUnsignedTx...)
	if len(chargeValue.Bits()) != 0 {
		prepareTestAccountsAndContracts(accGrp)
	}
	// After charging accounts, cut the slice to the desired length, calculated by ActiveAccountPercent.
	if activeUserPercent > 100 {
		log.Fatalf("ActiveAccountPercent should be less than or equal to 100, but it is %v", activeUserPercent)
	}
	numActiveAccounts := len(accGrp) * activeUserPercent / 100
	// Not to assign 0 account for some cases.
	if numActiveAccounts == 0 {
		numActiveAccounts = 1
	}
	accGrp = accGrp[:numActiveAccounts]
	prepareERC721Transfer(accGrp)

	if len(filteredTask) == 0 {
		log.Fatal("No Tc is set. Please set TcList. \nExample argument) -tc='" + tcNames + "'")
	}

	println("Initializing tasks")
	var filteredBoomerTask []*boomer.Task
	for _, task := range filteredTask {
		task.Init(accGrp, gEndpoint, gasPrice)
		filteredBoomerTask = append(filteredBoomerTask, &boomer.Task{10, task.Fn, task.Name})
		println("=> " + task.Name + " task is initialized.")
	}

	setRLimit(syscall.RLIMIT_NOFILE, 1024*400)

	// Locust Slave Run
	// 봐봐, task에는 이름이 있지 이름 list를 만들어보자.
	// 결국 filteredBoomerTask 가 만들어지는게 목적이네.
	boomer.Run(filteredBoomerTask...)
	//boomer.Run(cpuHeavyTx)
}
