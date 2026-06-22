// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"time"

	"github.com/KyoUK4n/etherscan/pkg/ethx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql DBConfig
	Eth   EthConfig
}

type DBConfig struct {
	Datasource string
}

type EthConfig struct {
	RpcUrl           string
	WSUrl            string
	MaxBlockRange    int64
	Subscriptions    []*ethx.SubscriptionConfig
	RequestPeriod    time.Duration
	ReplayStartBlock uint64
	ReplayEndBlock   uint64
	ReplayCount      uint64
	ReplayStep       uint64
	Tokens           []string
}
