interface BadgeProps {
  variant: "success" | "danger" | "pending" | "info";
  children: React.ReactNode;
}

const styles = {
  success: "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20",
  danger: "bg-red-500/10 text-red-400 border border-red-500/20",
  pending: "bg-amber-500/10 text-amber-400 border border-amber-500/20",
  info: "bg-cyan-500/10 text-cyan-400 border border-cyan-500/20",
};

export default function Badge({ variant, children }: BadgeProps) {
  return (
    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium ${styles[variant]}`}>
      {children}
    </span>
  );
}
