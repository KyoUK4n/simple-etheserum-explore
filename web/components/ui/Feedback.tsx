import { Loader2 } from "lucide-react";

export function EmptyState({ message }: { message: string }) {
  return (
    <div className="text-center py-16 text-muted-foreground text-sm">
      {message}
    </div>
  );
}

export function Spinner() {
  return <Loader2 size={16} className="animate-spin text-muted-foreground" />;
}
