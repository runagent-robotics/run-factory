package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
)

var (
	ErrFactoryExists   = errors.New("factory already exists")
	ErrFactoryNotFound = errors.New("factory not found")
	ErrRobotExists     = errors.New("robot already exists")
	ErrRobotNotFound   = errors.New("robot not found")
)

type FactoryRepository interface {
	Create(ctx context.Context, factory domain.Factory) (domain.Factory, error)
	List(ctx context.Context) ([]domain.Factory, error)
	Get(ctx context.Context, id string) (domain.Factory, error)
	UpdateMap(ctx context.Context, id string, map3D json.RawMessage) (domain.Factory, error)
	AddRobot(ctx context.Context, factoryID string, robot domain.Robot) (domain.Factory, error)
	RemoveRobot(ctx context.Context, factoryID, robotID string) (domain.Factory, error)
}
