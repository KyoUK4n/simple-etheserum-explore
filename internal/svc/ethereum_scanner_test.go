package svc

import (
	"context"
	"testing"

	"github.com/KyoUK4n/etherscan/pkg/ethx"
	"github.com/stretchr/testify/assert"
)

func TestScanner_scanTransactionsBatch(t *testing.T) {
	ethClient, err := ethx.NewEthClient("https://sepolia.infura.io/v3/79271c24cfeb47489cdca58648c6ae48")
	if err != nil {
		t.Fatal(err)
	}
	scanner := NewScanner(ethClient, nil, nil, nil, 10)
	batch, err := scanner.scanTransactionsBatch(context.Background(), 11085765, 11085775)
	if err != nil {
		t.Fatal(err)
	}
	assert.Greater(t, len(batch), 0)
}
