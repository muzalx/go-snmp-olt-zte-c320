import { http } from "@/lib/http";
import type { ApiMeta, ApiSuccess } from "@/types/api";
import type { OnuItem } from "@/types/onu";

export const getOnusPaginated = async (params: {
  boardId: number;
  ponId: number;
  page?: number;
  limit?: number;
}) => {
  const { boardId, ponId, page = 1, limit = 10 } = params;
  const res = await http.get<ApiSuccess<OnuItem[], ApiMeta>>(
    `/api/v1/paginate/board/${boardId}/pon/${ponId}`,
    { params: { page, limit } },
  );
  return res.data;
};

export const getOnuDetail = async (params: { boardId: number; ponId: number; onuId: number }) => {
  const { boardId, ponId, onuId } = params;
  const res = await http.get<ApiSuccess<Record<string, unknown>>>(
    `/api/v1/board/${boardId}/pon/${ponId}/onu/${onuId}`,
  );
  return res.data;
};

export const getEmptyOnuIds = async (params: { boardId: number; ponId: number }) => {
  const { boardId, ponId } = params;
  const res = await http.get<ApiSuccess<Array<{ onu_id: number }>>>(
    `/api/v1/board/${boardId}/pon/${ponId}/onu_id/empty`,
  );
  return res.data;
};

export const getOnuIdSn = async (params: { boardId: number; ponId: number }) => {
  const { boardId, ponId } = params;
  const res = await http.get<ApiSuccess<Array<{ onu_id: number; serial_number: string }>>>(
    `/api/v1/board/${boardId}/pon/${ponId}/onu_id_sn`,
  );
  return res.data;
};
