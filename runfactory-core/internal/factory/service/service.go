package service

import (
	"context"
	"encoding/json"
	stderrors "errors"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository"
	apperrors "github.com/runagent-robotics/run-factory/runfactory-core/internal/platform/errors"
)

type CreateFactoryInput struct {
	ID    string
	Name  string
	Map3D json.RawMessage
}

type AddRobotInput struct {
	ID       string
	Name     string
	Position *domain.Position3D
}

type FactoryService interface {
	CreateFactory(ctx context.Context, input CreateFactoryInput) (domain.Factory, error)
	ListFactories(ctx context.Context) ([]domain.Factory, error)
	GetFactory(ctx context.Context, id string) (domain.Factory, error)
	UpdateFactoryMap(ctx context.Context, id string, map3D json.RawMessage) (domain.Factory, error)
	AddRobot(ctx context.Context, factoryID string, input AddRobotInput) (domain.Factory, error)
	RemoveRobot(ctx context.Context, factoryID, robotID string) (domain.Factory, error)
}

type factoryService struct {
	repo repository.FactoryRepository
}

func NewFactoryService(repo repository.FactoryRepository) FactoryService {
	return &factoryService{repo: repo}
}

func (s *factoryService) CreateFactory(ctx context.Context, input CreateFactoryInput) (domain.Factory, error) {
	if input.ID == "" || input.Name == "" {
		return domain.Factory{}, apperrors.New(apperrors.KindValidation, "id and name are required", nil)
	}
	if len(input.Map3D) == 0 {
		return domain.Factory{}, apperrors.New(apperrors.KindValidation, "map3d is required", nil)
	}

	factory, err := s.repo.Create(ctx, domain.Factory{
		ID:    input.ID,
		Name:  input.Name,
		Map3D: domain.CloneRawMessage(input.Map3D),
	})
	if err != nil {
		return domain.Factory{}, mapRepoError(err)
	}
	return factory, nil
}

func (s *factoryService) ListFactories(ctx context.Context) ([]domain.Factory, error) {
	factories, err := s.repo.List(ctx)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return factories, nil
}

func (s *factoryService) GetFactory(ctx context.Context, id string) (domain.Factory, error) {
	factory, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Factory{}, mapRepoError(err)
	}
	return factory, nil
}

func (s *factoryService) UpdateFactoryMap(ctx context.Context, id string, map3D json.RawMessage) (domain.Factory, error) {
	if len(map3D) == 0 {
		return domain.Factory{}, apperrors.New(apperrors.KindValidation, "map3d is required", nil)
	}

	factory, err := s.repo.UpdateMap(ctx, id, map3D)
	if err != nil {
		return domain.Factory{}, mapRepoError(err)
	}
	return factory, nil
}

func (s *factoryService) AddRobot(ctx context.Context, factoryID string, input AddRobotInput) (domain.Factory, error) {
	if input.ID == "" {
		return domain.Factory{}, apperrors.New(apperrors.KindValidation, "id is required", nil)
	}

	factory, err := s.repo.AddRobot(ctx, factoryID, domain.Robot{
		ID:       input.ID,
		Name:     input.Name,
		Position: domain.ClonePosition(input.Position),
	})
	if err != nil {
		return domain.Factory{}, mapRepoError(err)
	}
	return factory, nil
}

func (s *factoryService) RemoveRobot(ctx context.Context, factoryID, robotID string) (domain.Factory, error) {
	factory, err := s.repo.RemoveRobot(ctx, factoryID, robotID)
	if err != nil {
		return domain.Factory{}, mapRepoError(err)
	}
	return factory, nil
}

func mapRepoError(err error) error {
	switch {
	case stderrors.Is(err, repository.ErrFactoryExists), stderrors.Is(err, repository.ErrRobotExists):
		return apperrors.New(apperrors.KindConflict, err.Error(), err)
	case stderrors.Is(err, repository.ErrFactoryNotFound), stderrors.Is(err, repository.ErrRobotNotFound):
		return apperrors.New(apperrors.KindNotFound, err.Error(), err)
	default:
		return apperrors.New(apperrors.KindInternal, "internal server error", err)
	}
}
