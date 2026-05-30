package httptransport

import (
	"net/http"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/service"
	apperrors "github.com/runagent-robotics/run-factory/runfactory-core/internal/platform/errors"
	httputil "github.com/runagent-robotics/run-factory/runfactory-core/internal/platform/http"
)

type Handler struct {
	service service.FactoryService
}

func NewHandler(service service.FactoryService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateFactory(w http.ResponseWriter, r *http.Request) {
	var req createFactoryRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	factory, err := h.service.CreateFactory(r.Context(), service.CreateFactoryInput{
		ID:    req.ID,
		Name:  req.Name,
		Map3D: req.Map3D,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, toFactoryResponse(factory))
}

func (h *Handler) ListFactories(w http.ResponseWriter, r *http.Request) {
	factories, err := h.service.ListFactories(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]factoryResponse, 0, len(factories))
	for _, factory := range factories {
		response = append(response, toFactoryResponse(factory))
	}
	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetFactory(w http.ResponseWriter, r *http.Request) {
	factory, err := h.service.GetFactory(r.Context(), r.PathValue("factoryID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, toFactoryResponse(factory))
}

func (h *Handler) UpdateFactoryMap(w http.ResponseWriter, r *http.Request) {
	var req updateMapRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	factory, err := h.service.UpdateFactoryMap(r.Context(), r.PathValue("factoryID"), req.Map3D)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, toFactoryResponse(factory))
}

func (h *Handler) AddRobot(w http.ResponseWriter, r *http.Request) {
	var req addRobotRequest
	if err := httputil.DecodeJSONBody(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	factory, err := h.service.AddRobot(r.Context(), r.PathValue("factoryID"), service.AddRobotInput{
		ID:       req.ID,
		Name:     req.Name,
		Position: req.Position,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, toFactoryResponse(factory))
}

func (h *Handler) RemoveRobot(w http.ResponseWriter, r *http.Request) {
	factory, err := h.service.RemoveRobot(r.Context(), r.PathValue("factoryID"), r.PathValue("robotID"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, toFactoryResponse(factory))
}

func writeServiceError(w http.ResponseWriter, err error) {
	httputil.WriteError(w, apperrors.HTTPStatus(err), apperrors.MessageOf(err))
}
