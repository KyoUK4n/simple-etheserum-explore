"use client";

import { useState } from "react";
import { Copy, CheckCircle } from "lucide-react";

interface CopyButtonProps {
  text: string;
}

export default function CopyButton({ text }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    });
  };

  return (
    <button
      onClick={handleCopy}
      className="text-muted-foreground hover:text-foreground transition-colors ml-1"
      title="Copy"
    >
      {copied
        ? <CheckCircle size={13} className="text-emerald-400" />
        : <Copy size={13} />}
    </button>
  );
}
