"use client";

import { useState, useCallback } from "react";
import { CheckCircle, XCircle, Clock, ChevronDown } from "lucide-react";
import { apiGetTransactionInfo, apiGetEventLogs } from "@/api";
import { formatNumber, formatTimestamp, timeAgo } from "@/utils";
import type { TxData, TxLog } from "@/types";
import SearchInput from "@/components/ui/SearchInput";
import Card from "@/components/ui/Card";
import FieldRow from "@/components/ui/FieldRow";
import Mono from "@/components/ui/Mono";
import Badge from "@/components/ui/Badge";
import CopyButton from "@/components/ui/CopyButton";
import { EmptyState, Spinner } from "@/components/ui/Feedback";

function statusVariant(s: TxData["status"]) {
  return s === "success" ? "success" : s === "failed" ? "danger" : "pending";
}

export default function TxView() {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<TxData | null>(null);
  const [notFound, setNotFound] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showFullInput, setShowFullInput] = useState(false);
  const [logsExpanded, setLogsExpanded] = useState(false);
  const [logs, setLogs] = useState<TxLog[]>([]);
  const [logsLoading, setLogsLoading] = useState(false);
  const [expandedLogIndex, setExpandedLogIndex] = useState<number | null>(null);

  const handleSearch = useCallback(async () => {
    if (!query.trim()) return;
    setLoading(true);
    setNotFound(false);
    setData(null);
    setLogs([]);
    setLogsExpanded(false);
    try {
      const result = await apiGetTransactionInfo(query.trim());
      if (result) setData(result);
      else setNotFound(true);
      setLoading(false);
    } catch (error: any) {
      setError(error.message);
    }
  }, [query]);

  const handleToggleLogs = useCallback(async (tx: TxData) => {
    if (logsExpanded) { setLogsExpanded(false); return; }
    setLogsExpanded(true);
    if (logs.length === 0 && tx.logsCount > 0) {
      setLogsLoading(true);
      try {
        const [eventLogs, total] = await apiGetEventLogs(tx.hash, "", 1, 100);
        setLogs(eventLogs);
        setLogsLoading(false);
      } catch (error: any) {
        console.log(error);
      }

    }
  }, [logsExpanded, logs]);

  return (
    <div className="space-y-5">
      <SearchInput
        value={query}
        onChange={setQuery}
        onSubmit={handleSearch}
        placeholder="Enter transaction hash (0x...)"
        loading={loading}
      />

      {loading && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
          <Spinner /> Fetching transaction...
        </div>
      )}

      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <Badge variant="danger">{error}</Badge>
        </div>
      )}

      {notFound && !loading && (
        <Card>
          <EmptyState message="Transaction not found. Ensure you entered a valid tx hash." />
        </Card>
      )}

      {data && !loading && (
        <Card title="Transaction Details">
          <FieldRow label="Tx Hash">
            <Mono className="text-cyan-400">{data.hash}</Mono>
            <CopyButton text={data.hash} />
          </FieldRow>
          <FieldRow label="Status">
            <div className="flex items-center gap-2">
              {data.status === "success"
                ? <CheckCircle size={14} className="text-emerald-400" />
                : data.status === "failed"
                  ? <XCircle size={14} className="text-red-400" />
                  : <Clock size={14} className="text-amber-400" />}
              <Badge variant={statusVariant(data.status)}>
                {data.status.charAt(0).toUpperCase() + data.status.slice(1)}
              </Badge>
            </div>
          </FieldRow>
          <FieldRow label="Block">
            <Mono className="text-cyan-400">#{data.blockNumber}</Mono>
          </FieldRow>
          <FieldRow label="Timestamp">
            <span className="text-foreground">{formatTimestamp(data.timestamp)}</span>
            <span className="text-muted-foreground ml-2 text-xs">({timeAgo(data.timestamp)})</span>
          </FieldRow>
          <FieldRow label="From">
            <Mono className="text-cyan-400">{data.from}</Mono>
            <CopyButton text={data.from} />
          </FieldRow>
          <FieldRow label="To">
            {data.to ? (
              <>
                <Mono className="text-cyan-400">{data.to}</Mono>
                <CopyButton text={data.to} />
              </>
            ) : (
              <Badge variant="info">Contract Creation</Badge>
            )}
          </FieldRow>
          <FieldRow label="Value">
            <Mono className="text-foreground font-semibold">{data.value} ETH</Mono>
          </FieldRow>
          <FieldRow label="Gas Limit / Used">
            <Mono className="text-foreground">{formatNumber(data.gas)}</Mono>
            <span className="text-muted-foreground mx-2">/</span>
            <Mono className="text-foreground">{formatNumber(data.gasUsed)}</Mono>
            <Badge variant="info">{((data.gasUsed / data.gas) * 100).toFixed(1)}%</Badge>
          </FieldRow>
          <FieldRow label="Gas Price">
            <Mono className="text-foreground">{data.gasPrice} Gwei</Mono>
          </FieldRow>
          <FieldRow label="Nonce">
            <Mono className="text-foreground">{data.nonce}</Mono>
          </FieldRow>
          <FieldRow label="Event Logs">
            <button onClick={() => handleToggleLogs(data)} className="flex items-center gap-2 group">
              <Badge variant={data.logsCount > 0 ? "info" : "pending"}>
                {data.logsCount} log{data.logsCount !== 1 ? "s" : ""}
              </Badge>
              {data.logsCount > 0 && (
                <span className="flex items-center gap-1 text-xs text-cyan-400 group-hover:text-cyan-300 transition-colors">
                  {logsExpanded ? "Hide" : "View logs"}
                  <ChevronDown size={12} className={`transition-transform duration-200 ${logsExpanded ? "rotate-180" : ""}`} />
                </span>
              )}
            </button>
          </FieldRow>
          <FieldRow label="Input Data">
            <div className="space-y-2">
              <Mono className={`text-muted-foreground break-all ${showFullInput ? "" : "line-clamp-2"}`}>
                {data.inputData}
              </Mono>
              {data.inputData.length > 40 && (
                <button
                  onClick={() => setShowFullInput((v) => !v)}
                  className="text-xs text-cyan-400 hover:text-cyan-300 transition-colors"
                >
                  {showFullInput ? "Show less" : "Show more"}
                </button>
              )}
            </div>
          </FieldRow>
        </Card>
      )}

      {/* Expandable Event Logs Panel */}
      {data && logsExpanded && (
        <div className="bg-card border border-border rounded-lg overflow-hidden">
          <div className="flex items-center justify-between px-5 py-3 border-b border-border">
            <h3 className="text-sm font-semibold text-foreground">
              Event Logs
              <span className="ml-2 text-xs font-normal text-muted-foreground">({data.logsCount})</span>
            </h3>
          </div>

          {logsLoading && (
            <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
              <Spinner /> Loading logs...
            </div>
          )}

          {!logsLoading && logs.map((log, i) => (
            <div key={i} className="border-b border-border/60 last:border-0">
              <button
                onClick={() => setExpandedLogIndex(expandedLogIndex === log.logIndex ? null : log.logIndex)}
                className="w-full flex items-start gap-4 px-5 py-4 hover:bg-muted/30 transition-colors text-left"
              >
                <div className="shrink-0 w-6 h-6 rounded border border-border bg-muted flex items-center justify-center text-xs font-mono text-muted-foreground mt-0.5">
                  {log.logIndex}
                </div>
                <div className="flex-1 min-w-0 space-y-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    {log.eventName && <Badge variant="info">{log.eventName}</Badge>}
                    <Mono className="text-cyan-400 text-xs">{log.address}</Mono>
                  </div>
                  <div className="text-xs text-muted-foreground font-mono">{log.txTimestamp}</div>
                </div>
                <ChevronDown
                  size={14}
                  className={`shrink-0 text-muted-foreground transition-transform duration-200 mt-1 ${expandedLogIndex === log.logIndex ? "rotate-180" : ""}`}
                />
              </button>

              {expandedLogIndex === log.logIndex && (
                <div className="px-5 pb-5 space-y-4 bg-muted/20">
                  {/* Topics */}
                  {log.topics.length > 0 && (
                    <div>
                      <div className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground mb-2">Topics</div>
                      <div className="bg-muted rounded-md overflow-hidden">
                        <table className="w-full text-xs">
                          <thead>
                            <tr className="border-b border-border/60">
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Name</th>
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Type</th>
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Value</th>
                            </tr>
                          </thead>
                          <tbody>
                            {log.topics.map((t, ti) => (
                              <tr key={ti} className="border-b border-border/40 last:border-0">
                                <td className="px-3 py-2 font-medium text-foreground">{t.name}</td>
                                <td className="px-3 py-2"><Badge variant="pending">{t.type}</Badge></td>
                                <td className="px-3 py-2">
                                  <Mono className="text-cyan-400 break-all">{t.value}</Mono>
                                  <CopyButton text={t.value} />
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  )}

                  {/* Data */}
                  {log.data.length > 0 && (
                    <div>
                      <div className="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground mb-2">Data</div>
                      <div className="bg-muted rounded-md overflow-hidden">
                        <table className="w-full text-xs">
                          <thead>
                            <tr className="border-b border-border/60">
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Name</th>
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Type</th>
                              <th className="text-left px-3 py-2 text-muted-foreground font-medium">Value</th>
                            </tr>
                          </thead>
                          <tbody>
                            {log.data.map((d, di) => (
                              <tr key={di} className="border-b border-border/40 last:border-0">
                                <td className="px-3 py-2 font-medium text-foreground">{d.name}</td>
                                <td className="px-3 py-2"><Badge variant="pending">{d.type}</Badge></td>
                                <td className="px-3 py-2">
                                  <Mono className="text-cyan-400 break-all">{d.value}</Mono>
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
