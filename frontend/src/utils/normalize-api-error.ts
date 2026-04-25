import type { AxiosError } from "axios";
import type { ApiErrorResponse } from "@/types/api";

export type NormalizedApiError = {
  title: string;
  message: string;
  requestId?: string;
};

export function normalizeApiError(error: unknown): NormalizedApiError {
  const fallback: NormalizedApiError = {
    title: "Request failed",
    message: "Terjadi kesalahan yang tidak diketahui.",
  };

  if (!error || typeof error !== "object") {
    return fallback;
  }

  const axiosError = error as AxiosError<ApiErrorResponse>;
  const payload = axiosError.response?.data;
  if (!payload) {
    return fallback;
  }

  const message =
    typeof payload.data === "string"
      ? payload.data
      : payload.data?.message || "Terjadi kesalahan pada API.";

  return {
    title: payload.error_code || payload.status || "API Error",
    message,
    requestId: payload.request_id,
  };
}
