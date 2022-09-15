package klapTC

import (
	"fmt"
	"testing"

<<<<<<< HEAD
	"github.com/klaytn/klaytn/common"
)

func TestHowmuch(t *testing.T) {
	fmt.Println(len(Klap_account_list))
	var data []byte
	data = common.FromHex("0x976fafc500000000000000000000000078b6adde60a9181c1889913d31906bbf5c3795dd")
	fmt.Println(data)
=======
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
>>>>>>> e3713d5 (add klapAppCallAdditionalTC)
}
