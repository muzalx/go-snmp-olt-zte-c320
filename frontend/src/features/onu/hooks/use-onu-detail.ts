import { useQuery } from "@tanstack/react-query";
import { getOnuDetail } from "@/features/onu/api/onu.api";

export function useOnuDetail(params: { boardId: number; ponId: number; onuId: number }) {
  return useQuery({
    queryKey: ["onu-detail", params.boardId, params.ponId, params.onuId],
    queryFn: () => getOnuDetail(params),
    enabled: Boolean(params.boardId && params.ponId && params.onuId),
  });
}
