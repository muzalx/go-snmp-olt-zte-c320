package usecase

import (
	"context"
	"strings"
	"sync"
	"time"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
)

// OLTUseCaseInterface manages OLT profiles for multi-OLT operations.
type OLTUseCaseInterface interface {
	Create(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error)
	Update(ctx context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error)
	SoftDelete(ctx context.Context, oltID string) error
	List(ctx context.Context) []model.OLTProfile
	GetActive(ctx context.Context, oltID string) (model.OLTProfile, error)
}

type oltUsecase struct {
	mu   sync.RWMutex
	data map[string]model.OLTProfile
}

func NewOLTUsecase() OLTUseCaseInterface {
	return &oltUsecase{data: make(map[string]model.OLTProfile)}
}

func (u *oltUsecase) Create(_ context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error) {
	if strings.TrimSpace(oltID) == "" {
		return model.OLTProfile{}, apperrors.NewValidationError("olt_id is required", nil)
	}
	if strings.TrimSpace(req.Host) == "" {
		return model.OLTProfile{}, apperrors.NewValidationError("host is required", nil)
	}
	now := time.Now().UTC()
	profile := model.OLTProfile{
		ID:        oltID,
		Name:      req.Name,
		Host:      req.Host,
		Port:      req.Port,
		Community: req.Community,
		TimeoutMS: req.TimeoutMS,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	if _, exists := u.data[oltID]; exists {
		return model.OLTProfile{}, apperrors.NewValidationError("olt_id already exists", map[string]any{"olt_id": oltID})
	}
	u.data[oltID] = profile
	return profile, nil
}

func (u *oltUsecase) Update(_ context.Context, oltID string, req model.OLTProfileRequest) (model.OLTProfile, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	profile, exists := u.data[oltID]
	if !exists {
		return model.OLTProfile{}, apperrors.NewNotFoundError("OLT profile", oltID)
	}
	profile.Name = req.Name
	if req.Host != "" {
		profile.Host = req.Host
	}
	if req.Port != 0 {
		profile.Port = req.Port
	}
	if req.Community != "" {
		profile.Community = req.Community
	}
	if req.TimeoutMS != 0 {
		profile.TimeoutMS = req.TimeoutMS
	}
	profile.UpdatedAt = time.Now().UTC()
	u.data[oltID] = profile
	return profile, nil
}

func (u *oltUsecase) SoftDelete(_ context.Context, oltID string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	profile, exists := u.data[oltID]
	if !exists {
		return apperrors.NewNotFoundError("OLT profile", oltID)
	}
	profile.Status = "deleted_soft"
	profile.UpdatedAt = time.Now().UTC()
	u.data[oltID] = profile
	return nil
}

func (u *oltUsecase) List(_ context.Context) []model.OLTProfile {
	u.mu.RLock()
	defer u.mu.RUnlock()
	result := make([]model.OLTProfile, 0, len(u.data))
	for _, profile := range u.data {
		result = append(result, profile)
	}
	return result
}

func (u *oltUsecase) GetActive(_ context.Context, oltID string) (model.OLTProfile, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	profile, exists := u.data[oltID]
	if !exists {
		return model.OLTProfile{}, apperrors.NewNotFoundError("OLT profile", oltID)
	}
	if profile.Status != "active" {
		return model.OLTProfile{}, apperrors.NewValidationError("olt profile is not active", map[string]any{"olt_id": oltID, "status": profile.Status})
	}
	return profile, nil
}
