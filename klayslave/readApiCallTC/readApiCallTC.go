package readApiCallTC

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/myzhan/boomer"
	"github.com/tidwall/gjson"
)

var contracts = []common.Address{common.HexToAddress("0xc6a2ad8cc6e4a7e08fc37cc5954be07d499e7654"),common.HexToAddress("0x422abb57e4bb7d46032852b884b7bb4cc4a39cc7"),common.HexToAddress("0x86b0f3dfddadb2e57b00a6d740f1a464f79179bf"),
	common.HexToAddress("0xcee8faf64bb97a73bb51e115aa89c17ffa8dd167"),common.HexToAddress("0xf7cbb00b62376f41ddaba34fe16248ed4d82d422"),common.HexToAddress("0x97b4e13114ce2c9bf289be1ffd1268be5b2ed7c2"),
	common.HexToAddress("0x2b5065d6049099295c68f5fcb97b8b0d3c354df7"),common.HexToAddress("0xbe7377db700664331beb28023cfbd46de079efac"),common.HexToAddress("0x9657fb399847d85a9c1a234ece9ca09d5c00f466")}

var (
	gasPrice *big.Int
	endPoint string
	cliPool  clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	nAcc   int
	accGrp []*account.Account

	latestBlockNumber *big.Int
	getLogsBlockNumber *big.Int
	count             uint64
)

func Init(accs []*account.Account, ep string, gp *big.Int) {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return
	}
	initialized = true

	count = 0
	gasPrice = gp
	endPoint = ep

	cliCreate := func() interface{} {
		c, err := rpc.Dial(endPoint)
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

	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)
	latestBlockNumber, err := cli.BlockNumber(ctx)
	getLogsBlockNumber = latestBlockNumber
	if err != nil {
		log.Fatalf("Failed to get latest block number: %v", err)
	}
}

func sendBoomerEvent(tcName string, logString string, elapsed int64, cli *rpc.Client, err error) {
	if err == nil {
		boomer.Events.Publish("request_success", "http", tcName+" to "+endPoint, elapsed, int64(10))
	} else {
		log.Printf("[TC] %s: %s, err=%v\n", tcName, logString, err)
		boomer.Events.Publish("request_failure", "http", tcName+" to "+endPoint, elapsed, err.Error())
	}
	cliPool.Free(cli)
}

func getRandomBlockNumber(cli *client.Client, ctx context.Context) *big.Int {
	mutex.Lock()
	defer mutex.Unlock()

	count %= 10000000
	if count%10000 == 0 {
		bn, err := cli.BlockNumber(ctx)
		if err != nil {
			log.Printf("Failed to update the current block number. err=%s\n", err)
		} else {
			log.Printf("Update the current block number. blockNumber=0x%s\n", bn.Text(16))
			latestBlockNumber.Set(bn)
		}
	}
	count += 1

	return big.NewInt(0).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), latestBlockNumber)
}

func GasPrice() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()
	gp, err := cli.SuggestGasPrice(ctx)
	elapsed := boomer.Now() - start
	if err == nil && gp.Cmp(gasPrice) != 0 {
		err = errors.New("wrong gas price: " + gp.String() + ", answer: " + gasPrice.String())
	}
	sendBoomerEvent("readGasPrice", "Failed to call klay_gasPrice", elapsed, rpcCli, err)
}

func BlockNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()

	bn, err := cli.BlockNumber(ctx)
	if err == nil && bn.Cmp(big.NewInt(0)) != 1 {
		err = errors.New("wrong block number: 0x" + bn.Text(16) + ", answer: smaller than 0")
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readBlockNumber", "Failed to call klay_blockNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := getRandomBlockNumber(cli, ctx)
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockByNumberLatest() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	bn, err := cli.BlockNumber(ctx)
	if err != nil {
		sendBoomerEvent("readGetBlockByNumberLatest", "Failed to get latest block number", 0, rpcCli, err)
		return
	}

	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, bn) //read the random block
	if err == nil && block.Header().Number.Cmp(bn) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + bn.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumberLatest", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber_1() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := big.NewInt(80771176) //size 1864312
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber_1", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber_2() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := big.NewInt(78771176) //size 927300
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber_2", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber_3() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := big.NewInt(78771189) //size 378906
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber_3", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber_4() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := big.NewInt(78771188) //size 156057
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber_4", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetBlockReceiptsLatest() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	bn, err := cli.BlockNumber(ctx)
	if err != nil {
		sendBoomerEvent("readBlockReceiptsLatest", "Failed to get the latest block number", 0, rpcCli, err)
		return
	}
	block, err := cli.BlockByNumber(ctx, bn) // read the latest block
	if err != nil {
		sendBoomerEvent("readBlockReceiptsLatest", "Failed to call klay_getBlockByNumber", 0, rpcCli, err)
		return
	}

	start := boomer.Now()
	var j json.RawMessage
	err = rpcCli.CallContext(ctx, &j, "klay_getBlockReceipts", block.Hash().String())
	elapsed := boomer.Now() - start
	sendBoomerEvent("readBlockReceiptsLatest", "Failed to call klay_getBlockReceipts", elapsed, rpcCli, err)
}

func GetTraceBlockLatest() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	bn, err := cli.BlockNumber(ctx)
	if err != nil {
		sendBoomerEvent("readGetTraceLatestBlock", "Failed to get the latest block number", 0, rpcCli, err)
		return
	}
	block, err := cli.BlockByNumber(ctx, bn) // read the latest block
	if err != nil {
		sendBoomerEvent("readGetTraceLatestBlock", "Failed to call klay_getBlockByNumber", 0, rpcCli, err)
		return
	}

	start := boomer.Now()
	var j json.RawMessage
	err = rpcCli.CallContext(ctx, &j, "debug_traceBlockByNumber", block.NumberU64(), map[string]interface{}{"tracer": "fastCallTracer"})
	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetTraceLatestBlock", "Failed to call debug_traceBlockByNumber", elapsed, rpcCli, err)
}

func GetAccount() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)

	var j json.RawMessage
	fromAccount := accGrp[rand.Int()%nAcc]
	start := boomer.Now()

	err := rpcCli.CallContext(ctx, &j, "klay_getAccount", fromAccount.GetAddress(), "latest")
	if err == nil {
		ret := gjson.Get(string(j), "accType").String()
		if ret != "1" {
			err = errors.New("wrong account type: " + ret + ", answer: 1")
		}
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetAccount", "Failed to call klay_getAccount", elapsed, rpcCli, err)
}

func GetBlockWithConsensusInfoByNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := getRandomBlockNumber(cli, ctx)
	start := boomer.Now()

	var j json.RawMessage
	err := rpcCli.CallContext(ctx, &j, "klay_getBlockWithConsensusInfoByNumber", "0x"+ansBN.Text(16))
	if err == nil {
		ret := gjson.Get(string(j), "number").String()
		if !strings.Contains(ret, "0x"+ansBN.Text(16)) {
			err = errors.New("wrong block: " + ret + ", answer: " + "0x" + ansBN.Text(16))
		}
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockWithConsensusInfoByNumber",
		"Failed to call klay_GetBlockWithConsensusInfoByNumber", elapsed, rpcCli, err)
}

func GetBlockWithConsensusInfoByNumberRange() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBNFrom := getRandomBlockNumber(cli, ctx)
	ansBNTo := ansBNFrom.Add(ansBNFrom,big.NewInt(5))
	if ansBNTo.Cmp(latestBlockNumber) == 1{
		ansBNTo = latestBlockNumber
	}

	start := boomer.Now()

	var j json.RawMessage
	err := rpcCli.CallContext(ctx, &j, "klay_getBlockWithConsensusInfoByNumberRange", "0x"+ansBNFrom.Text(16), "0x"+ansBNTo.Text(16))

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockWithConsensusInfoByNumberRange",
		"Failed to call klay_GetBlockWithConsensusInfoByNumberRange", elapsed, rpcCli, err)
}

func GetLogsLatest() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()
	bn, err := cli.BlockNumber(ctx)
	if err != nil {
		sendBoomerEvent("readGetLogsLatest", "Failed to get the latest block number", 0, rpcCli, err)
		return
	}
	filter := klaytn.FilterQuery{ FromBlock:bn, ToBlock:bn}
	_, err = cli.FilterLogs(ctx, filter)
	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetLogsLatest", "Failed to call klay_getLogs", elapsed, rpcCli, err)
}

// GetLogsHeavy assumes next case
// * The service collects the log once a day.
// * So, it calls the klay_getLogs for ranges.
func GetLogsHeavy() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()

	mutex.Lock()
	to := big.NewInt(getLogsBlockNumber.Int64())
	getLogsBlockNumber.Sub(getLogsBlockNumber, big.NewInt(100))
	from := getLogsBlockNumber
	mutex.Unlock()

	_, err := cli.FilterLogs(ctx, klaytn.FilterQuery{FromBlock:from, ToBlock:to})

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetLogsHeavy", "Failed to call klay_getLogs", elapsed, rpcCli, err)
}
