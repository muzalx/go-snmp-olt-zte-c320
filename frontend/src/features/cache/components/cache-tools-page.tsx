import { useState } from "react";
import { ApiErrorAlert } from "@/components/common/api-error-alert";
import { LoadingState } from "@/components/common/loading-state";
import { useCacheActions } from "@/features/cache/hooks/use-cache-actions";
import { normalizeApiError } from "@/utils/normalize-api-error";
import { isValidBoardId, isValidPonId } from "@/utils/validators";

export function CacheToolsPage() {
  const [boardId, setBoardId] = useState(1);
  const [ponId, setPonId] = useState(1);

  const invalidFilter = !isValidBoardId(boardId) || !isValidPonId(ponId);
  const { emptyIds, idSn, clear, refresh } = useCacheActions({ boardId, ponId });

  return (
    <section>
      <h2>Cache Tools</h2>
      <div className="card inline">
        <div>
          <label htmlFor="cache-board">Board ID</label>
          <input id="cache-board" type="number" min={1} max={2} value={boardId} onChange={(e) => setBoardId(Number(e.target.value))} />
        </div>
        <div>
          <label htmlFor="cache-pon">PON ID</label>
          <input id="cache-pon" type="number" min={1} max={16} value={ponId} onChange={(e) => setPonId(Number(e.target.value))} />
        </div>
        <button
          onClick={() => clear.mutate()}
          disabled={invalidFilter || clear.isPending}
        >
          {clear.isPending ? "Clearing..." : "Clear Cache"}
        </button>
        <button
          className="secondary"
          onClick={() => refresh.mutate()}
          disabled={invalidFilter || refresh.isPending}
        >
          {refresh.isPending ? "Refreshing..." : "Refresh Empty ONU ID"}
        </button>
      </div>

      {invalidFilter ? <ApiErrorAlert error={{ title: "VALIDATION_ERROR", message: "board_id harus 1-2 dan pon_id harus 1-16." }} /> : null}
      {clear.error ? <ApiErrorAlert error={normalizeApiError(clear.error)} /> : null}
      {refresh.error ? <ApiErrorAlert error={normalizeApiError(refresh.error)} /> : null}
      {emptyIds.error ? <ApiErrorAlert error={normalizeApiError(emptyIds.error)} /> : null}
      {idSn.error ? <ApiErrorAlert error={normalizeApiError(idSn.error)} /> : null}

      {emptyIds.isLoading || idSn.isLoading ? <LoadingState label="Mengambil data ONU ID utility..." /> : null}

      <div className="grid">
        <div className="card">
          <h3>Empty ONU IDs</h3>
          <pre className="code">{JSON.stringify(emptyIds.data?.data || [], null, 2)}</pre>
        </div>
        <div className="card">
          <h3>ONU ID + Serial Number</h3>
          <pre className="code">{JSON.stringify(idSn.data?.data || [], null, 2)}</pre>
        </div>
      </div>
    </section>
  );
}
