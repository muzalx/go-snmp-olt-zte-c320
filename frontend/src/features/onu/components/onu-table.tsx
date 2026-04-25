import { Link } from "react-router-dom";
import type { OnuItem } from "@/types/onu";

type Props = {
  rows: OnuItem[];
  boardId: number;
  ponId: number;
};

const excludedTooltipKeys = new Set(["serial_number", "rx_power", "tx_power"]);

function formatFieldLabel(value: string): string {
  return value
    .replaceAll("_", " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function toReadable(value: unknown): string {
  if (value === null || value === undefined || value === "") {
    return "-";
  }

  if (typeof value === "object") {
    return JSON.stringify(value);
  }

  return String(value);
}

export function OnuTable({ rows, boardId, ponId }: Props) {
  if (!rows.length) {
    return <p className="muted">Tidak ada ONU ditemukan untuk filter ini.</p>;
  }

  return (
    <div className="card">
      <h3>ONU List</h3>
      <div className="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ONU ID</th>
              <th>Serial Number</th>
              <th>RX Power</th>
              <th>TX Power</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((item, index) => {
              const onuId = Number(item.onu_id ?? index + 1);
              const serialNumber = toReadable(item.serial_number);
              const rxPower = toReadable(item.rx_power);
              const txPower = toReadable(item.tx_power);

              const tooltipDetails = Object.entries(item)
                .filter(([key]) => !excludedTooltipKeys.has(key))
                .map(([key, value]) => ({ key, value: toReadable(value) }));

              return (
                <tr key={`${onuId}-${index}`}>
                  <td>{onuId}</td>
                  <td>
                    <span className="sn-tooltip" tabIndex={0} aria-label={`Detail tambahan ONU ${serialNumber}`}>
                      {serialNumber}
                      <span className="sn-tooltip-content" role="tooltip">
                        {tooltipDetails.length ? (
                          <ul>
                            {tooltipDetails.map((entry) => (
                              <li key={`${onuId}-${entry.key}`}>
                                <strong>{formatFieldLabel(entry.key)}:</strong> {entry.value}
                              </li>
                            ))}
                          </ul>
                        ) : (
                          <p>Tidak ada detail tambahan.</p>
                        )}
                      </span>
                    </span>
                  </td>
                  <td>{rxPower}</td>
                  <td>{txPower}</td>
                  <td>{toReadable(item.status)}</td>
                  <td>
                    <Link to={`/onu/${boardId}/${ponId}/${onuId}`}>Detail</Link>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}
