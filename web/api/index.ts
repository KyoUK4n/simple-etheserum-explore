// ============================================================
// API SERVICE LAYER
// Replace mock implementations with real fetch calls to your backend.
// Base URL is read from NEXT_PUBLIC_API_BASE_URL env var (default: localhost:8080).
// ============================================================

import type {
  BlockData,
  TxData,
  TxLog,
  AddressData,
  TokenBalance,
  BaseResponse,
} from "@/types";

// 前端独立开发时通过环境变量修改后端服务地址，通过golang服务挂载打包后的静态资源时，通过url直接访问接口即可
export const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api/v1";

// ============================================================
// API Functions
// ============================================================

export async function apiGetBlock(query: string): Promise<BlockData | null> {

  let apiUrl: string;
  const blockNumber = parseInt(query);
  if (query.startsWith("0x")) {
    apiUrl = `${API_BASE}/blocks/info?hash=${query}`;
  } else if (blockNumber > 0) {
    apiUrl = `${API_BASE}/blocks/info?number=${blockNumber}`;
  } else {
    if (["latest", "pending", "earliest", "safe", "finalized"].includes(query)) {
      apiUrl = `${API_BASE}/blocks/info?tag=${query}`;
    } else {
      console.log("Invalid block query:", query);
      throw new Error("Invalid block query");
    }
  }

  const res = await fetch(apiUrl);
  if (!res.ok) return null;
  const json = await res.json();
  if (json.code !== 1) return null;
  const d = json.data;
  return {
    number: d.number,
    hash: d.hash,
    parentHash: d.parentHash,
    timestamp: d.time,
    miner: d.miner ?? "",
    txCount: d.transactionCount,
    gasUsed: d.gasUsed,
    gasLimit: d.gasLimit,
    size: d.size ?? 0,
    baseFeePerGas: d.baseFee,
  };
}

export async function apiGetTransactionInfo(hash: string): Promise<TxData | null> {
  if (hash && !hash.startsWith("0x")) hash = `0x${hash}`;
  const res = await fetch(`${API_BASE}/transactions/${encodeURIComponent(hash)}`);
  if (!res.ok) return null;
  const json = await res.json();
  if (json.code !== 1) return null;
  const d = json.data;

  let status: TxData["status"];
  if (d.isPending) status = "pending";
  else status = d.status === 1 ? "success" : "failed";

  return {
    hash: d.hash,
    blockNumber: d.blockNumber,
    from: d.from,
    to: d.to ?? null,
    value: (d.value / 1e18).toFixed(18),
    gas: d.gas,
    gasPrice: (d.gasPrice / 1e9).toFixed(2),
    nonce: d.nonce,
    inputData: d.dataLen > 0 ? `${d.dataLen} bytes` : "0x",
    status,
    gasUsed: d.gasUsed,
    logsCount: d.logs,
    timestamp: d.timestamp,
  };
}

export async function apiGetTransactions(address: string, pageIndex: number, pageSize: number): Promise<TxData[] | null> {
  if (address && !address.startsWith("0x")) address = `0x${address}`;
  const res = await fetch(`${API_BASE}/transactions?pageIndex=${pageIndex}&pageSize=${pageSize}&address=${encodeURIComponent(address)}`);
  if (!res.ok) return null;
  const json = await res.json();
  if (json.code !== 1) return null;
  const txs = json.data;

  let txData: TxData[] = [];
  txs.map((d: any) => {
    let status: TxData["status"];
    if (d.isPending) status = "pending";
    else status = d.status === 1 ? "success" : "failed";

    txData.push({
      hash: d.hash,
      blockNumber: d.blockNumber,
      from: d.from,
      to: d.to ?? null,
      value: (d.value / 1e18).toFixed(18),
      gas: d.gas,
      gasPrice: (d.gasPrice / 1e9).toFixed(2),
      nonce: d.nonce,
      inputData: d.dataLen > 0 ? `${d.dataLen} bytes` : "0x",
      status,
      gasUsed: d.gasUsed,
      logsCount: d.logs,
      timestamp: d.timestamp,
    });
  });

  return txData;
}

export async function apiGetEvents(page = 1): Promise<TxLog[]> {
  const res = await fetch(`${API_BASE}/api/events?page=${page}`);
  return res.json();
}

export async function apiSubmitScanTask(fromBlock: number, toBlock: number): Promise<BaseResponse> {
  const res = await fetch(`${API_BASE}/transactions/pull?startBlock=${fromBlock}&endBlock=${toBlock}`);
  return res.json();
}

export async function apiGetBlance(address: string, tokenAddress: string): Promise<TokenBalance | null> {
  const res = await fetch(`${API_BASE}/balances?address=${encodeURIComponent(address)}&tokenAddress=${encodeURIComponent(tokenAddress)}`);
  if (!res.ok) return null;
  const json = await res.json();
  if (json.code !== 1) return null;
  return json.data;
}

export async function apiGetTokensBalances(address: string): Promise<AddressData | null> {
  const res = await fetch(`${API_BASE}/balances/tokens?address=${encodeURIComponent(address)}`);
  if (!res.ok) return null;
  const json = await res.json();
  if (json.code !== 1) return null;
  const d = json.data;
  const addressData: AddressData = {
    address: address,
    ethAmount: d.ethBalance.amount,
    ethDecimals: d.ethBalance.decimals,
    txCount: 0,
    tokens: d.tokenBalances,
    recentTxs: [],
  }
  return addressData;
}

export async function apiGetEventLogs(txHash: string, address: string, pageIndex: number, pageSize: number): Promise<[TxLog[], number]> {
  if (txHash && !txHash.startsWith("0x")) txHash = `0x${txHash}`;
  const res = await fetch(`${API_BASE}/transactions/events?txHash=${txHash}&address=${address}&pageIndex=${pageIndex}&pageSize=${pageSize}`);
  const json = await res.json();
  if (json.code !== 1) {
    throw new Error(json.msg);
  };
  return [
    json.data.list.map((d: any) => ({
      txHash: d.tx_hash,
      address: d.address,
      eventName: d.event_name,
      blockNumber: d.block_number,
      logIndex: d.log_index,
      topics: d.topics,
      data: d.data,
      txTimestamp: d.tx_timestamp,
    })),
    json.data.total
  ];
}
