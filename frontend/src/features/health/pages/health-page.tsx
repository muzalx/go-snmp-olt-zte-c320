import { ApiErrorAlert } from "@/components/common/api-error-alert";
import { LoadingState } from "@/components/common/loading-state";
import { useHealth } from "@/features/health/hooks/use-health";
import { normalizeApiError } from "@/utils/normalize-api-error";

export function HealthPage() {
  const { healthz, readyz, version } = useHealth();

  const anyError = healthz.error || readyz.error || version.error;
  if (anyError) {
    return <ApiErrorAlert error={normalizeApiError(anyError)} />;
  }

  if (healthz.isLoading || readyz.isLoading || version.isLoading) {
    return <LoadingState label="Mengambil status health..." />;
  }

  return (
    <section>
      <h2>Health Dashboard</h2>
      <div className="grid">
        <div className="card">
          <h3>Liveness</h3>
          <p>{healthz.data?.status || "unknown"}</p>
        </div>
        <div className="card">
          <h3>Readiness</h3>
          <p>{readyz.data?.status || "unknown"}</p>
        </div>
      </div>
      <div className="card">
        <h3>Dependencies</h3>
        <pre className="code">{JSON.stringify(readyz.data?.dependencies || {}, null, 2)}</pre>
      </div>
      <div className="card">
        <h3>Version</h3>
        <pre className="code">{JSON.stringify(version.data || {}, null, 2)}</pre>
      </div>
    </section>
  );
}
