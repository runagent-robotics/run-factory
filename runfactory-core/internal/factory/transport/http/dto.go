package httptransport

import (
	"encoding/json"
	"slices"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
)

type createFactoryRequest struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Map3D json.RawMessage `json:"map3d"`
}

type updateMapRequest struct {
	Map3D json.RawMessage `json:"map3d"`
}

type addRobotRequest struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Position *domain.Position3D `json:"position"`
}

type factoryResponse struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Map3D  json.RawMessage `json:"map3d,omitempty"`
	Robots []robotResponse `json:"robots"`
}

type robotResponse struct {
	ID       string             `json:"id"`
	Name     string             `json:"name,omitempty"`
	Position *domain.Position3D `json:"position,omitempty"`
}

func toFactoryResponse(factory domain.Factory) factoryResponse {
	robots := make([]robotResponse, 0, len(factory.Robots))
	for _, robot := range factory.Robots {
		robots = append(robots, robotResponse{
			ID:       robot.ID,
			Name:     robot.Name,
			Position: domain.ClonePosition(robot.Position),
		})
	}

	slices.SortFunc(robots, func(a, b robotResponse) int {
		switch {
		case a.ID < b.ID:
			return -1
		case a.ID > b.ID:
			return 1
		default:
			return 0
		}
	})

	return factoryResponse{
		ID:     factory.ID,
		Name:   factory.Name,
		Map3D:  domain.CloneRawMessage(factory.Map3D),
		Robots: robots,
	}
}
