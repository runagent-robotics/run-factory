package httptransport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository/memory"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/service"
)

func TestFactoryHTTPContractLifecycle(t *testing.T) {
	h := newTestRouter()

	createResp := requestJSON(t, h, http.MethodPost, "/factories", map[string]any{
		"id":   "factory-a",
		"name": "Factory A",
		"map3d": map[string]any{
			"mesh": "factory-a.glb",
		},
	})
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createResp.Code)
	}

	addRobotResp := requestJSON(t, h, http.MethodPost, "/factories/factory-a/robots", map[string]any{
		"id":   "robot-1",
		"name": "Robot 1",
		"position": map[string]any{
			"x": 1.25,
			"y": 2.5,
			"z": 3.75,
		},
	})
	if addRobotResp.Code != http.StatusCreated {
		t.Fatalf("expected add robot status %d, got %d", http.StatusCreated, addRobotResp.Code)
	}

	updateMapResp := requestJSON(t, h, http.MethodPut, "/factories/factory-a/map", map[string]any{
		"map3d": map[string]any{
			"mesh": "factory-a-v2.glb",
		},
	})
	if updateMapResp.Code != http.StatusOK {
		t.Fatalf("expected update map status %d, got %d", http.StatusOK, updateMapResp.Code)
	}

	getResp := request(t, h, http.MethodGet, "/factories/factory-a", nil)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get factory status %d, got %d", http.StatusOK, getResp.Code)
	}

	var factory factoryResponse
	decodeResponse(t, getResp, &factory)
	if factory.ID != "factory-a" || factory.Name != "Factory A" {
		t.Fatalf("unexpected factory payload: %+v", factory)
	}
	if len(factory.Robots) != 1 || factory.Robots[0].ID != "robot-1" {
		t.Fatalf("unexpected robots payload: %+v", factory.Robots)
	}

	var mapPayload map[string]string
	if err := json.Unmarshal(factory.Map3D, &mapPayload); err != nil {
		t.Fatalf("failed to decode map3d payload: %v", err)
	}
	if mapPayload["mesh"] != "factory-a-v2.glb" {
		t.Fatalf("expected updated map mesh, got %q", mapPayload["mesh"])
	}

	removeResp := request(t, h, http.MethodDelete, "/factories/factory-a/robots/robot-1", nil)
	if removeResp.Code != http.StatusOK {
		t.Fatalf("expected remove robot status %d, got %d", http.StatusOK, removeResp.Code)
	}

	listResp := request(t, h, http.MethodGet, "/factories", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listResp.Code)
	}
	var factories []factoryResponse
	decodeResponse(t, listResp, &factories)
	if len(factories) != 1 {
		t.Fatalf("expected 1 factory, got %d", len(factories))
	}
	if len(factories[0].Robots) != 0 {
		t.Fatalf("expected robots to be empty after remove, got %+v", factories[0].Robots)
	}
}

func TestFactoryHTTPValidationAndNotFound(t *testing.T) {
	h := newTestRouter()

	missingMapResp := requestJSON(t, h, http.MethodPost, "/factories", map[string]any{
		"id":   "factory-b",
		"name": "Factory B",
	})
	if missingMapResp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, missingMapResp.Code)
	}

	addMissingFactoryResp := requestJSON(t, h, http.MethodPost, "/factories/nope/robots", map[string]any{"id": "robot-1"})
	if addMissingFactoryResp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, addMissingFactoryResp.Code)
	}
}

func newTestRouter() http.Handler {
	repo := memory.NewFactoryStore()
	svc := service.NewFactoryService(repo)
	return NewRouter(NewHandler(svc))
}

func requestJSON(t *testing.T, h http.Handler, method, path string, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request payload: %v", err)
	}
	return request(t, h, method, path, bytes.NewReader(payload))
}

func request(t *testing.T, h http.Handler, method, path string, body *bytes.Reader) *httptest.ResponseRecorder {
	t.Helper()
	if body == nil {
		body = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, body)
	if body.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	return res
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
