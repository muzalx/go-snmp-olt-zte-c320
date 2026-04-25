import { useState } from "react";
import { ApiErrorAlert } from "@/components/common/api-error-alert";
import { LoadingState } from "@/components/common/loading-state";
import { OnuFilterForm } from "@/features/onu/components/onu-filter-form";
import { OnuTable } from "@/features/onu/components/onu-table";
import { useOnuList } from "@/features/onu/hooks/use-onu-list";
import { normalizeApiError } from "@/utils/normalize-api-error";
import { isValidBoardId, isValidPonId } from "@/utils/validators";

export function OnuExplorerPage() {
  const [params, setParams] = useState({ boardId: 1, ponId: 1, page: 1, limit: 10 });
  const query = useOnuList(params);

  const invalidFilter = !isValidBoardId(params.boardId) || !isValidPonId(params.ponId);

  return (
    <section>
      <h2>ONU Explorer</h2>
      <OnuFilterForm onApply={setParams} />

      {invalidFilter ? <ApiErrorAlert error={{ title: "VALIDATION_ERROR", message: "board_id harus 1-2 dan pon_id harus 1-16." }} /> : null}

      {query.isLoading ? <LoadingState label="Mengambil daftar ONU..." /> : null}
      {query.error ? <ApiErrorAlert error={normalizeApiError(query.error)} /> : null}

      <OnuTable rows={query.data?.data || []} boardId={params.boardId} ponId={params.ponId} />

      <div className="inline">
        <button
          className="secondary"
          disabled={params.page <= 1}
          onClick={() => setParams((prev) => ({ ...prev, page: Math.max(prev.page - 1, 1) }))}
        >
          Prev
        </button>
        <p>Page: {params.page}</p>
        <button onClick={() => setParams((prev) => ({ ...prev, page: prev.page + 1 }))}>Next</button>
      </div>
    </section>
  );
}
