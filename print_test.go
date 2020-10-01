package printing_test

import (
	"fmt"
	"testing"

	p "github.com/githomework/apps-util-printing"
)

func TestNetworkNames(t *testing.T) {

	for k, v := range p.NetworkPrinters() {
		fmt.Println(k, v)
	}
}
