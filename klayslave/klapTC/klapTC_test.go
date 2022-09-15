package klapTC

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHowmuch(t *testing.T) {
	fmt.Println("klap_account_list length:", len(klap_account_list))
	fmt.Println("dataList length:", len(dataList))
	fmt.Println("contractAddrList length:", len(contractAddrList))
	assert.Equal(t, len(dataList), len(contractAddrList))
	fmt.Println("dataListAdditional length:", len(dataListAdditional))
	fmt.Println("contractAddrListAdditional length:", len(contractAddrListAdditional))
	assert.Equal(t, len(dataListAdditional),len(contractAddrListAdditional))
}
