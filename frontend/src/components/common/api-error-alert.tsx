import type { NormalizedApiError } from "@/utils/normalize-api-error";

export function ApiErrorAlert({ error }: { error: NormalizedApiError }) {
  return (
    <div className="error-box" role="alert">
      <p>
        <strong>{error.title}</strong>
      </p>
      <p>{error.message}</p>
      {error.requestId ? <p>request_id: {error.requestId}</p> : null}
    </div>
  );
}
