import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "@/components/common/app-layout";
import { HealthPage } from "@/features/health/pages/health-page";
import { OnuExplorerPage } from "@/features/onu/pages/onu-explorer-page";
import { OnuDetailPage } from "@/features/onu/pages/onu-detail-page";
import { CacheToolsPage } from "@/features/cache/components/cache-tools-page";

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route path="/" element={<Navigate to="/health" replace />} />
          <Route path="/health" element={<HealthPage />} />
          <Route path="/onu" element={<OnuExplorerPage />} />
          <Route path="/onu/:boardId/:ponId/:onuId" element={<OnuDetailPage />} />
          <Route path="/cache" element={<CacheToolsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
