package ethx

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SubInfo struct {
	ethClient *EthClient
	logsCh    chan types.Log
	filter    ethereum.FilterQuery
	address   string
	Sub       *ethereum.Subscription
	ABI       *abi.ABI
}

type SubscriptionConfig struct {
	AbiPath string
	Address string
	Topics  []string `json:"topics,optional"`
}

// NewSubscribe 需要用ws监听
func NewSubscribe(ethClient *EthClient, sc *SubscriptionConfig, logsCh chan types.Log) (*SubInfo, error) {

	parsedABI, err := ParseABI(sc.AbiPath)
	if err != nil {
		logx.Errorf("Parse abi JSON failed: %s", err)
		return nil, err
	}
	nftAddress := common.HexToAddress(sc.Address)
	filter := ethereum.FilterQuery{
		Addresses: []common.Address{nftAddress},
	}
	topicHashes := make([]common.Hash, 0)
	for _, topic := range sc.Topics {
		if strings.TrimSpace(topic) == "" {
			continue
		}
		topicHashes = append(topicHashes, common.HexToHash(topic))
	}
	filter.Topics = [][]common.Hash{topicHashes}

	subInfo := &SubInfo{
		ethClient: ethClient,
		filter:    filter,
		address:   sc.Address,
		logsCh:    logsCh,
		ABI:       parsedABI,
	}
	if err = subInfo.subscribe(); err != nil {
		return nil, err
	}
	subInfo.errCather()
	return subInfo, nil
}

func (s *SubInfo) Unsubscribe() {
	(*s.Sub).Unsubscribe()
	s.ethClient.Close()
}

func (s *SubInfo) subscribe() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	sub, err := s.ethClient.SubscribeFilterLogs(ctx, s.filter, s.logsCh)
	if err != nil {
		return err
	}
	s.Sub = &sub
	return nil
}

func (s *SubInfo) errCather() {
	go func() {
		for {
			select {
			case err := <-(*s.Sub).Err():
				logx.Errorf("catch an error on subscription [%s] : %s", s.address, err)
				s.resubscribe()
			}
		}
	}()
}

func (s *SubInfo) resubscribe() {

	exponentialBackOff := backoff.NewExponentialBackOff()
	var maxTries uint = 5
	tries := 0
	_, err := backoff.Retry(context.Background(), func() (bool, error) {

		if err := s.ethClient.connect(); err != nil {
			tries++
			logx.Errorf("[%d/%d, next:%v] attempting to reconnect eth client failed: %s", tries, maxTries, exponentialBackOff.NextBackOff(), err)
			return false, err
		}
		if err := s.subscribe(); err != nil {
			logx.Errorf("[%d/%d, next:%v] attempting to resubscribe on [%s] failed: %s", tries, maxTries, exponentialBackOff.NextBackOff(), s.address, err)
			return false, err
		}

		return true, nil
	}, backoff.WithBackOff(exponentialBackOff), backoff.WithMaxTries(maxTries))
	if err != nil {
		logx.Errorf("resubscribe on [%s] failed: %s", s.address, err)
	}
}
