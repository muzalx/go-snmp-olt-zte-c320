package usecase

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
)

type mockOLTUsecase struct {
	getActiveFunc func(ctx context.Context, oltID string) (model.OLTProfile, error)
}

func (m *mockOLTUsecase) Create(context.Context, string, model.OLTProfileRequest) (model.OLTProfile, error) {
	return model.OLTProfile{}, nil
}
func (m *mockOLTUsecase) Update(context.Context, string, model.OLTProfileRequest) (model.OLTProfile, error) {
	return model.OLTProfile{}, nil
}
func (m *mockOLTUsecase) SoftDelete(context.Context, string) error { return nil }
func (m *mockOLTUsecase) List(context.Context) []model.OLTProfile  { return nil }
func (m *mockOLTUsecase) GetActive(ctx context.Context, oltID string) (model.OLTProfile, error) {
	if m.getActiveFunc != nil {
		return m.getActiveFunc(ctx, oltID)
	}
	return model.OLTProfile{ID: oltID, Status: "active"}, nil
}

func TestONUConfigUsecaseExecute_DryRunSuccess(t *testing.T) {
	uc := NewONUConfigUsecase(&mockOLTUsecase{})
	res, err := uc.Execute(context.Background(), model.ONUExecuteRequest{
		OLTID:     "olt-jkt-1",
		Operation: "ADD",
		Target:    "ONT_CONFIG",
		DryRun:    true,
		OnuRef:    model.ONUReference{Slot: 1, PON: 1, OnuID: 1},
		Payload:   map[string]any{"device_info": map[string]any{"name": "ONU-1"}},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.Mode != "dry_run" || res.Applied {
		t.Fatalf("unexpected result: %+v", res)
	}
	if len(res.Plan) != 3 {
		t.Fatalf("expected operation plan, got %+v", res.Plan)
	}
}

func TestONUConfigUsecaseExecute_DeleteWithoutConfirm(t *testing.T) {
	uc := NewONUConfigUsecase(&mockOLTUsecase{})
	_, err := uc.Execute(context.Background(), model.ONUExecuteRequest{
		OLTID:     "olt-jkt-1",
		Operation: "delete",
		Target:    "interface",
		OnuRef:    model.ONUReference{Slot: 1, PON: 1, OnuID: 1},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %T", err)
	}
	if appErr.Type != apperrors.ErrorTypeValidation {
		t.Fatalf("expected validation error type, got %s", appErr.Type)
	}
}

func TestONUConfigUsecaseExecute_RequirePayloadOnEdit(t *testing.T) {
	uc := NewONUConfigUsecase(&mockOLTUsecase{})
	_, err := uc.Execute(context.Background(), model.ONUExecuteRequest{
		OLTID:     "olt-jkt-1",
		Operation: "edit",
		Target:    "registered",
		OnuRef:    model.ONUReference{Slot: 1, PON: 1, OnuID: 1},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestONUConfigUsecaseExecute_RequireActiveOLT(t *testing.T) {
	uc := NewONUConfigUsecase(&mockOLTUsecase{
		getActiveFunc: func(_ context.Context, _ string) (model.OLTProfile, error) {
			return model.OLTProfile{}, apperrors.NewNotFoundError("OLT profile", "olt-1")
		},
	})
	_, err := uc.Execute(context.Background(), model.ONUExecuteRequest{
		OLTID:     "olt-1",
		Operation: "add",
		Target:    "ont_config",
		OnuRef:    model.ONUReference{Slot: 1, PON: 1, OnuID: 1},
		Payload:   map[string]any{"foo": "bar"},
	})
	if err == nil {
		t.Fatal("expected not found error")
	}
}
