package memory

import (
	"context"
	"encoding/json"
	"slices"
	"sync"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository"
)

type FactoryStore struct {
	mu        sync.RWMutex
	factories map[string]domain.Factory
}

func NewFactoryStore() *FactoryStore {
	return &FactoryStore{
		factories: make(map[string]domain.Factory),
	}
}

func (s *FactoryStore) Create(_ context.Context, factory domain.Factory) (domain.Factory, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.factories[factory.ID]; exists {
		return domain.Factory{}, repository.ErrFactoryExists
	}

	factory.Robots = make(map[string]domain.Robot)
	cloned := domain.CloneFactory(factory)
	s.factories[factory.ID] = cloned
	return domain.CloneFactory(cloned), nil
}

func (s *FactoryStore) List(_ context.Context) ([]domain.Factory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]domain.Factory, 0, len(s.factories))
	for _, factory := range s.factories {
		items = append(items, domain.CloneFactory(factory))
	}

	slices.SortFunc(items, func(a, b domain.Factory) int {
		switch {
		case a.ID < b.ID:
			return -1
		case a.ID > b.ID:
			return 1
		default:
			return 0
		}
	})

	return items, nil
}

func (s *FactoryStore) Get(_ context.Context, id string) (domain.Factory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	factory, found := s.factories[id]
	if !found {
		return domain.Factory{}, repository.ErrFactoryNotFound
	}

	return domain.CloneFactory(factory), nil
}

func (s *FactoryStore) UpdateMap(_ context.Context, id string, map3D json.RawMessage) (domain.Factory, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	factory, found := s.factories[id]
	if !found {
		return domain.Factory{}, repository.ErrFactoryNotFound
	}

	factory.Map3D = domain.CloneRawMessage(map3D)
	s.factories[id] = factory
	return domain.CloneFactory(factory), nil
}

func (s *FactoryStore) AddRobot(_ context.Context, factoryID string, robot domain.Robot) (domain.Factory, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	factory, found := s.factories[factoryID]
	if !found {
		return domain.Factory{}, repository.ErrFactoryNotFound
	}
	if _, exists := factory.Robots[robot.ID]; exists {
		return domain.Factory{}, repository.ErrRobotExists
	}

	factory.Robots[robot.ID] = domain.CloneRobot(robot)
	s.factories[factoryID] = factory
	return domain.CloneFactory(factory), nil
}

func (s *FactoryStore) RemoveRobot(_ context.Context, factoryID, robotID string) (domain.Factory, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	factory, found := s.factories[factoryID]
	if !found {
		return domain.Factory{}, repository.ErrFactoryNotFound
	}
	if _, exists := factory.Robots[robotID]; !exists {
		return domain.Factory{}, repository.ErrRobotNotFound
	}

	delete(factory.Robots, robotID)
	s.factories[factoryID] = factory
	return domain.CloneFactory(factory), nil
}

var _ repository.FactoryRepository = (*FactoryStore)(nil)
