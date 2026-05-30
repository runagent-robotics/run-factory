package memory

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/domain"
	"github.com/runagent-robotics/run-factory/runfactory-core/internal/factory/repository"
)

func TestFactoryStoreLifecycle(t *testing.T) {
	store := NewFactoryStore()
	ctx := context.Background()

	created, err := store.Create(ctx, domain.Factory{
		ID:    "factory-a",
		Name:  "Factory A",
		Map3D: json.RawMessage(`{"mesh":"factory-a.glb"}`),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := store.AddRobot(ctx, created.ID, domain.Robot{ID: "robot-1", Name: "Robot 1"}); err != nil {
		t.Fatalf("add robot failed: %v", err)
	}

	got, err := store.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if len(got.Robots) != 1 {
		t.Fatalf("expected 1 robot, got %d", len(got.Robots))
	}

	if _, err := store.RemoveRobot(ctx, created.ID, "robot-1"); err != nil {
		t.Fatalf("remove robot failed: %v", err)
	}
}

func TestFactoryStoreNotFound(t *testing.T) {
	store := NewFactoryStore()
	_, err := store.Get(context.Background(), "missing")
	if !errors.Is(err, repository.ErrFactoryNotFound) {
		t.Fatalf("expected ErrFactoryNotFound, got %v", err)
	}
}
