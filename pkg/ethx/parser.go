package ethx

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/KyoUK4n/etherscan/pkg/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeromicro/go-zero/core/logx"
)

type EventLog struct {
	TxHash      string
	Address     string
	EventName   string
	BlockNumber int64
	LogIndex    int64
	Topics      string
	Data        string
	TxTimestamp time.Time
}

type EventLogData struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func ParseABI(abiPath string) (*abi.ABI, error) {
	abiFile, err := os.OpenFile(abiPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer abiFile.Close()

	parsedABI, err := abi.JSON(abiFile)
	if err != nil {
		logx.Errorf("Parse abi JSON failed: %s", err)
		return nil, err
	}
	return &parsedABI, nil
}

func ParseLogEvent(vLog *types.Log, parsedABI abi.ABI) *EventLog {
	// 检查是否有 Topics（没有 Topics 的日志可能是无效的）
	if len(vLog.Topics) == 0 {
		return nil
	}

	// 步骤 1: 识别事件类型
	// Topics[0] 是事件签名的 keccak256 哈希值
	// 例如: Transfer(address,address,uint256) 的哈希
	eventTopic := vLog.Topics[0]

	// 尝试识别是哪个事件（通过比较 Topics[0] 和事件签名的哈希）
	var eventName string
	var eventSig abi.Event

	// 遍历 ABI 中定义的所有事件，查找匹配的事件签名
	for name, event := range parsedABI.Events {
		// 计算事件的签名哈希
		eventSigHash := crypto.Keccak256Hash([]byte(event.Sig))
		if eventSigHash == eventTopic {
			eventName = name
			eventSig = event
			break
		}
	}

	if eventName == "" {
		// 如果无法识别事件类型，打印原始信息
		logx.Infof("[%s] Unknown Event - Block: %d, Tx: %s, Topic[0]: %s\n",
			time.Now().Format(time.RFC3339),
			vLog.BlockNumber,
			vLog.TxHash.Hex(),
			eventTopic.Hex(),
		)
		return nil
	}

	// 步骤 2: 解析事件参数
	eventLog := &EventLog{
		TxHash:      vLog.TxHash.Hex(),
		Address:     vLog.Address.Hex(),
		EventName:   eventName,
		BlockNumber: int64(vLog.BlockNumber),
		LogIndex:    int64(vLog.Index),
		TxTimestamp: time.Unix(int64(vLog.BlockTimestamp), 0),
	}

	// 步骤 3: 解析 indexed 参数（从 Topics 中解析）
	// Topics[0] 是事件签名哈希，Topics[1..N] 是 indexed 参数
	// 注意：只有前 3 个 indexed 参数会放在 Topics 中（Ethereum 限制）
	// Topics[0] 是事件签名，所以 indexed 参数从 Topics[1] 开始
	// 注意：topicIndex 只针对 indexed 参数计数，不考虑非 indexed 参数

	indexedData := make([]*EventLogData, 0)

	indexedParamIndex := 0
	for i, input := range eventSig.Inputs {
		if !input.Indexed {
			continue
		}
		// indexed 参数在 Topics 中的位置 = 1 + indexed 参数的索引
		topicIndex := 1 + indexedParamIndex
		indexedParamIndex++

		if topicIndex >= len(vLog.Topics) {
			continue
		}

		topic := vLog.Topics[topicIndex]
		eventLogData := &EventLogData{
			Index: i + 1,
			Name:  input.Name,
			Type:  input.Type.String(),
		}

		// 根据类型解析 indexed 参数
		switch input.Type.T {
		case abi.AddressTy:
			// address 类型：去除前 12 字节的 0 填充，后 20 字节是地址
			addr := common.BytesToAddress(topic.Bytes())
			eventLogData.Value = addr.Hex()
		case abi.IntTy, abi.UintTy:
			// 整数类型：直接转换为 big.Int
			value := new(big.Int).SetBytes(topic.Bytes())
			eventLogData.Value = value.String()
		case abi.BoolTy:
			// bool 类型：检查最后一个字节
			eventLogData.Value = fmt.Sprintf("%t", topic[31] == 1)
		case abi.BytesTy, abi.FixedBytesTy:
			// bytes 类型：直接显示十六进制
			eventLogData.Value = topic.Hex()
		default:
			// 其他类型：显示原始十六进制
			eventLogData.Value = topic.Hex()
		}

		indexedData = append(indexedData, eventLogData)
	}

	nonIndexedData := make([]*EventLogData, 0)

	// 步骤 4: 解析非 indexed 参数（从 Data 字段中解析）
	// Data 字段包含所有非 indexed 参数的编码数据
	if len(vLog.Data) > 0 {
		// 创建一个结构体来接收解码后的参数
		// 注意：这里使用通用方法，实际应用中可能需要根据具体事件定义结构体
		nonIndexedInputs := make([]abi.Argument, 0)
		for _, input := range eventSig.Inputs {
			if !input.Indexed {
				nonIndexedInputs = append(nonIndexedInputs, input)
			}
		}

		if len(nonIndexedInputs) > 0 {
			// 使用 ABI 解码 Data 字段
			// 方法 1: 使用 UnpackIntoInterface（需要预定义结构体）
			// 方法 2: 使用 Unpack（返回 []interface{}）
			values, err := parsedABI.Unpack(eventName, vLog.Data)
			if err != nil {
				logx.Errorf("Unpack ABI failed: %s", err)
			} else {
				// 只输出非 indexed 参数
				nonIndexedIdx := 0
				for i, input := range eventSig.Inputs {
					if !input.Indexed {
						if nonIndexedIdx < len(values) {
							value := values[nonIndexedIdx]
							eventLogData := &EventLogData{
								Index: i + 1,
								Name:  input.Name,
								Type:  input.Type.String(),
							}
							// 根据类型格式化输出
							switch v := value.(type) {
							case *big.Int:
								eventLogData.Value = v.String()
							case common.Address:
								eventLogData.Value = v.Hex()
							case []byte:
								eventLogData.Value = fmt.Sprintf("0x%x", v)
							default:
								eventLogData.Value = fmt.Sprintf("%v", v)
							}
							nonIndexedData = append(nonIndexedData, eventLogData)
							nonIndexedIdx++
						}
					}
				}
			}
		}
	}

	indexedDataBytes, err := json.Marshal(&indexedData)
	if err != nil {
		logx.Errorf("marshal indexed data failed: %s", err)
	} else {
		eventLog.Topics = utils.BytesToString(indexedDataBytes)
	}

	nonIndexedDataBytes, err := json.Marshal(&nonIndexedData)
	if err != nil {
		logx.Errorf("marshal nonIndexed data failed: %s", err)
	} else {
		eventLog.Data = utils.BytesToString(nonIndexedDataBytes)
	}
	return eventLog
}
