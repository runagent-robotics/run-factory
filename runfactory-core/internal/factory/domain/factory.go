package domain

import "encoding/json"

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

type Factory struct {
	ID     string
	Name   string
	Map3D  json.RawMessage
	Robots map[string]Robot
}

func CloneRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	out := make(json.RawMessage, len(raw))
	copy(out, raw)
	return out
}

func ClonePosition(position *Position3D) *Position3D {
	if position == nil {
		return nil
	}
	copyPosition := *position
	return &copyPosition
}

func CloneRobot(robot Robot) Robot {
	return Robot{
		ID:       robot.ID,
		Name:     robot.Name,
		Position: ClonePosition(robot.Position),
	}
}

func CloneFactory(factory Factory) Factory {
	copyFactory := Factory{
		ID:     factory.ID,
		Name:   factory.Name,
		Map3D:  CloneRawMessage(factory.Map3D),
		Robots: make(map[string]Robot, len(factory.Robots)),
	}
	for id, robot := range factory.Robots {
		copyFactory.Robots[id] = CloneRobot(robot)
	}
	return copyFactory
}
