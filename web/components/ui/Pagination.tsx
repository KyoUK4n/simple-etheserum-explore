"use client";

import { useState, useEffect } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
    page: number;
    pageSize: number;
    pageCount: number;
    hasMore: boolean;
    onPrev: () => void;
    onNext: () => void;
    onPageSizeChange: (size: number) => void;
    onPageIndexChange: (index: number) => void;
    loading?: boolean;
}

export default function Pagination({
    page,
    pageSize,
    pageCount,
    hasMore,
    onPrev,
    onNext,
    onPageSizeChange,
    onPageIndexChange,
    loading,
}: PaginationProps) {
    const [inputSize, setInputSize] = useState(String(pageSize));
    const [inputIndex, setInputIndex] = useState(String(page));

    useEffect(() => {
        setInputIndex(String(page));
        setInputSize(String(pageSize));
    }, [page, pageSize]);

    const handleSizeBlur = () => {
        const n = parseInt(inputSize);
        if (!isNaN(n) && n > 0 && n !== pageSize) {
            onPageSizeChange(n);
        } else {
            setInputSize(String(pageSize)); // 输入非法时还原
        }
    };

    const handleIndexBlur = () => {
        const n = parseInt(inputIndex);
        if (!isNaN(n) && n > 0 && n !== page) {
            onPageIndexChange(n);
        } else {
            setInputIndex(String(page)); // 输入非法时还原
        }
    };

    return (
        <div className="flex items-center justify-end gap-3 pt-2">
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <span>Show</span>
                <input
                    value={inputSize}
                    onChange={(e) => setInputSize(e.target.value)}
                    onBlur={handleSizeBlur}
                    onKeyDown={(e) => e.key === "Enter" && handleSizeBlur()}
                    className="w-12 text-center bg-input border border-border rounded-md px-1.5 py-1 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary/50 transition-colors"
                />
                <span>records</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <span>Page</span>
                <input
                    value={inputIndex}
                    onChange={(e) => setInputIndex(e.target.value)}
                    onBlur={handleIndexBlur}
                    onKeyDown={(e) => e.key === "Enter" && handleIndexBlur()}
                    className="w-12 text-center bg-input border border-border rounded-md px-1.5 py-1 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary/50 transition-colors"
                />
                <span>of {pageCount}</span>
            </div>
            <button
                onClick={onPrev}
                disabled={page <= 1 || loading}
                className="p-1.5 rounded-md border border-border hover:bg-muted/50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
                <ChevronLeft size={14} />
            </button>
            <button
                onClick={onNext}
                disabled={!hasMore || loading}
                className="p-1.5 rounded-md border border-border hover:bg-muted/50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
                <ChevronRight size={14} />
            </button>
        </div>
    );
}