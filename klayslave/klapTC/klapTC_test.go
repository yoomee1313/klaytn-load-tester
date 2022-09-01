package klapTC

import (
	"fmt"
	"testing"

	"github.com/klaytn/klaytn/common"
)

func TestHowmuch(t *testing.T) {
	fmt.Println(len(Klap_account_list))
	var data []byte
	data = common.FromHex("0x976fafc500000000000000000000000078b6adde60a9181c1889913d31906bbf5c3795dd")
	fmt.Println(data)
}
