import { http } from "@/lib/http";
import type { ApiSuccess } from "@/types/api";
import { getEmptyOnuIds, getOnuIdSn } from "@/features/onu/api/onu.api";

export { getEmptyOnuIds, getOnuIdSn };

export const clearCache = async (params: { boardId: number; ponId: number }) => {
  const { boardId, ponId } = params;
  const res = await http.delete<ApiSuccess<{ message: string }>>(
    `/api/v1/board/${boardId}/pon/${ponId}/cache/clear`,
  );
  return res.data;
};

export const refreshEmptyOnuIdCache = async (params: { boardId: number; ponId: number }) => {
  const { boardId, ponId } = params;
  const res = await http.post<ApiSuccess<string>>(`/api/v1/board/${boardId}/pon/${ponId}/onu_id/update`);
  return res.data;
};
