package erc721TransferTC

import (
	account2 "github.com/klaytn/klaytn-load-tester/account"
	"github.com/klaytn/klaytn-load-tester/clipool"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/klaytn/klaytn/client"
	"github.com/myzhan/boomer"
)

const Name = "erc721TransferTC"

var (
	endPoint string
	nAcc     int
	accGrp   []*account2.Account
	cliPool  clipool.ClientPool
	gasPrice *big.Int

	// multinode tester
	transferedValue *big.Int
	expectedFee     *big.Int

	fromAccount     *account2.Account
	prevBalanceFrom *big.Int

	toAccount     *account2.Account
	prevBalanceTo *big.Int

	SmartContractAccount *account2.Account
)

func Init(accs []*account2.Account, endpoint string, gp *big.Int) {
	gasPrice = gp

	endPoint = endpoint

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

	rand.Seed(time.Now().UnixNano())
}

func Run() {
	cli := cliPool.Alloc().(*client.Client)

	fromAcc := accGrp[rand.Intn(nAcc)]
	toAcc := accGrp[rand.Intn(nAcc)]

	// Get token ID from the channel
	// Here is an assumption that it won't be blocked by the channel
	// Although this go routine can be blocked, other can send a NFT to this account
	fromNFTs := account2.ERC721Ledger[fromAcc.GetAddress()]
	tokenId := <-fromNFTs

	start := boomer.Now()
	_, _, err := fromAcc.TransferERC721(false, cli, SmartContractAccount.GetAddress(), toAcc, tokenId)
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", "http", Name+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
		toNFTs := account2.ERC721Ledger[toAcc.GetAddress()]
		toNFTs <- tokenId // push the token to the new owner's queue, it it does not fail

	} else {
		boomer.Events.Publish("request_failure", "http", Name+" to "+endPoint, elapsed, err.Error())
		fromNFTs <- tokenId // push back to the original owner, if it fails
	}
}
