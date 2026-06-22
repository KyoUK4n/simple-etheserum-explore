"use client";

import { Search } from "lucide-react";
import { Spinner } from "./Feedback";

interface SearchInputProps {
  value: string;
  onChange: (v: string) => void;
  onSubmit: () => void;
  placeholder: string;
  loading?: boolean;
}

export default function SearchInput({
  value,
  onChange,
  onSubmit,
  placeholder,
  loading,
}: SearchInputProps) {
  return (
    <form
      onSubmit={(e) => { e.preventDefault(); onSubmit(); }}
      className="flex gap-2"
    >
      <div className="relative flex-1">
        <Search
          size={15}
          className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none"
        />
        <input
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="w-full bg-input border border-border rounded-md pl-9 pr-4 py-2.5 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary/50 font-mono transition-colors"
        />
      </div>
      <button
        type="submit"
        disabled={loading || !value.trim()}
        className="px-4 py-2.5 bg-primary text-primary-foreground rounded-md text-sm font-semibold hover:bg-primary/90 disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
      >
        {loading ? <Spinner /> : <Search size={14} />}
        Search
      </button>
    </form>
  );
}
