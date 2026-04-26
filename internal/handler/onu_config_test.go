package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
)

type mockONUConfigUsecase struct {
	executeFunc func(ctx context.Context, req model.ONUExecuteRequest) (model.ONUExecuteResult, error)
}

func (m *mockONUConfigUsecase) Execute(ctx context.Context, req model.ONUExecuteRequest) (model.ONUExecuteResult, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, req)
	}
	return model.ONUExecuteResult{Mode: "dry_run"}, nil
}

type mockOLTCRUDUsecase struct {
	createFunc     func(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error)
	updateFunc     func(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error)
	softDeleteFunc func(ctx context.Context, oltID string) error
	listFunc       func(ctx context.Context) []model.OLTProfile
}

func (m *mockOLTCRUDUsecase) Create(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, oltID, req)
	}
	return model.OLTProfile{ID: oltID}, nil
}
func (m *mockOLTCRUDUsecase) Update(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, oltID, req)
	}
	return model.OLTProfile{ID: oltID}, nil
}
func (m *mockOLTCRUDUsecase) SoftDelete(ctx context.Context, oltID string) error {
	if m.softDeleteFunc != nil {
		return m.softDeleteFunc(ctx, oltID)
	}
	return nil
}
func (m *mockOLTCRUDUsecase) List(ctx context.Context) []model.OLTProfile {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil
}
func (m *mockOLTCRUDUsecase) GetActive(context.Context, string) (model.OLTProfile, error) {
	return model.OLTProfile{}, apperrors.NewNotFoundError("OLT", "x")
}

func TestOnuConfigHandler_ExecuteForbiddenForViewer(t *testing.T) {
	h := NewOnuConfigHandler(&mockONUConfigUsecase{}, &mockOLTCRUDUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/onu/config/execute", bytes.NewBufferString(`{"olt_id":"o1","operation":"add","target":"ont_config"}`))
	req.Header.Set("X-Role", "viewer")
	rr := httptest.NewRecorder()

	h.Execute(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestOnuConfigHandler_ExecuteSuccessForOperator(t *testing.T) {
	h := NewOnuConfigHandler(&mockONUConfigUsecase{
		executeFunc: func(_ context.Context, req model.ONUExecuteRequest) (model.ONUExecuteResult, error) {
			if req.Operation != "add" {
				t.Fatalf("expected add operation, got %s", req.Operation)
			}
			return model.ONUExecuteResult{Mode: "execute", Applied: true}, nil
		},
	}, &mockOLTCRUDUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/onu/config/execute", bytes.NewBufferString(`{"olt_id":"o1","operation":"add","target":"ont_config","onu_ref":{"slot":1,"pon":1,"onu_id":1},"payload":{"x":1}}`))
	req.Header.Set("X-Role", "operator")
	rr := httptest.NewRecorder()

	h.Execute(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestOnuConfigHandler_CreateOLTAdminOnly(t *testing.T) {
	h := NewOnuConfigHandler(&mockONUConfigUsecase{}, &mockOLTCRUDUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/olts", bytes.NewBufferString(`{"olt_id":"olt-1","host":"10.0.0.1"}`))
	req.Header.Set("X-Role", "operator")
	rr := httptest.NewRecorder()

	h.CreateOLT(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}
