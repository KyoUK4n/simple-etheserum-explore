"use client";

import { useState, useCallback } from "react";
import { apiGetBlock } from "@/api";
import { formatNumber, formatTimestamp, timeAgo, gasPercent } from "@/utils";
import type { BlockData } from "@/types";
import SearchInput from "@/components/ui/SearchInput";
import Card from "@/components/ui/Card";
import FieldRow from "@/components/ui/FieldRow";
import Mono from "@/components/ui/Mono";
import Badge from "@/components/ui/Badge";
import CopyButton from "@/components/ui/CopyButton";
import { EmptyState, Spinner } from "@/components/ui/Feedback";

export default function BlockView() {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<BlockData | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);


  const handleSearch = useCallback(async () => {
    if (!query.trim()) return;
    setLoading(true);
    setNotFound(false);
    setData(null);
    setError(null);
    try {
      const result = await apiGetBlock(query.trim());
      if (result) setData(result);
      else setNotFound(true);
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
        placeholder="Enter block number or block hash (0x...)"
        loading={loading}
      />
      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <Badge variant="danger">{error}</Badge>
        </div>
      )}

      {loading && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
          <Spinner /> Fetching block...
        </div>
      )}

      {notFound && !loading && (
        <Card>
          <EmptyState message="Block not found. Check the number or hash and try again." />
        </Card>
      )}

      {data && !loading && (
        <Card title={`Block #${data.number}`}>
          <FieldRow label="Block Number">
            <span className="text-foreground font-semibold">#{data.number}</span>
          </FieldRow>
          <FieldRow label="Hash">
            <Mono className="text-cyan-400">{data.hash}</Mono>
            <CopyButton text={data.hash} />
          </FieldRow>
          <FieldRow label="Parent Hash">
            <Mono className="text-muted-foreground">{data.parentHash}</Mono>
            <CopyButton text={data.parentHash} />
          </FieldRow>
          <FieldRow label="Timestamp">
            <span className="text-foreground">{formatTimestamp(data.timestamp)}</span>
            <span className="text-muted-foreground ml-2 text-xs">({timeAgo(data.timestamp)})</span>
          </FieldRow>
          <FieldRow label="Transactions">
            <Badge variant="info">{data.txCount} txns</Badge>
          </FieldRow>
          <FieldRow label="Gas Used / Limit">
            <div className="space-y-1.5">
              <div className="flex items-center gap-3">
                <Mono className="text-foreground">{formatNumber(data.gasUsed)}</Mono>
                <span className="text-muted-foreground">/</span>
                <Mono className="text-muted-foreground">{formatNumber(data.gasLimit)}</Mono>
                <Badge variant={parseFloat(gasPercent(data.gasUsed, data.gasLimit)) > 80 ? "danger" : "success"}>
                  {gasPercent(data.gasUsed, data.gasLimit)}%
                </Badge>
              </div>
              <div className="w-48 h-1.5 bg-muted rounded-full overflow-hidden">
                <div
                  className="h-full bg-emerald-500 rounded-full transition-all"
                  style={{ width: `${gasPercent(data.gasUsed, data.gasLimit)}%` }}
                />
              </div>
            </div>
          </FieldRow>
          {data.baseFeePerGas && (
            <FieldRow label="Base Fee">
              <Mono className="text-foreground">
                {(data.baseFeePerGas / 1e18).toFixed(18)} ETH
                <span className="text-muted-foreground ml-2 text-xs">
                  ({(data.baseFeePerGas / 1e9).toFixed(2)} Gwei)
                </span>
              </Mono>
            </FieldRow>
          )}
          <FieldRow label="Miner">
            <Mono className="text-cyan-400">{data.miner}</Mono>
            <CopyButton text={data.miner} />
          </FieldRow>
          <FieldRow label="Size">
            <Mono className="text-foreground">{formatNumber(data.size)} bytes</Mono>
          </FieldRow>
        </Card>
      )}
    </div>
  );
}
