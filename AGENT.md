# AGENT.md — go-snmp-olt-zte-c320

Dokumen ini adalah **quick-operational map** untuk agent agar bisa eksekusi task dengan cepat, konsisten, dan hemat token.

## SYSTEM ROLE
Anda adalah engineering agent untuk service **Monitoring OLT ZTE C320 via SNMP**. Fokus utama:
1. Menjaga reliability endpoint read-oriented untuk data ONU.
2. Menjaga kontrak API (format response/error + request_id).
3. Menjaga performa (Redis cache + SNMP pooling + dedup/singleflight).
4. Menjaga observability (structured logging + metrics + health).

## Project Stack
- **Language**: Go 1.26 (`go.mod`)
- **HTTP**: `chi/v5`
- **SNMP**: `gosnmp`
- **Cache/Store**: `redis/go-redis/v9`
- **Scheduler**: `robfig/cron/v3`
- **Logging**: `zap`
- **Metrics**: Prometheus `client_golang`
- **Dev tooling**: Taskfile, Docker Compose, Air, k6

## Workflow & Isolation (Critical)
1. **Mulai dari peta ini**, jangan langsung scan repo penuh.
2. **Baca file minimum** sesuai task (lihat “Core Code Structure Map”).
3. **Isolasi scope perubahan**: ubah hanya package relevan.
4. **Jangan ubah kontrak lintas-layer** tanpa alasan jelas:
   - handler ↔ usecase ↔ repository
   - format response/error
   - request_id propagation
5. **Prefer small atomic patch** + test package terkait dulu, baru `go test ./...` jika perlu.
6. Jika task menyentuh konfigurasi/runtime, validasi juga pada:
   - `.env.example`
   - `README.md`
   - `api/openapi.yaml` (jika kontrak API berubah)

## Core Directives
- Pertahankan format response JSON yang sudah distandarkan.
- Pertahankan middleware chain keamanan + audit + request_id.
- Semua error HTTP lewat helper utilitas yang sudah ada (hindari format custom per-handler).
- Jangan introduce dependency berat tanpa kebutuhan kuat.
- Kompatibilitas backward untuk endpoint existing adalah default.

## Code Quality Standards
- Ikuti idiom Go: small functions, early return, error wrapping yang jelas.
- Wajib ada test untuk perubahan logic (minimal unit test pada package terkait).
- Hindari side effect global yang membuat test flaky.
- Gunakan structured logging (`zap`) dan field snake_case.
- Hindari magic number; gunakan konstanta/config.
- Jika menyentuh concurrent flow: cek race risk dan lock scope.

## Core Code Structure Map (Token-Saving Entry Points)
Baca berurutan sesuai kebutuhan:

### 1) Bootstrapping & app wiring
- `cmd/api/main.go` → entrypoint
- `app/app.go` → inisialisasi dependency (Redis, SNMP, usecase, handler, health)
- `app/routes.go` → router + middleware + route registration

### 2) Domain flow utama (ONU)
- `internal/handler/onu.go` → HTTP handler
- `internal/usecase/onu.go` → business logic
- `internal/repository/snmp.go` + `internal/repository/redis.go` → data access
- `internal/model/onu.go` → model utama

### 3) Cross-cutting wajib
- `internal/middleware/*` → auth, request_id, logger, audit, validation, security
- `internal/utils/response.go` + `internal/utils/error.go` → response/error helpers
- `internal/errors/errors.go` → typed app error
- `pkg/logger/logger.go` → logger wrapper
- `pkg/metrics/prometheus.go` → metric instrumentation

### 4) Health, trap, dan background behavior
- `internal/health/health.go`
- `internal/trap/listener.go`, `internal/trap/handler.go`, `internal/trap/power_monitor.go`, `internal/trap/webhook.go`
- `internal/usecase/prewarm.go`

### 5) Config & contract
- `config/config.go`
- `.env.example`
- `api/openapi.yaml`
- `README.md`

## Fast Task Routing (What to open first)
- **Tambah/ubah endpoint**: `app/routes.go` → `internal/handler/onu.go` → `internal/usecase/onu.go` → test terkait.
- **Bug response format/error**: `internal/utils/response.go` + `internal/utils/error.go` + `internal/errors/errors.go`.
- **Bug auth/header/request_id**: `internal/middleware/auth.go`, `requestid.go`, `logger.go`, `audit.go`.
- **Bug SNMP data**: `internal/repository/snmp.go` + `pkg/snmp/snmp.go` + `config/oid_generator.go`.
- **Bug cache**: `internal/repository/redis.go` + `pkg/redis/redis.go` + prewarm/usecase.
- **Metrics/health**: `pkg/metrics/prometheus.go` + `internal/health/health.go`.

## Validation Checklist Before Finish
1. Build/test minimal scope lulus.
2. Tidak ada kontrak API existing yang rusak tanpa dokumentasi.
3. Jika endpoint berubah: update OpenAPI + README + test.
4. Logging dan metrics tetap konsisten.
5. Tidak ada hardcoded secret/credential.

## Recommended Commands
- Unit test semua: `go test ./... -cover`
- Unit test package spesifik: `go test ./internal/usecase -v -cover`
- Race check: `go test ./... -race`
- Local dev: `task dev`
- Lihat task tersedia: `task --list`

## Non-Goals / Avoid
- Refactor besar lintas package tanpa diminta.
- Mengubah style API secara global tanpa migration path.
- Menambah abstraction baru jika belum ada kebutuhan nyata.

## Notes for Future Agents
Jika prompt user ambigu, prioritaskan:
1. correctness kontrak API,
2. kestabilan runtime,
3. observability,
4. efisiensi perubahan.
