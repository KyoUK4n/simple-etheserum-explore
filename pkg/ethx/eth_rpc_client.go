package ethx

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type EthClient struct {
	logx.Logger
	*ethclient.Client
	url string
}

func NewEthClient(rpcUrl string) (*EthClient, error) {

	client := &EthClient{
		Logger: logx.WithContext(context.Background()).WithFields(logx.LogField{
			Key:   "module",
			Value: "EthRpcClient",
		}),
		url: rpcUrl,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *EthClient) connect() error {

	var maxTries uint = 5
	tries := 0
	exponentialBackOff := backoff.NewExponentialBackOff()
	client, err := backoff.Retry(context.Background(), func() (*ethclient.Client, error) {
		tries++
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		client, err := ethclient.DialContext(ctx, c.url)
		if err != nil {
			logx.Errorf(
				"[%d/%d, next %v] attempting to dial on [%s] failed: %s",
				tries,
				maxTries,
				exponentialBackOff.NextBackOff(),
				c.url,
				err,
			)
			return nil, err
		}
		return client, nil
	}, backoff.WithMaxTries(maxTries), backoff.WithBackOff(exponentialBackOff))
	if err != nil {
		logx.Errorf("dial on [%s] failed: %s", c.url, err)
		return err
	}
	c.Client = client
	return nil
}
