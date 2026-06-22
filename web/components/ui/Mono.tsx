interface MonoProps {
  children: React.ReactNode;
  className?: string;
}

export default function Mono({ children, className = "" }: MonoProps) {
  return (
    <span className={`font-mono text-xs ${className}`}>{children}</span>
  );
}
