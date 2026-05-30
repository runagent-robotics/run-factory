package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
)

type Position3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Robot struct {
	ID       string      `json:"id"`
	Name     string      `json:"name,omitempty"`
	Position *Position3D `json:"position,omitempty"`
}

type factoryRecord struct {
	ID     string
	Name   string
	Map3D  json.RawMessage
	Robots map[string]Robot
}

type factoryResponse struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Map3D  json.RawMessage `json:"map3d,omitempty"`
	Robots []Robot         `json:"robots"`
}

type createFactoryRequest struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Map3D json.RawMessage `json:"map3d"`
}

type updateMapRequest struct {
	Map3D json.RawMessage `json:"map3d"`
}

type addRobotRequest struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Position *Position3D `json:"position"`
}

var (
	errFactoryExists   = errors.New("factory already exists")
	errFactoryNotFound = errors.New("factory not found")
	errRobotExists     = errors.New("robot already exists")
	errRobotNotFound   = errors.New("robot not found")
)

type factoryStore struct {
	mu        sync.RWMutex
	factories map[string]*factoryRecord
}

func newFactoryStore() *factoryStore {
	return &factoryStore{factories: make(map[string]*factoryRecord)}
}

func (s *factoryStore) create(req createFactoryRequest) (*factoryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.factories[req.ID]; exists {
		return nil, errFactoryExists
	}

	record := &factoryRecord{
		ID:     req.ID,
		Name:   req.Name,
		Map3D:  cloneRawMessage(req.Map3D),
		Robots: make(map[string]Robot),
	}
	s.factories[req.ID] = record

	return record.clone(), nil
}

func (s *factoryStore) list() []*factoryRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*factoryRecord, 0, len(s.factories))
	for _, factory := range s.factories {
		items = append(items, factory.clone())
	}

	slices.SortFunc(items, func(a, b *factoryRecord) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return items
}

func (s *factoryStore) get(id string) (*factoryRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, found := s.factories[id]
	if !found {
		return nil, errFactoryNotFound
	}

	return record.clone(), nil
}

func (s *factoryStore) updateMap(id string, map3D json.RawMessage) (*factoryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, found := s.factories[id]
	if !found {
		return nil, errFactoryNotFound
	}

	record.Map3D = cloneRawMessage(map3D)
	return record.clone(), nil
}

func (s *factoryStore) addRobot(factoryID string, robot Robot) (*factoryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, found := s.factories[factoryID]
	if !found {
		return nil, errFactoryNotFound
	}
	if _, exists := record.Robots[robot.ID]; exists {
		return nil, errRobotExists
	}

	record.Robots[robot.ID] = robot
	return record.clone(), nil
}

func (s *factoryStore) removeRobot(factoryID, robotID string) (*factoryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, found := s.factories[factoryID]
	if !found {
		return nil, errFactoryNotFound
	}
	if _, exists := record.Robots[robotID]; !exists {
		return nil, errRobotNotFound
	}

	delete(record.Robots, robotID)
	return record.clone(), nil
}

func (f *factoryRecord) clone() *factoryRecord {
	copyRecord := &factoryRecord{
		ID:     f.ID,
		Name:   f.Name,
		Map3D:  cloneRawMessage(f.Map3D),
		Robots: make(map[string]Robot, len(f.Robots)),
	}
	for id, robot := range f.Robots {
		copyRecord.Robots[id] = robot
	}
	return copyRecord
}

func (f *factoryRecord) toResponse() factoryResponse {
	robots := make([]Robot, 0, len(f.Robots))
	for _, robot := range f.Robots {
		robots = append(robots, robot)
	}
	slices.SortFunc(robots, func(a, b Robot) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})

	return factoryResponse{
		ID:     f.ID,
		Name:   f.Name,
		Map3D:  cloneRawMessage(f.Map3D),
		Robots: robots,
	}
}

type apiServer struct {
	store *factoryStore
}

func newServer() http.Handler {
	api := &apiServer{store: newFactoryStore()}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /factories", api.createFactory)
	mux.HandleFunc("GET /factories", api.listFactories)
	mux.HandleFunc("GET /factories/{factoryID}", api.getFactory)
	mux.HandleFunc("PUT /factories/{factoryID}/map", api.updateFactoryMap)
	mux.HandleFunc("POST /factories/{factoryID}/robots", api.addRobot)
	mux.HandleFunc("DELETE /factories/{factoryID}/robots/{robotID}", api.removeRobot)

	return mux
}

func main() {
	addr := ":8080"
	fmt.Printf("runfactory-core listening on %s\n", addr)
	if err := http.ListenAndServe(addr, newServer()); err != nil {
		panic(err)
	}
}

func (a *apiServer) createFactory(w http.ResponseWriter, r *http.Request) {
	var req createFactoryRequest
	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ID == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "id and name are required")
		return
	}
	if len(req.Map3D) == 0 {
		writeError(w, http.StatusBadRequest, "map3d is required")
		return
	}

	record, err := a.store.create(req)
	if err != nil {
		if errors.Is(err, errFactoryExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, record.toResponse())
}

func (a *apiServer) listFactories(w http.ResponseWriter, _ *http.Request) {
	records := a.store.list()
	response := make([]factoryResponse, 0, len(records))
	for _, record := range records {
		response = append(response, record.toResponse())
	}
	writeJSON(w, http.StatusOK, response)
}

func (a *apiServer) getFactory(w http.ResponseWriter, r *http.Request) {
	record, err := a.store.get(r.PathValue("factoryID"))
	if err != nil {
		if errors.Is(err, errFactoryNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, record.toResponse())
}

func (a *apiServer) updateFactoryMap(w http.ResponseWriter, r *http.Request) {
	var req updateMapRequest
	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Map3D) == 0 {
		writeError(w, http.StatusBadRequest, "map3d is required")
		return
	}

	record, err := a.store.updateMap(r.PathValue("factoryID"), req.Map3D)
	if err != nil {
		if errors.Is(err, errFactoryNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, record.toResponse())
}

func (a *apiServer) addRobot(w http.ResponseWriter, r *http.Request) {
	var req addRobotRequest
	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ID == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	record, err := a.store.addRobot(r.PathValue("factoryID"), Robot{
		ID:       req.ID,
		Name:     req.Name,
		Position: req.Position,
	})
	if err != nil {
		switch {
		case errors.Is(err, errFactoryNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, errRobotExists):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, record.toResponse())
}

func (a *apiServer) removeRobot(w http.ResponseWriter, r *http.Request) {
	record, err := a.store.removeRobot(r.PathValue("factoryID"), r.PathValue("robotID"))
	if err != nil {
		switch {
		case errors.Is(err, errFactoryNotFound), errors.Is(err, errRobotNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, record.toResponse())
}

func decodeJSONBody(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if decoder.More() {
		return fmt.Errorf("request body must contain exactly one JSON object")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func cloneRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	copyRaw := make(json.RawMessage, len(raw))
	copy(copyRaw, raw)
	return copyRaw
}
