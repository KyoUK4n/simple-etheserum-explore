import { Blocks, ChevronRight } from "lucide-react";

export default function Header() {
  return (
    <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-10">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center gap-4">
        <div className="flex items-center gap-2.5 shrink-0">
          <div className="w-7 h-7 bg-primary rounded-md flex items-center justify-center">
            <Blocks size={14} className="text-primary-foreground" />
          </div>
          <span className="font-semibold text-foreground text-sm tracking-tight">
            ChainExplorer
          </span>
          <span className="hidden sm:inline-flex items-center gap-1 text-[10px] font-medium px-1.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
            Sepolia
          </span>
        </div>
        <div className="flex-1" />
        <div className="text-xs text-muted-foreground hidden sm:block font-mono">
          
        </div>
      </div>
    </header>
  );
}
