interface CardProps {
  title?: string;
  children: React.ReactNode;
  action?: React.ReactNode;
}

export default function Card({ title, children, action }: CardProps) {
  return (
    <div className="bg-card border border-border rounded-lg overflow-hidden">
      {title && (
        <div className="flex items-center justify-between px-5 py-3 border-b border-border">
          <h3 className="text-sm font-semibold text-foreground">{title}</h3>
          {action}
        </div>
      )}
      <div className="px-5 py-4">{children}</div>
    </div>
  );
}
