import { useState } from "react";

type Props = {
  onApply: (value: { boardId: number; ponId: number; page: number; limit: number }) => void;
};

export function OnuFilterForm({ onApply }: Props) {
  const [boardId, setBoardId] = useState(1);
  const [ponId, setPonId] = useState(1);
  const [limit, setLimit] = useState(10);

  return (
    <div className="card">
      <h3>Filter Board / PON</h3>
      <form
        className="inline"
        onSubmit={(event) => {
          event.preventDefault();
          onApply({ boardId, ponId, page: 1, limit });
        }}
      >
        <div>
          <label htmlFor="board">Board ID</label>
          <input id="board" type="number" min={1} max={2} value={boardId} onChange={(e) => setBoardId(Number(e.target.value))} />
        </div>
        <div>
          <label htmlFor="pon">PON ID</label>
          <input id="pon" type="number" min={1} max={16} value={ponId} onChange={(e) => setPonId(Number(e.target.value))} />
        </div>
        <div>
          <label htmlFor="limit">Limit</label>
          <select id="limit" value={limit} onChange={(e) => setLimit(Number(e.target.value))}>
            {[10, 25, 50].map((value) => (
              <option key={value} value={value}>
                {value}
              </option>
            ))}
          </select>
        </div>
        <button type="submit">Load ONU</button>
      </form>
    </div>
  );
}
