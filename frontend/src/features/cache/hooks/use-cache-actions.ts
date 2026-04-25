import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  clearCache,
  getEmptyOnuIds,
  getOnuIdSn,
  refreshEmptyOnuIdCache,
} from "@/features/cache/api/cache.api";

export function useCacheActions(params: { boardId: number; ponId: number }) {
  const queryClient = useQueryClient();

  const emptyIds = useQuery({
    queryKey: ["empty-onu-ids", params.boardId, params.ponId],
    queryFn: () => getEmptyOnuIds(params),
  });

  const idSn = useQuery({
    queryKey: ["onu-id-sn", params.boardId, params.ponId],
    queryFn: () => getOnuIdSn(params),
  });

  const clear = useMutation({
    mutationFn: () => clearCache(params),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["onu-list"] });
    },
  });

  const refresh = useMutation({
    mutationFn: () => refreshEmptyOnuIdCache(params),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["empty-onu-ids", params.boardId, params.ponId] });
    },
  });

  return { emptyIds, idSn, clear, refresh };
}
