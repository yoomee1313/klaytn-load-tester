package klapTC

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHowmuch(t *testing.T) {
	fmt.Println("klap_account_list length:", len(Klap_account_list))
	fmt.Println("dataList length:", len(DataList))
	fmt.Println("contractAddrList length:", len(ContractAddrList))
	assert.Equal(t, len(DataList), len(ContractAddrList))
	fmt.Println("dataListAdditional length:", len(DataListAdditional))
	fmt.Println("contractAddrListAdditional length:", len(ContractAddrListAdditional))
	assert.Equal(t, len(DataListAdditional),len(ContractAddrListAdditional))
}
