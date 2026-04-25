import { Link } from "react-router-dom";
import type { OnuItem } from "@/types/onu";

type Props = {
  rows: OnuItem[];
  boardId: number;
  ponId: number;
};

export function OnuTable({ rows, boardId, ponId }: Props) {
  if (!rows.length) {
    return <p className="muted">Tidak ada ONU ditemukan untuk filter ini.</p>;
  }

  return (
    <div className="card">
      <h3>ONU List</h3>
      <table>
        <thead>
          <tr>
            <th>ONU ID</th>
            <th>Serial Number</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((item, index) => {
            const onuId = Number(item.onu_id ?? index + 1);
            return (
              <tr key={`${onuId}-${index}`}>
                <td>{onuId}</td>
                <td>{String(item.serial_number ?? "-")}</td>
                <td>{String(item.status ?? "-")}</td>
                <td>
                  <Link to={`/onu/${boardId}/${ponId}/${onuId}`}>Detail</Link>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
