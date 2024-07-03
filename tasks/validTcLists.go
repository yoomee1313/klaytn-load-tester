package tasks

import (
	"github.com/klaytn/klaytn-load-tester/testcases/analyticTC"

	"github.com/myzhan/boomer"
)

var validTasks = map[string]*ExtendedTask{
	"analyticTx": {
		Task: boomer.Task{Weight: 10, Fn: analyticTC.Run, Name: "analyticTx"},
		Init: analyticTC.Init,
	},
}

/*
// initTCList initializes TCs and returns a slice of TCs.
func initTCList() (taskSet []*ExtendedTask) {
	taskSet = append(taskSet, &ExtendedTask{
		Name:    "analyticTx",
		Fn:      analyticTC.Run,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "analyticQueryLargestAccBalTx",

		Fn:      analyticTC.QueryLargestAccBal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "analyticQueryLargestTxValTx",

		Fn:      analyticTC.QueryLargestTxVal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "analyticQueryTotalTxValTx",

		Fn:      analyticTC.QueryTotalTxVal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "cpuHeavyTx",

		Fn:      cpuHeavyTC.Run,
		Init:    cpuHeavyTC.Init,
		AccGrp:  accGrpForSignedTx, //[nUserForSigned/2:],
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "doNothingTx",

		Fn:      doNothingTC.Run,
		Init:    doNothingTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: internalTxTC.Name,

		Fn:      internalTxTC.Run,
		Init:    internalTxTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: internalTxTC.NameMintNFT,

		Fn:      internalTxTC.RunMintNFT,
		Init:    internalTxTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ioHeavyTx",

		Fn:      ioHeavyTC.Run,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ioHeavyScanTx",

		Fn:      ioHeavyTC.Scan,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ioHeavyWriteTx",

		Fn:      ioHeavyTC.Write,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "largeMemoTC",

		Fn:      largeMemoTC.Run,
		Init:    largeMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: readAPITCs.Name,

		Fn:      readAPITCs.Run,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankTx",

		Fn:      smallBankTC.Run,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankAlmagateTx",

		Fn:      smallBankTC.Almagate,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankGetBalanceTx",

		Fn:      smallBankTC.GetBalance,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankSendPaymentTx",

		Fn:      smallBankTC.SendPayment,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankUpdateBalanceTx",

		Fn:      smallBankTC.UpdateBalance,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankUpdateSavingTx",

		Fn:      smallBankTC.UpdateSaving,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "smallBankWriteCheckTx",

		Fn:      smallBankTC.WriteCheck,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "transferSignedTx",

		Fn:      transferSignedTc.Run,
		Init:    transferSignedTc.Init,
		AccGrp:  accGrpForSignedTx, //[:nUserForSigned/2-1],
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newValueTransferTC",

		Fn:      newValueTransferTC.Run,
		Init:    newValueTransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newValueTransferWithCancelTC",

		Fn:      newValueTransferWithCancelTC.Run,
		Init:    newValueTransferWithCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedValueTransferTC",

		Fn:      newFeeDelegatedValueTransferTC.Run,
		Init:    newFeeDelegatedValueTransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedValueTransferWithRatioTC",

		Fn:      newFeeDelegatedValueTransferWithRatioTC.Run,
		Init:    newFeeDelegatedValueTransferWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newValueTransferMemoTC",

		Fn:      newValueTransferMemoTC.Run,
		Init:    newValueTransferMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newValueTransferLargeMemoTC",

		Fn:      newValueTransferLargeMemoTC.Run,
		Init:    newValueTransferLargeMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newValueTransferSmallMemoTC",

		Fn:      newValueTransferSmallMemoTC.Run,
		Init:    newValueTransferSmallMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedValueTransferMemoTC",

		Fn:      newFeeDelegatedValueTransferMemoTC.Run,
		Init:    newFeeDelegatedValueTransferMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedValueTransferMemoWithRatioTC",

		Fn:      newFeeDelegatedValueTransferMemoWithRatioTC.Run,
		Init:    newFeeDelegatedValueTransferMemoWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newAccountCreationTC",

		Fn:      newAccountCreationTC.Run,
		Init:    newAccountCreationTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newAccountUpdateTC",

		Fn:      newAccountUpdateTC.Run,
		Init:    newAccountUpdateTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedAccountUpdateTC",

		Fn:      newFeeDelegatedAccountUpdateTC.Run,
		Init:    newFeeDelegatedAccountUpdateTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedAccountUpdateWithRatioTC",

		Fn:      newFeeDelegatedAccountUpdateWithRatioTC.Run,
		Init:    newFeeDelegatedAccountUpdateWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newSmartContractDeployTC",

		Fn:      newSmartContractDeployTC.Run,
		Init:    newSmartContractDeployTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedSmartContractDeployTC",

		Fn:      newFeeDelegatedSmartContractDeployTC.Run,
		Init:    newFeeDelegatedSmartContractDeployTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedSmartContractDeployWithRatioTC",

		Fn:      newFeeDelegatedSmartContractDeployWithRatioTC.Run,
		Init:    newFeeDelegatedSmartContractDeployWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newSmartContractExecutionTC",

		Fn:      newSmartContractExecutionTC.Run,
		Init:    newSmartContractExecutionTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: storageTrieWriteTC.Name,

		Fn:      storageTrieWriteTC.Run,
		Init:    storageTrieWriteTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedSmartContractExecutionTC",

		Fn:      newFeeDelegatedSmartContractExecutionTC.Run,
		Init:    newFeeDelegatedSmartContractExecutionTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedSmartContractExecutionWithRatioTC",

		Fn:      newFeeDelegatedSmartContractExecutionWithRatioTC.Run,
		Init:    newFeeDelegatedSmartContractExecutionWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newCancelTC",

		Fn:      newCancelTC.Run,
		Init:    newCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedCancelTC",

		Fn:      newFeeDelegatedCancelTC.Run,
		Init:    newFeeDelegatedCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newFeeDelegatedCancelWithRatioTC",

		Fn:      newFeeDelegatedCancelWithRatioTC.Run,
		Init:    newFeeDelegatedCancelWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "transferSignedWithCheckTx",

		Fn:      transferSignedWithCheckTc.Run,
		Init:    transferSignedWithCheckTc.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "transferUnsignedTx",

		Fn:      transferUnsignedTc.Run,
		Init:    transferUnsignedTc.Init,
		AccGrp:  accGrpForUnsignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "userStorageSetTx",

		Fn:      userStorageTC.RunSet,
		Init:    userStorageTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "userStorageSetGetTx",

		Fn:      userStorageTC.RunSetGet,
		Init:    userStorageTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ycsbTx",

		Fn:      ycsbTC.Run,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ycsbGetTx",

		Fn:      ycsbTC.Get,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ycsbSetTx",

		Fn:      ycsbTC.Set,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: erc20TransferTC.Name,

		Fn:      erc20TransferTC.Run,
		Init:    erc20TransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: erc721TransferTC.Name,

		Fn:      erc721TransferTC.Run,
		Init:    erc721TransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readGasPrice",

		Fn:      readAPITCs.GasPrice,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readBlockNumber",

		Fn:      readAPITCs.BlockNumber,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readGetBlockByNumber",

		Fn:      readAPITCs.GetBlockByNumber,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readGetAccount",

		Fn:      readAPITCs.GetAccount,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readGetBlockWithConsensusInfoByNumber",

		Fn:      readAPITCs.GetBlockWithConsensusInfoByNumber,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readGetStorageAt",

		Fn:      readAPITCs.GetStorageAt,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readCall",

		Fn:      readAPITCs.Call,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "readEstimateGas",

		Fn:      readAPITCs.EstimateGas,
		Init:    readAPITCs.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ethereumTxLegacyTC",

		Fn:      ethereumTxLegacyTC.Run,
		Init:    ethereumTxLegacyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ethereumTxAccessListTC",

		Fn:      ethereumTxAccessListTC.Run,
		Init:    ethereumTxAccessListTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "ethereumTxDynamicFeeTC",

		Fn:      ethereumTxDynamicFeeTC.Run,
		Init:    ethereumTxDynamicFeeTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newEthereumAccessListTC",

		Fn:      newEthereumAccessListTC.Run,
		Init:    newEthereumAccessListTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name: "newEthereumDynamicFeeTC",

		Fn:      newEthereumDynamicFeeTC.Run,
		Init:    newEthereumDynamicFeeTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	return taskSet
}
*/
