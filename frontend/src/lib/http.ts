import axios from "axios";
import { env } from "@/lib/env";

export const http = axios.create({
  baseURL: env.apiBaseUrl,
  timeout: 20000,
  headers: {
    "Content-Type": "application/json",
  },
});

http.interceptors.request.use((config) => {
  if (env.apiKey) {
    config.headers["X-API-Key"] = env.apiKey;
  }

  return config;
});
