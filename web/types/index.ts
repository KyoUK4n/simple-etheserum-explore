export interface BlockData {
  number: number;
  hash: string;
  parentHash: string;
  timestamp: number; // unix seconds
  miner: string;
  txCount: number;
  gasUsed: number;
  gasLimit: number;
  size: number;
  baseFeePerGas?: number;
}

export interface TxData {
  hash: string;
  blockNumber: number;
  from: string;
  to: string;
  value: string; // in ETH
  gas: number;
  gasPrice: string; // in Gwei
  nonce: number;
  inputData: string;
  status: "success" | "failed" | "pending";
  gasUsed: number;
  logsCount: number;
  timestamp: number;
}

export interface TxLogParam {
  index: number;
  name: string;
  type: string;
  value: string;
}

export interface TxLog {
  txHash: string;
  address: string;
  eventName?: string;
  blockNumber: number;
  logIndex: number;
  topics: TxLogParam[];
  data: TxLogParam[];
  txTimestamp: string;
}

export interface TokenBalance {
  symbol: string;
  name: string;
  address: string;
  amount: bigint;
  decimals: number;
}

export interface AddressData {
  address: string;
  ethAmount: bigint;
  ethDecimals: number;
  txCount: number;
  tokens: TokenBalance[];
  recentTxs: TxData[] | null;
}

export interface BaseResponse {
  code: number;
  msg: string;
}
