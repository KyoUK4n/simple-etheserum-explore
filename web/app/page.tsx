"use client";

import { useState } from "react";
import { Blocks, ArrowRightLeft, Activity, Wallet } from "lucide-react";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import BlockView from "@/components/views/BlockView";
import TxView from "@/components/views/TxView";
import EventsView from "@/components/views/EventsView";
import AddressView from "@/components/views/AddressView";

const TABS = [
  { id: "blocks", label: "Block Query", icon: Blocks },
  { id: "txs", label: "Transaction", icon: ArrowRightLeft },
  { id: "events", label: "Events", icon: Activity },
  { id: "address", label: "Address", icon: Wallet },
] as const;

type TabId = typeof TABS[number]["id"];

export default function Page() {
  const [activeTab, setActiveTab] = useState<TabId>("blocks");

  return (
    <div className="min-h-screen bg-background text-foreground">
      <Header />

      <main className="max-w-screen-2xl mx-auto px-4 sm:px-6 py-6">
        {/* Tab navigation */}
        <nav className="flex gap-1 mb-6 bg-card border border-border rounded-lg p-1">
          {TABS.map(({ id, label, icon: Icon }) => (
            <button
              key={id}
              onClick={() => setActiveTab(id)}
              className={`flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-all ${activeTab === id
                  ? "bg-muted text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
                }`}
            >
              <Icon size={14} />
              <span className="hidden sm:inline">{label}</span>
            </button>
          ))}
        </nav>

        {/* Tab content */}
        <div className={activeTab === "blocks" ? "" : "hidden"}><BlockView /></div>
        <div className={activeTab === "txs" ? "" : "hidden"}><TxView /></div>
        <div className={activeTab === "events" ? "" : "hidden"}><EventsView /></div>
        <div className={activeTab === "address" ? "" : "hidden"}><AddressView /></div>
      </main>

      <Footer />
    </div>
  );
}
