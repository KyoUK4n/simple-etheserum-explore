import { ExternalLink } from "lucide-react";
import { API_BASE } from "@/api";

export default function Footer() {
  return (
    <footer className="border-t border-border mt-12 py-4 px-4">
      <div className="max-w-5xl mx-auto flex items-center justify-between text-xs text-muted-foreground">
        <span>ChainExplorer &mdash; Lightweight Blockchain Browser</span>
        {/* <span className="font-mono flex items-center gap-1.5">
          <ExternalLink size={11} />
          {API_BASE}
        </span> */}
      </div>
    </footer>
  );
}
