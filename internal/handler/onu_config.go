package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	apperrors "github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/errors"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/usecase"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/utils"
	"github.com/go-chi/chi/v5"
)

// OnuConfigHandler serves generic ONU config execute and OLT management endpoints.
type OnuConfigHandler struct {
	onuConfigUsecase usecase.ONUConfigUseCaseInterface
	oltUsecase       usecase.OLTUseCaseInterface
}

func NewOnuConfigHandler(onuConfigUsecase usecase.ONUConfigUseCaseInterface, oltUsecase usecase.OLTUseCaseInterface) *OnuConfigHandler {
	return &OnuConfigHandler{onuConfigUsecase: onuConfigUsecase, oltUsecase: oltUsecase}
}

func (h *OnuConfigHandler) Execute(w http.ResponseWriter, r *http.Request) {
	var req model.ONUExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, r, apperrors.NewValidationError("invalid JSON body", map[string]any{"error": err.Error()}))
		return
	}

	if !allowONUOperation(r.Header.Get("X-Role"), req.Operation) {
		writeForbidden(w, r, "insufficient role for operation")
		return
	}

	result, err := h.onuConfigUsecase.Execute(r.Context(), req)
	if err != nil {
		utils.HandleError(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{Code: http.StatusOK, Status: "success", Data: result})
}

func (h *OnuConfigHandler) CreateOLT(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Header.Get("X-Role")) {
		writeForbidden(w, r, "admin role is required")
		return
	}
	var body struct {
		OLTID string `json:"olt_id"`
		model.OLTProfileRequest
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.HandleError(w, r, apperrors.NewValidationError("invalid JSON body", map[string]any{"error": err.Error()}))
		return
	}
	created, err := h.oltUsecase.Create(r.Context(), body.OLTID, body.OLTProfileRequest)
	if err != nil {
		utils.HandleError(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, utils.WebResponse{Code: http.StatusCreated, Status: "success", Data: created})
}

func (h *OnuConfigHandler) UpdateOLT(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Header.Get("X-Role")) {
		writeForbidden(w, r, "admin role is required")
		return
	}
	oltID := chi.URLParam(r, "olt_id")
	var req model.OLTProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, r, apperrors.NewValidationError("invalid JSON body", map[string]any{"error": err.Error()}))
		return
	}
	updated, err := h.oltUsecase.Update(r.Context(), oltID, req)
	if err != nil {
		utils.HandleError(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{Code: http.StatusOK, Status: "success", Data: updated})
}

func (h *OnuConfigHandler) DeleteOLT(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r.Header.Get("X-Role")) {
		writeForbidden(w, r, "admin role is required")
		return
	}
	oltID := chi.URLParam(r, "olt_id")
	if err := h.oltUsecase.SoftDelete(r.Context(), oltID); err != nil {
		utils.HandleError(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{Code: http.StatusOK, Status: "success", Data: map[string]string{"message": "OLT soft deleted"}})
}

func (h *OnuConfigHandler) ListOLT(w http.ResponseWriter, r *http.Request) {
	list := h.oltUsecase.List(r.Context())
	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{Code: http.StatusOK, Status: "success", Data: list})
}

func allowONUOperation(role, operation string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	operation = strings.ToLower(strings.TrimSpace(operation))
	if role == "admin" {
		return true
	}
	if role == "operator" {
		return operation == "add" || operation == "edit"
	}
	if role == "viewer" {
		return false
	}
	return false
}

func isAdmin(role string) bool {
	return strings.EqualFold(strings.TrimSpace(role), "admin")
}

func writeForbidden(w http.ResponseWriter, r *http.Request, msg string) {
	utils.SendJSONResponse(w, http.StatusForbidden, utils.ErrorResponse{
		Code:      http.StatusForbidden,
		Status:    "Forbidden",
		ErrorCode: "FORBIDDEN",
		Data:      msg,
		RequestID: utils.RequestIDFromContext(r.Context()),
	})
}
