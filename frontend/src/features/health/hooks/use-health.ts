import { useQuery } from "@tanstack/react-query";
import { getHealthz, getReadyz, getVersion } from "@/features/health/api/health.api";

export function useHealth() {
  const healthz = useQuery({ queryKey: ["healthz"], queryFn: getHealthz, refetchInterval: 20_000 });
  const readyz = useQuery({ queryKey: ["readyz"], queryFn: getReadyz, refetchInterval: 20_000 });
  const version = useQuery({ queryKey: ["version"], queryFn: getVersion });

  return { healthz, readyz, version };
}
