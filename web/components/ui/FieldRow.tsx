interface FieldRowProps {
  label: string;
  children: React.ReactNode;
}

export default function FieldRow({ label, children }: FieldRowProps) {
  return (
    <div className="flex gap-4 py-3 border-b border-border last:border-0">
      <div className="w-40 shrink-0 text-muted-foreground text-xs font-medium uppercase tracking-wide pt-0.5">
        {label}
      </div>
      <div className="flex-1 text-sm break-all">{children}</div>
    </div>
  );
}
