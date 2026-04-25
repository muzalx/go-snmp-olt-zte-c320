import { Link, useParams } from "react-router-dom";
import { ApiErrorAlert } from "@/components/common/api-error-alert";
import { LoadingState } from "@/components/common/loading-state";
import { useOnuDetail } from "@/features/onu/hooks/use-onu-detail";
import { normalizeApiError } from "@/utils/normalize-api-error";
import { isValidBoardId, isValidOnuId, isValidPonId } from "@/utils/validators";

export function OnuDetailPage() {
  const params = useParams();

  const boardId = Number(params.boardId);
  const ponId = Number(params.ponId);
  const onuId = Number(params.onuId);

  if (!isValidBoardId(boardId) || !isValidPonId(ponId) || !isValidOnuId(onuId)) {
    return <ApiErrorAlert error={{ title: "VALIDATION_ERROR", message: "Parameter board/pon/onu tidak valid." }} />;
  }

  const query = useOnuDetail({ boardId, ponId, onuId });

  return (
    <section>
      <p>
        <Link to="/onu">← Kembali ke ONU Explorer</Link>
      </p>
      <h2>ONU Detail</h2>

      {query.isLoading ? <LoadingState label="Mengambil detail ONU..." /> : null}
      {query.error ? <ApiErrorAlert error={normalizeApiError(query.error)} /> : null}

      {query.data ? (
        <div className="card">
          <pre className="code">{JSON.stringify(query.data.data, null, 2)}</pre>
        </div>
      ) : null}
    </section>
  );
}
