import { useQuery } from "@tanstack/react-query";
import { getOnusPaginated } from "@/features/onu/api/onu.api";

export function useOnuList(params: { boardId: number; ponId: number; page: number; limit: number }) {
  return useQuery({
    queryKey: ["onu-list", params.boardId, params.ponId, params.page, params.limit],
    queryFn: () => getOnusPaginated(params),
    enabled: Boolean(params.boardId && params.ponId),
  });
}
