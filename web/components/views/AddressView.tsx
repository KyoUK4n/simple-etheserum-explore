"use client";

import { useState, useCallback } from "react";
import { apiGetTokensBalances, apiGetTransactions } from "@/api";
import { formatNumber, shortHash, timeAgo } from "@/utils";
import type { AddressData, TxData } from "@/types";
import SearchInput from "@/components/ui/SearchInput";
import Card from "@/components/ui/Card";
import Mono from "@/components/ui/Mono";
import Badge from "@/components/ui/Badge";
import CopyButton from "@/components/ui/CopyButton";
import { EmptyState, Spinner } from "@/components/ui/Feedback";
import { formatUnits } from 'viem';

function statusVariant(s: TxData["status"]): "success" | "danger" | "pending" {
  return s === "success" ? "success" : s === "failed" ? "danger" : "pending";
}

export default function AddressView() {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<AddressData | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeSubTab, setActiveSubTab] = useState<"txs" | "tokens">("txs");

  const handleSearch = useCallback(async () => {
    if (!query.trim()) return;
    setLoading(true);
    setNotFound(false);
    setData(null);

    try {
      // 获取代币余额
      const addressData = await apiGetTokensBalances(query.trim());

      // 获取交易记录
      const txs = await apiGetTransactions(query.trim(), 1, 20);
      if (addressData && txs) {
        addressData.recentTxs = txs;
      }

      setData(addressData);
    } catch (error: any) {
      setNotFound(true);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  }, [query]);

  return (
    <div className="space-y-5">
      <SearchInput
        value={query}
        onChange={setQuery}
        onSubmit={handleSearch}
        placeholder="Enter wallet or contract address (0x...)"
        loading={loading}
      />

      {loading && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
          <Spinner /> Fetching address...
        </div>
      )}

      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <Badge variant="danger">{error}</Badge>
        </div>
      )}

      {notFound && !loading && (
        <Card>
          <EmptyState message="Address not found or has no activity." />
        </Card>
      )}

      {data && !loading && (
        <>
          {/* Overview stats */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div className="bg-card border border-border rounded-lg p-4">
              <div className="text-xs text-muted-foreground uppercase tracking-wide mb-1.5">ETH Balance</div>
              <div className="text-xl font-semibold text-foreground">
                {formatUnits(data.ethAmount, data.ethDecimals)}{" "}
                <span className="text-muted-foreground text-sm font-normal">ETH</span>
              </div>
            </div>
            <div className="bg-card border border-border rounded-lg p-4">
              <div className="text-xs text-muted-foreground uppercase tracking-wide mb-1.5">Transactions</div>
              <div className="text-xl font-semibold text-foreground">{data.txCount}</div>
            </div>
            <div className="bg-card border border-border rounded-lg p-4">
              <div className="text-xs text-muted-foreground uppercase tracking-wide mb-1.5">ERC-20 Tokens</div>
              <div className="text-xl font-semibold text-foreground">
                {data.tokens?.length}{" "}
                <span className="text-muted-foreground text-sm font-normal">types</span>
              </div>
            </div>
          </div>

          {/* Address line */}
          <div className="bg-card border border-border rounded-lg p-4 flex items-center gap-2">
            <span className="text-xs text-muted-foreground uppercase tracking-wide shrink-0">Address</span>
            <Mono className="text-cyan-400 text-sm ml-2">{data.address}</Mono>
            <CopyButton text={data.address} />
          </div>

          {/* Sub-tabs */}
          <div className="flex gap-0 border-b border-border">
            {(["txs", "tokens"] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveSubTab(tab)}
                className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${activeSubTab === tab
                  ? "border-primary text-foreground"
                  : "border-transparent text-muted-foreground hover:text-foreground"
                  }`}
              >
                {tab === "txs" ? "Transactions" : "Token Balances"}
              </button>
            ))}
          </div>

          {/* Transactions tab */}
          {activeSubTab === "txs" && (
            <div className="bg-card border border-border rounded-lg overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full text-xs">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      {["Tx Hash", "Block", "From","", "To", "Value", "Age"].map((h) => (
                        <th key={h} className="text-left px-4 py-3 text-muted-foreground font-medium uppercase tracking-wide">
                          {h}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {data.recentTxs?.map((tx, i) => {
                      const isOut = tx.from.toLowerCase() === data.address.toLowerCase();
                      return (
                        <tr key={i} className="border-b border-border/60 hover:bg-muted/30 transition-colors">
                          <td className="px-4 py-3">
                            <Mono className="text-cyan-400">{shortHash(tx.hash, 6)}</Mono>
                            <CopyButton text={tx.hash} />
                          </td>
                          <td className="px-4 py-3">
                            <Mono className="text-foreground">#{tx.blockNumber}</Mono>
                          </td>
                          <td className="px-4 py-3 space-y-0.5">
                            <div className="flex items-center gap-1.5">
                              <Mono className="text-muted-foreground">
                                {shortHash(tx.from, 6)}
                              </Mono>
                              <CopyButton text={tx.from} />
                            </div>
                          </td>
                          <td className="px-4 py-3 space-y-0.5">
                            <div className="flex items-center gap-1.5">
                            <span className={`text-[10px] font-semibold px-1.5 py-0.5 rounded ${isOut ? "bg-red-500/10 text-red-400" : "bg-emerald-500/10 text-emerald-400"}`}>
                              {isOut ? "OUT" : "IN"}
                            </span>
                          </div>
                          </td>
                          <td className="px-4 py-3 space-y-0.5">
                            <div className="flex items-center gap-1.5">
                              <Mono className="text-muted-foreground">
                                {shortHash(tx.to, 6)}
                              </Mono>
                              <CopyButton text={tx.to} />
                            </div>
                          </td>
                          <td className="px-4 py-3">
                            <Mono className="text-foreground">{tx.value} ETH</Mono>
                          </td>
                          <td className="px-4 py-3 text-muted-foreground whitespace-nowrap">
                            {tx.timestamp ? timeAgo(tx.timestamp): "?"}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {/* Tokens tab */}
          {activeSubTab === "tokens" && (
            <div className="bg-card border border-border rounded-lg overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full text-xs">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      {["Token", "Contract", "Balance"].map((h, i) => (
                        <th key={h} className={`px-4 py-3 text-muted-foreground font-medium uppercase tracking-wide ${i === 2 ? "text-right" : "text-left"}`}>
                          {h}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {data.tokens?.map((tok, i) => (
                      <tr key={i} className="border-b border-border/60 hover:bg-muted/30 transition-colors">
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <div className="w-6 h-6 rounded-full bg-muted flex items-center justify-center text-[10px] font-bold text-muted-foreground">
                              {tok.symbol.slice(0, 2)}
                            </div>
                            <div>
                              <div className="font-semibold text-foreground">{tok.symbol}</div>
                              <div className="text-muted-foreground">{tok.name}</div>
                            </div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <Mono className="text-muted-foreground">{shortHash(tok.address, 6)}</Mono>
                          <CopyButton text={tok.address} />
                        </td>
                        <td className="px-4 py-3 text-right">
                          <Mono className="text-foreground font-semibold">{formatUnits(tok.amount, tok.decimals)}</Mono>
                          <span className="text-muted-foreground ml-1">{tok.symbol}</span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
