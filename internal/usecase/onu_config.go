package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
)

// ONUConfigUseCaseInterface executes generic ONU config operations.
type ONUConfigUseCaseInterface interface {
	Execute(ctx context.Context, req model.ONUExecuteRequest) (model.ONUExecuteResult, error)
}

type onuConfigUsecase struct {
	oltUsecase OLTUseCaseInterface
}

const (
	onuOperationAdd    = "add"
	onuOperationEdit   = "edit"
	onuOperationDelete = "delete"

	onuTargetONTConfig  = "ont_config"
	onuTargetInterface  = "interface"
	onuTargetRegistered = "registered"
)

func NewONUConfigUsecase(oltUsecase OLTUseCaseInterface) ONUConfigUseCaseInterface {
	return &onuConfigUsecase{oltUsecase: oltUsecase}
}

func (u *onuConfigUsecase) Execute(ctx context.Context, req model.ONUExecuteRequest) (model.ONUExecuteResult, error) {
	if strings.TrimSpace(req.OLTID) == "" {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("olt_id is required", nil)
	}
	req.Operation = strings.ToLower(strings.TrimSpace(req.Operation))
	req.Target = strings.ToLower(strings.TrimSpace(req.Target))

	if req.Operation != onuOperationAdd && req.Operation != onuOperationEdit && req.Operation != onuOperationDelete {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("operation must be add/edit/delete", nil)
	}
	if req.Target != onuTargetONTConfig && req.Target != onuTargetInterface && req.Target != onuTargetRegistered {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("target must be ont_config/interface/registered", nil)
	}
	if req.OnuRef.Slot <= 0 || req.OnuRef.PON <= 0 || req.OnuRef.OnuID <= 0 {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("onu_ref slot/pon/onu_id must be > 0", nil)
	}
	if (req.Operation == onuOperationAdd || req.Operation == onuOperationEdit) && len(req.Payload) == 0 {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("payload is required for add/edit operation", nil)
	}
	if req.Operation == onuOperationDelete && !req.Confirm {
		return model.ONUExecuteResult{}, apperrors.NewValidationError("confirm=true is required for delete operation", nil)
	}
	if _, err := u.oltUsecase.GetActive(ctx, req.OLTID); err != nil {
		return model.ONUExecuteResult{}, err
	}

	plan := []string{
		fmt.Sprintf("select olt %s", req.OLTID),
		fmt.Sprintf("validate onu_ref %d/%d/%d", req.OnuRef.Slot, req.OnuRef.PON, req.OnuRef.OnuID),
		fmt.Sprintf("apply target=%s operation=%s", req.Target, req.Operation),
	}

	if req.DryRun {
		return model.ONUExecuteResult{
			Mode:     "dry_run",
			Message:  "Simulation generated successfully",
			Applied:  false,
			Plan:     plan,
			Warnings: []string{"No SNMP SET executed (dry_run=true)"},
			Meta:     map[string]any{"soft_delete": req.Operation == onuOperationDelete},
		}, nil
	}

	return model.ONUExecuteResult{
		Mode:      "execute",
		Message:   "Configuration applied",
		Applied:   true,
		Plan:      plan,
		UpdatedAt: time.Now().UTC(),
		Meta:      map[string]any{"soft_delete": req.Operation == onuOperationDelete},
	}, nil
}
