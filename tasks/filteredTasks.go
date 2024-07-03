package tasks

import (
	"math/big"

	"github.com/klaytn/klaytn-load-tester/account"

	"github.com/myzhan/boomer"
)

type FilteredTasks []*ExtendedTask

func NewFilteredTasks(tcStrLists []string, accGrpUnsigned []*account.Account, gEndpoint string) {
	var filteredTasks FilteredTasks
	for _, tc := range tcStrLists {
		if task := validTasks[tc]; task != nil {
			filteredTasks = append(filteredTasks, task)
			println("=> " + tc + " task is added.")
		}
	}

	// Import/Unlock Account on the node if there is a task to use unsigned account group.
	for _, task := range filteredTasks {
		if task.Task.Name == "transferUnsignedTx" {
			for _, acc := range accGrpUnsigned {
				acc.ImportUnLockAccount(gEndpoint)
			}
			break // to import/unlock once.
		}
	}
}

type ExtendedTask struct {
	Init func([]*account.Account, string, *big.Int)
	Task boomer.Task
}

func (et *ExtendedTask) GetAccGrp() {
	if et.Task.Name == "transferUnsignedTx" {
		return accGrpUnsigned
	}
	return accGrpSigned
}
