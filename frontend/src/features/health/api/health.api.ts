import { http } from "@/lib/http";

export type ReadyzResponse = {
  status: string;
  dependencies?: Record<string, { state?: string; error?: string }>;
};

export type VersionResponse = {
  version: string;
  api_version: string;
  commit: string;
  build_time: string;
  uptime: string;
};

export const getHealthz = async () => (await http.get<{ status: string }>("/healthz")).data;
export const getReadyz = async () => (await http.get<ReadyzResponse>("/readyz")).data;
export const getVersion = async () => (await http.get<VersionResponse>("/version")).data;
