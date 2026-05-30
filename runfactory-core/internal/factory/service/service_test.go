package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository"
	apperrors "github.com/runagent-robotics/run-factory/runfactory-core/internal/platform/errors"
)

type repoMock struct {
	createFn func(context.Context, domain.Factory) (domain.Factory, error)
}

func (r repoMock) Create(ctx context.Context, factory domain.Factory) (domain.Factory, error) {
	if r.createFn == nil {
		return domain.Factory{}, nil
	}
	return r.createFn(ctx, factory)
}

func (r repoMock) List(context.Context) ([]domain.Factory, error) { return nil, nil }
func (r repoMock) Get(context.Context, string) (domain.Factory, error) {
	return domain.Factory{}, nil
}
func (r repoMock) UpdateMap(context.Context, string, json.RawMessage) (domain.Factory, error) {
	return domain.Factory{}, nil
}
func (r repoMock) AddRobot(context.Context, string, domain.Robot) (domain.Factory, error) {
	return domain.Factory{}, nil
}
func (r repoMock) RemoveRobot(context.Context, string, string) (domain.Factory, error) {
	return domain.Factory{}, nil
}

func TestCreateFactoryValidation(t *testing.T) {
	svc := NewFactoryService(repoMock{})

	_, err := svc.CreateFactory(context.Background(), CreateFactoryInput{ID: "", Name: "x"})
	if apperrors.KindOf(err) != apperrors.KindValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateFactoryRepoConflict(t *testing.T) {
	svc := NewFactoryService(repoMock{
		createFn: func(context.Context, domain.Factory) (domain.Factory, error) {
			return domain.Factory{}, repository.ErrFactoryExists
		},
	})

	_, err := svc.CreateFactory(context.Background(), CreateFactoryInput{
		ID:    "factory-a",
		Name:  "Factory A",
		Map3D: json.RawMessage(`{"mesh":"a.glb"}`),
	})
	if !errors.Is(err, repository.ErrFactoryExists) {
		t.Fatalf("expected wrapped conflict error, got %v", err)
	}
	if apperrors.KindOf(err) != apperrors.KindConflict {
		t.Fatalf("expected conflict kind, got %s", apperrors.KindOf(err))
	}
}
