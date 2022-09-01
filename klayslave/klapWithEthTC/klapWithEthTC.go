package klapWithEthTC

import (
	"encoding/hex"
	"github.com/klaytn/klaytn-load-tester/klayslave/klapTC"
	"log"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
)

var (
	endPoint string
	cliPool  clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	nAcc   int
	accGrp []*account.Account

	gasPrice       *big.Int
	executablePath string
)

func Init(accs []*account.Account, ep string, gp *big.Int) {
	mutex.Lock()
	defer mutex.Unlock()

	gasPrice = gp
	endPoint = ep
	cliCreate := func() interface{} {
		c, err := client.Dial(endPoint)
		if err != nil {
			log.Fatalf("Failed to connect RPC: %v", err)
		}
		return c
	}
	cliPool.Init(20, 300, cliCreate)

	for _, acc := range accs {
		accGrp = append(accGrp, acc)
	}
	nAcc = len(accGrp)

	// Path to executable file that call the message with eth client.
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	executablePath = exPath + "/ethTxGenerator"
}

func KlapBalanceOfCallWithEthClient() {
	cli := cliPool.Alloc().(*client.Client)

	fromAccount := klapTC.Klap_account_list[rand.Int()%3612]
	contractAddr := common.HexToAddress("0xd109065ee17e2dc20b3472a4d4fb5907bd687d09")
	data, err := klapTC.GenerateData(fromAccount)
	if err != nil {
		panic("failed to generate data field of klaytn argument")
	}
	gas := "1100000"
	value := "0"

	var callopts bind.CallOpts
	callopts.Pending = false
	callopts.From = fromAccount

	start := boomer.Now()
	// To test this, you need to update submodule and build executable file.
	// ./ethTxGenerator endpoint eth_call fromAddress toAddress gas gasPrice value data [blockNumber]
	cmd := exec.Command(executablePath, endPoint, "eth_call", fromAccount.String(), contractAddr.String(), gasPrice.String(), gas, value, hex.EncodeToString(data))
	elapsed := boomer.Now() - start

	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to create and call with eth client: %v", err)
	}

	strResult := string(result[:])
	if strings.Contains(strResult, "Error") {
		boomer.Events.Publish("request_failure", "http", "KlapWithEthTC to "+endPoint, elapsed, err.Error())
		cli.Close()
	} else {
		log.Printf("[TC] KlapWithEthTC: Failed to call eth_call, err=%v, from:%x\n", err, fromAccount)
		boomer.Events.Publish("request_success", "http", "KlapWithEthTC to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	}
}
func KlapAppCallWithEthClient() {
	cli := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc].GetAddress()

	idx := rand.Int() % 7
	contractAddr := klapTC.ContractAddrList[idx]
	data := klapTC.DataList[idx]
	gas := "50000000"
	value := "0"

	start := boomer.Now()
	// To test this, you need to update submodule and build executable file.
	// ./ethTxGenerator endpoint eth_call fromAddress toAddress gas gasPrice value data [blockNumber]
	cmd := exec.Command(executablePath, endPoint, "eth_call", fromAccount.String(), contractAddr.String(), gasPrice.String(), gas, value, hex.EncodeToString(data))
	elapsed := boomer.Now() - start

	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to create and call with eth client: %v", err)
	}

	strResult := string(result[:])
	if strings.Contains(strResult, "Error") {
		boomer.Events.Publish("request_success", "http", "KlapAppWithEthTC to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		log.Printf("[TC] KlapAppWithEthTC: Failed to call eth_call, err=%v, from:%x\n", err, fromAccount)
		boomer.Events.Publish("request_failure", "http", "KlapAppWithEthTC to "+endPoint, elapsed, err.Error())
		cli.Close()
	}
}
