# AGENT_FRONTEND.md

## Tujuan Dokumen
Dokumen ini adalah **blueprint sebelum implementasi** untuk membangun frontend website berbasis **ViteJS** yang terhubung ke API pada `api/openapi.yaml`.

Fokus tahap ini:
1. Menentukan plan delivery bertahap.
2. Menentukan daftar feature prioritas.
3. Menentukan rancangan kode (struktur folder, contract type, service layer) agar implementasi berikutnya konsisten.

---

## Ringkasan Kontrak API yang Menjadi Acuan
Berdasarkan `api/openapi.yaml`, frontend perlu menangani area berikut:

### Endpoint Health & Utility (tanpa auth)
- `GET /health`
- `GET /healthz`
- `GET /readyz`
- `GET /version`
- `GET /metrics` (opsional untuk UI karena format text Prometheus)

### Endpoint Domain ONU / Cache (dengan header `X-API-Key` bila diaktifkan backend)
- `GET /api/v1/board/{board_id}/pon/{pon_id}`
- `GET /api/v1/paginate/board/{board_id}/pon/{pon_id}?page=&limit=`
- `GET /api/v1/board/{board_id}/pon/{pon_id}/onu/{onu_id}`
- `GET /api/v1/board/{board_id}/pon/{pon_id}/onu_id/empty`
- `GET /api/v1/board/{board_id}/pon/{pon_id}/onu_id_sn`
- `POST /api/v1/board/{board_id}/pon/{pon_id}/onu_id/update`
- `DELETE /api/v1/board/{board_id}/pon/{pon_id}/cache/clear`

### Constraint penting dari API
- `board_id`: 1..2
- `pon_id`: 1..16
- `onu_id`: 1..128
- Success response memakai wrapper (`code`, `status`, `data`, opsional `meta`).
- Error response memakai wrapper (`code`, `status`, `error_code`, `data`, `request_id`).

---

## Rekomendasi Stack Frontend (ViteJS)
Gunakan stack berikut agar scalable:

- **Vite + React + TypeScript**
- **React Router** (navigasi halaman)
- **TanStack Query** (fetching, caching, retry, invalidate)
- **Axios** (HTTP client + interceptor)
- **Zod** (runtime validation untuk response penting)
- **Tailwind CSS** atau **shadcn/ui** (UI kit cepat)
- **Vitest + React Testing Library** (unit/component test)

> Catatan: kalau tim ingin footprint minimal, bisa pakai `fetch` native + Context, tapi untuk observability dan maintainability lebih baik TanStack Query.

---

## Plan Delivery (Sebelum Implementasi)

## Phase 0 — Inisialisasi Proyek
**Output:** scaffold siap coding.

Checklist:
- Create project: `npm create vite@latest frontend -- --template react-ts`
- Setup alias path `@/*` ke `src/*`
- Setup env:
  - `VITE_API_BASE_URL`
  - `VITE_API_KEY` (optional)
- Setup lint + format:
  - ESLint
  - Prettier

## Phase 1 — Fondasi Arsitektur Kode
**Output:** service layer, contract types, error normalizer.

Checklist:
- Buat `src/lib/http.ts` untuk axios instance + interceptor.
- Buat `src/types/api.ts` untuk generic wrapper response/error.
- Buat `src/features/.../api/*.ts` per domain endpoint.
- Buat utility `normalizeApiError()` untuk handling bentuk `data` string/object.

## Phase 2 — UI Shell + Halaman Utama
**Output:** layout dan routing dasar.

Halaman minimal:
- Dashboard Health
- ONU Explorer
- ONU Detail
- Cache Tools

## Phase 3 — Feature MVP (High Priority)
**Output:** flow operasional utama bisa dipakai NOC/ops.

1. Query ONU by Board/PON (paginated)
2. Lihat detail ONU by `onu_id`
3. Lihat empty ONU IDs
4. Lihat mapping ONU ID + serial number
5. Clear cache board/pon
6. Force refresh cache empty ONU ID

## Phase 4 — UX & Reliability
**Output:** aplikasi siap dipakai harian.

Checklist:
- Loading skeleton
- Empty state
- Retry action
- Toast notification sukses/gagal
- Debounce input
- Guard validasi ID sebelum request

## Phase 5 — Quality Gate
**Output:** siap handover.

Checklist:
- Unit test service + util minimal
- Component test untuk tabel/list utama
- E2E smoke (opsional)
- Build production check

---

## Daftar Feature & Kebutuhan Kode

### 1) Feature: Health Dashboard
**Tujuan:** memonitor status service cepat.

Kebutuhan kode:
- Service:
  - `getHealth()`
  - `getHealthz()`
  - `getReadyz()`
  - `getVersion()`
- UI:
  - Card status (`healthy`, `ready`, `not_ready`)
  - Dependency grid (redis/snmp)
  - Build info panel (version, commit, build_time, uptime)

### 2) Feature: ONU Explorer (Board/PON)
**Tujuan:** menampilkan daftar ONU per board/pon.

Kebutuhan kode:
- Form input board/pon dengan validasi range.
- Query endpoint paginated.
- Tabel ONU (status, serial number, power, dll sesuai payload).
- Pagination controls (`page`, `limit`).

### 3) Feature: ONU Detail
**Tujuan:** troubleshooting detail per ONU.

Kebutuhan kode:
- Route param: `/onu/:boardId/:ponId/:onuId`
- Fetch detail endpoint.
- Panel informasi detail + riwayat offline (jika tersedia di payload).

### 4) Feature: ONU ID Utilities
**Tujuan:** mempermudah provisioning dan audit.

Kebutuhan kode:
- Tab A: Empty ONU IDs.
- Tab B: ONU ID + Serial Number.
- Export CSV sederhana (nice-to-have).

### 5) Feature: Cache Control
**Tujuan:** operasi maintenance cache dari UI.

Kebutuhan kode:
- Tombol `Clear Cache` (DELETE) + confirm dialog.
- Tombol `Refresh Empty ONU ID Cache` (POST).
- Invalidate query terkait setelah aksi sukses.

### 6) Feature: Global Error Handling
**Tujuan:** menjaga UX konsisten saat API gagal.

Kebutuhan kode:
- Error adapter untuk menampilkan:
  - `error_code`
  - message dari `data`
  - `request_id` untuk troubleshooting
- Komponen `ApiErrorAlert` reusable.

---

## Rancangan Struktur Folder (Disarankan)

```txt
src/
  app/
    router.tsx
    providers.tsx
  lib/
    http.ts
    query-client.ts
    env.ts
  types/
    api.ts
    onu.ts
  features/
    health/
      api/health.api.ts
      hooks/use-health.ts
      components/health-cards.tsx
      pages/health-page.tsx
    onu/
      api/onu.api.ts
      hooks/use-onu-list.ts
      hooks/use-onu-detail.ts
      components/onu-filter-form.tsx
      components/onu-table.tsx
      components/onu-detail-panel.tsx
      pages/onu-explorer-page.tsx
      pages/onu-detail-page.tsx
    cache/
      api/cache.api.ts
      hooks/use-cache-actions.ts
      components/cache-actions.tsx
  components/
    ui/
    common/
      app-layout.tsx
      api-error-alert.tsx
      loading-state.tsx
  utils/
    normalize-api-error.ts
    validators.ts
  main.tsx
```

---

## Kontrak TypeScript (Draft Kode)

```ts
// src/types/api.ts
export type ApiSuccess<TData, TMeta = unknown> = {
  code: number;
  status: "success";
  data: TData;
  meta?: TMeta;
};

export type ApiErrorCode =
  | "VALIDATION_ERROR"
  | "NOT_FOUND"
  | "SNMP_ERROR"
  | "REDIS_ERROR"
  | "CONFIG_ERROR"
  | "INTERNAL_ERROR";

export type ApiErrorResponse = {
  code: number;
  status: string;
  error_code: ApiErrorCode;
  data: string | { message?: string; details?: Record<string, unknown> };
  request_id?: string;
};
```

```ts
// src/lib/http.ts
import axios from "axios";

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 20000,
  headers: {
    "Content-Type": "application/json",
  },
});

http.interceptors.request.use((config) => {
  const apiKey = import.meta.env.VITE_API_KEY;
  if (apiKey) {
    config.headers["X-API-Key"] = apiKey;
  }
  return config;
});
```

```ts
// src/utils/validators.ts
export const isValidBoardId = (n: number) => n >= 1 && n <= 2;
export const isValidPonId = (n: number) => n >= 1 && n <= 16;
export const isValidOnuId = (n: number) => n >= 1 && n <= 128;
```

```ts
// src/features/onu/api/onu.api.ts
import { http } from "@/lib/http";
import type { ApiSuccess } from "@/types/api";

export const getOnusPaginated = async (params: {
  boardId: number;
  ponId: number;
  page?: number;
  limit?: number;
}) => {
  const { boardId, ponId, page = 1, limit = 10 } = params;
  const res = await http.get<ApiSuccess<unknown, unknown>>(
    `/api/v1/paginate/board/${boardId}/pon/${ponId}`,
    { params: { page, limit } },
  );
  return res.data;
};
```

---

## Mapping Endpoint → UI Action

- `GET /healthz`, `GET /readyz`, `GET /version` → Dashboard Health.
- `GET /api/v1/paginate/...` → Tabel ONU + pagination.
- `GET /api/v1/board/{b}/pon/{p}/onu/{o}` → Halaman detail ONU.
- `GET /onu_id/empty` → Tab Empty IDs.
- `GET /onu_id_sn` → Tab ONU ID & SN.
- `DELETE /cache/clear` → Tombol clear cache + konfirmasi.
- `POST /onu_id/update` → Tombol refresh cache ONU ID.

---

## Acceptance Criteria (untuk lanjut implementasi)

1. User bisa input board/pon valid lalu melihat daftar ONU paginated.
2. User bisa buka detail ONU spesifik.
3. User bisa melihat empty ONU IDs dan ID+SN.
4. User bisa menjalankan aksi cache clear/refresh dengan feedback sukses/gagal.
5. Semua error API tampil konsisten beserta `request_id` bila tersedia.
6. Base URL API dan API key bisa diatur via environment.

---

## Risiko & Mitigasi

- **Risiko:** Shape field detail ONU bisa berbeda dari ekspektasi UI.
  - **Mitigasi:** pakai typing bertahap + fallback renderer + logging response di dev mode.
- **Risiko:** API key mandatory di environment tertentu.
  - **Mitigasi:** tampilkan warning banner jika `VITE_API_KEY` kosong.
- **Risiko:** Response error `data` punya 2 bentuk (string/object).
  - **Mitigasi:** wajib pakai `normalizeApiError()` sebelum ditampilkan.

---

## Next Step (Saat Mulai Implementasi)
1. Scaffold proyek Vite React TS.
2. Implement `lib/http.ts`, `types/api.ts`, dan validators.
3. Implement halaman `Health` + `ONU Explorer` terlebih dahulu (MVP core).
4. Lanjutkan ke `ONU Detail` dan `Cache Tools`.
5. Tambah testing dan hardening UX.

