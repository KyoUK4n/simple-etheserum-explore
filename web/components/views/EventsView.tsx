"use client";

import { useState, useCallback } from "react";
import { RefreshCw, ChevronDown, Zap } from "lucide-react";
import { apiGetEventLogs, apiSubmitScanTask } from "@/api";
import { formatNumber, shortHash } from "@/utils";
import type { TxLog } from "@/types";
import Card from "@/components/ui/Card";
import Mono from "@/components/ui/Mono";
import Pagination from "@/components/ui/Pagination";
import Badge from "@/components/ui/Badge";
import CopyButton from "@/components/ui/CopyButton";
import { EmptyState, Spinner } from "@/components/ui/Feedback";

function eventColor(name?: string): "success" | "pending" | "info" {
  if (!name) return "info";
  if (name === "Transfer") return "success";
  if (name === "Swap") return "pending";
  return "info";
}

export default function EventsView() {
  const [events, setEvents] = useState<TxLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [expandedIndex, setExpandedIndex] = useState<number | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [pageCount, setPageCount] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [fromBlock, setFromBlock] = useState("");
  const [toBlock, setToBlock] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [scanResult, setScanResult] = useState<string | null>(null);

  const loadEvents = useCallback(async (pageIndex: number, pageSize: number) => {
    setLoading(true);
    setError(null);
    setExpandedIndex(null);

    try {
      const [eventLogs, total] = await apiGetEventLogs("", "", pageIndex, pageSize);
      const calcPageCount = Math.ceil(total / pageSize);
      setPage(pageIndex);
      setPageCount(calcPageCount);
      setEvents(eventLogs);
      setHasMore(pageIndex < calcPageCount);
      setLoaded(true);
    } catch (error: any) {
      console.log(error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    loadEvents(1, size); // 改 pageSize 时回到第1页
  };

  const handlePageIndexChange = (index: number) => {
    loadEvents(index, pageSize);
  };

  const handleScan = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!fromBlock || !toBlock) return;
    setSubmitting(true);
    try {
      const res = await apiSubmitScanTask(parseInt(fromBlock), parseInt(toBlock));
      if (res.code === 1) {
        setScanResult("Task submitted successfully");
      } else {
        setSubmitError(res.msg);
      }
    } catch (error: any) {
      setSubmitError(error.message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-5">
      {/* Scan task form */}
      <Card title="Submit Block Scan Task">
        <form onSubmit={handleScan} className="flex gap-3 flex-wrap">
          <div className="flex-1 min-w-36">
            <label className="block text-xs text-muted-foreground mb-1.5 uppercase tracking-wide">From Block</label>
            <input
              type="number"
              value={fromBlock}
              onChange={(e) => setFromBlock(e.target.value)}
              placeholder="e.g. 22847000"
              className="w-full bg-input border border-border rounded-md px-3 py-2 text-sm font-mono text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary/50 transition-colors"
            />
          </div>
          <div className="flex-1 min-w-36">
            <label className="block text-xs text-muted-foreground mb-1.5 uppercase tracking-wide">To Block</label>
            <input
              type="number"
              value={toBlock}
              onChange={(e) => setToBlock(e.target.value)}
              placeholder="e.g. 22847391"
              className="w-full bg-input border border-border rounded-md px-3 py-2 text-sm font-mono text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary/50 transition-colors"
            />
          </div>
          <div className="self-end">
            <button
              type="submit"
              disabled={submitting || !fromBlock || !toBlock}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm font-semibold hover:bg-primary/90 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
            >
              {submitting ? <Spinner /> : <Zap size={14} />}
              Submit Scan
            </button>
          </div>
        </form>
        {submitError ? (
          <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
            <Badge variant="danger">{submitError}</Badge>
          </div>
        ) : (
          <div className="flex items-center gap-2 p-4 bg-green-500/10 border border-green-500/20 rounded-lg">
            <Badge variant="success">{scanResult}</Badge>
          </div>
        )}
      </Card>

      {/* Events list */}
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-foreground">Recent Events</h3>
        <button
          onClick={() => { setPage(1); loadEvents(1, pageSize); }}
          disabled={loading}
          className="flex items-center gap-1.5 text-xs text-cyan-400 hover:text-cyan-300 transition-colors disabled:opacity-50"
        >
          <RefreshCw size={12} className={loading ? "animate-spin" : ""} />
          {loaded ? "Refresh" : "Load Events"}
        </button>
      </div>

      {!loaded && !loading && (
        <Card>
          <EmptyState message='Click "Load Events" to fetch recent on-chain events.' />
        </Card>
      )}

      {error && (
        <div className="flex items-center gap-2 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <Badge variant="danger">{error}</Badge>
        </div>
      )}

      {loading && (
        <div className="flex items-center gap-2 text-muted-foreground text-sm py-8 justify-center">
          <Spinner /> Loading events...
        </div>
      )}

      {loaded && events.length > 0 && (
        <>
          <div className="bg-card border border-border rounded-lg overflow-hidden">
            {events.map((ev, i) => (
              <div key={i} className="border-b border-border/60 last:border-0">
                {/* 列表行，点击展开 */}
                <div
                  onClick={() => setExpandedIndex(expandedIndex === i ? null : i)}
                  className="w-full flex items-start gap-4 px-5 py-4 hover:bg-muted/30 transition-colors text-left"
                >
                  <div className="shrink-0 w-6 h-6 rounded border border-border bg-muted flex items-center justify-center text-xs font-mono text-muted-foreground mt-0.5">
                    {ev.logIndex}
                  </div>
                  <div className="flex-1 min-w-0 space-y-1">
                    <div className="flex items-center gap-2 flex-wrap">
                      {ev.eventName && (
                        <Badge variant={eventColor(ev.eventName)}>{ev.eventName}</Badge>
                      )}
                      <Mono className="text-cyan-400 text-xs">{ev.address}</Mono>
                    </div>
                    <div className="flex items-center gap-3 text-xs text-muted-foreground">
                      <span>
                        Tx: <Mono className="text-cyan-400">{shortHash(ev.txHash, 6)}</Mono>
                        <CopyButton text={ev.txHash} />
                      </span>
                      <span>Block: <Mono className="text-foreground">#{ev.blockNumber}</Mono></span>
                      <span>{ev.txTimestamp}</span>
                    </div>
                  </div>
                  <ChevronDown
                    size={14}
                    className={`shrink-0 text-muted-foreground transition-transform duration-200 mt-1 ${expandedIndex === i ? "rotate-180" : ""}`}
                  />
                </div>

                {/* 展开详情，和 TxView 里的 log 详情样式一致 */}
                {expandedIndex === i && (
                  <div className="px-5 pb-5 space-y-4 bg-muted/20">
                    {ev.topics.length > 0 && (
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
                              {ev.topics.map((t, ti) => (
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

                    {ev.data.length > 0 && (
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
                              {ev.data.map((d, di) => (
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
          <Pagination
            page={page}
            pageSize={pageSize}
            pageCount={pageCount}
            hasMore={hasMore}
            onPrev={() => loadEvents(page - 1, pageSize)}
            onNext={() => loadEvents(page + 1, pageSize)}
            onPageSizeChange={handlePageSizeChange}
            onPageIndexChange={handlePageIndexChange}
            loading={loading}
          />
        </>
      )}
    </div>
  );
}