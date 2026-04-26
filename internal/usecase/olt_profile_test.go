package usecase

import (
	"context"
	"testing"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
)

func TestOLTUsecase_CreateUpdateListSoftDelete(t *testing.T) {
	uc := NewOLTUsecase()

	created, err := uc.Create(context.Background(), "olt-1", model.OLTProfileRequest{
		Name:      "DC1 OLT",
		Host:      "10.0.0.1",
		Port:      161,
		Community: "public",
		TimeoutMS: 1000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != "olt-1" {
		t.Fatalf("unexpected id: %+v", created)
	}

	updated, err := uc.Update(context.Background(), "olt-1", model.OLTProfileRequest{
		Name: "DC1 OLT Updated",
		Host: "10.0.0.2",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "DC1 OLT Updated" || updated.Host != "10.0.0.2" {
		t.Fatalf("unexpected update: %+v", updated)
	}

	all := uc.List(context.Background())
	if len(all) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(all))
	}

	if err = uc.SoftDelete(context.Background(), "olt-1"); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	_, err = uc.GetActive(context.Background(), "olt-1")
	if err == nil {
		t.Fatal("expected inactive profile error")
	}
}

func TestOLTUsecase_CreateRequireHost(t *testing.T) {
	uc := NewOLTUsecase()
	_, err := uc.Create(context.Background(), "olt-1", model.OLTProfileRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if appErr, ok := err.(*apperrors.AppError); !ok || appErr.Type != apperrors.ErrorTypeValidation {
		t.Fatalf("expected validation app error, got %T (%v)", err, err)
	}
}
