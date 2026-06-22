export function delay(ms: number) {
  return new Promise((r) => setTimeout(r, ms));
}

export function shortHash(hash: string, chars = 8) {
  if (!hash) return "";
  return `${hash.slice(0, chars + 2)}...${hash.slice(-chars)}`;
}

export function formatTimestamp(ts: number) {
  const d = new Date(ts * 1000);
  return d.toLocaleString("en-US", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

export function timeAgo(ts: number) {
  const diff = Math.floor(Date.now() / 1000 - ts);
  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

export function formatNumber(n: number) {
  return n.toLocaleString("en-US");
}

export function gasPercent(used: number, limit: number) {
  return ((used / limit) * 100).toFixed(1);
}
