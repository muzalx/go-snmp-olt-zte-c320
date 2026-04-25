# Frontend - OLT ZTE C320 Monitor

## Development
```bash
cd frontend
npm install
npm run dev
```

By default app runs at `http://localhost:5173`.

## Required Environment Variables
Copy `.env.example` to `.env`.

- `VITE_API_BASE_URL` (example: `http://localhost:8081`)
- `VITE_API_KEY` (optional)

## Build
```bash
cd frontend
npm run build
npm run preview
```

## Docker Build
```bash
docker build \
  --build-arg VITE_API_BASE_URL=http://localhost:8081 \
  --build-arg VITE_API_KEY= \
  -t go-snmp-olt-zte-c320-frontend:local \
  ./frontend
```

## Docker Run
```bash
docker run --rm -p 5173:80 go-snmp-olt-zte-c320-frontend:local
```
