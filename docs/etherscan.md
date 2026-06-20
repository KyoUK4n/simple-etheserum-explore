### 1. "获取某个地址的ETH、ERC20代币余额，若TokenAddress未传，则获取ETH余额"

1. route definition

- Url: /api/v1/balances
- Method: GET
- Request: `GetBalanceReq`
- Response: `Response`

2. request definition



```golang
type GetBalanceReq struct {
	Address string `form:"address"`
	TokenAddress string `form:"tokenAddress,optional"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 2. "获取区间[start - end]区块信息"

1. route definition

- Url: /api/v1/blocks
- Method: GET
- Request: `GetBlockInfoReq`
- Response: `Response`

2. request definition



```golang
type GetBlockInfoReq struct {
	Number int64 `form:"number,optional"`
	Hash string `form:"hash,optional"`
	Tag string `form:"tag,optional"`
	Start int64 `form:"start,optional"`
	End int64 `form:"end,optional"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 3. "获取某个区块信息，按照number-&gt;hash-&gt;tag优先级获取"

1. route definition

- Url: /api/v1/blocks/info
- Method: GET
- Request: `GetBlockInfoReq`
- Response: `Response`

2. request definition



```golang
type GetBlockInfoReq struct {
	Number int64 `form:"number,optional"`
	Hash string `form:"hash,optional"`
	Tag string `form:"tag,optional"`
	Start int64 `form:"start,optional"`
	End int64 `form:"end,optional"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 4. "从DB分页查询交易列表"

1. route definition

- Url: /api/v1/transactions
- Method: GET
- Request: `GetTransactionsReq`
- Response: `Response`

2. request definition



```golang
type GetTransactionsReq struct {
	Address string `form:"address"`
	PageIndex int `form:"pageIndex"`
	PageSize int `form:"pageSize"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 5. "根据hash获取交易信息"

1. route definition

- Url: /api/v1/transactions/:hash
- Method: GET
- Request: `GetTransactionInfoReq`
- Response: `Response`

2. request definition



```golang
type GetTransactionInfoReq struct {
	Hash string `path:"hash"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 6. "获取交易的日志"

1. route definition

- Url: /api/v1/transactions/events
- Method: GET
- Request: `QueryEventLogReq`
- Response: `Response`

2. request definition



```golang
type QueryEventLogReq struct {
	Address string `form:"address,optional"`
	TxHash string `form:"txHash,optional"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

### 7. "提交区块同步任务，任务串行执行，上一个任务未完成时，接口直接返回失败"

1. route definition

- Url: /api/v1/transactions/pull
- Method: GET
- Request: `PullTransactionsReq`
- Response: `Response`

2. request definition



```golang
type PullTransactionsReq struct {
	StartBlock uint64 `form:"startBlock"`
	EndBlock uint64 `form:"endBlock"`
}
```


3. response definition



```golang
type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

