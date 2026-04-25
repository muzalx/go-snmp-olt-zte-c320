export type ApiSuccess<TData, TMeta = unknown> = {
  code: number;
  status: "success";
  data: TData;
  meta?: TMeta;
};

export type ApiErrorCode =
  | "VALIDATION_ERROR"
  | "NOT_FOUND"
  | "SNMP_ERROR"
  | "REDIS_ERROR"
  | "CONFIG_ERROR"
  | "INTERNAL_ERROR";

export type ApiErrorResponse = {
  code: number;
  status: string;
  error_code?: ApiErrorCode;
  data: string | { message?: string; details?: Record<string, unknown> };
  request_id?: string;
};

export type ApiMeta = {
  page?: number;
  limit?: number;
  total_rows?: number;
  total_page?: number;
};
