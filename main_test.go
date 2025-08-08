package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTruck(t *testing.T) {
	manager := NewTruckManager()
	manager.AddTruck("1", 100)

	tx := manager.(*truckManager)

	if len(tx.trucks) != 1 {
		t.Errorf("Expected 1 truck, got %d", len(tx.trucks))
	}
}

func TestGetTruck(t *testing.T) {
	manager := NewTruckManager()
	manager.AddTruck("1", 100)

	truck, err := manager.GetTruck("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if truck.ID != "1" {
		t.Errorf("Expected truck ID to be 1, got %s", truck.ID)
	}
}

func TestRemoveTruck(t *testing.T) {
	manager := NewTruckManager()
	manager.AddTruck("1", 100)

	manager.RemoveTruck("1")

	tx := manager.(*truckManager)
	_, err := tx.GetTruck("1")
	if err != ErrTruckNotFound {
		t.Errorf("Expected truck not found error, got %v", err)
	}

	if len(tx.trucks) != 0 {
		t.Errorf("Expected 0 trucks, got %d", len(tx.trucks))
	}
}

func TestUpdateTruckCargo(t *testing.T) {
	manager := NewTruckManager()
	manager.AddTruck("1", 100)

	manager.UpdateTruckCargo("1", 200)

	truck, err := manager.GetTruck("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if truck.Cargo != 200 {
		t.Errorf("Expected truck cargo to be 200, got %d", truck.Cargo)
	}
}

func TestConcurrentUpdate(t *testing.T) {
	manager := NewTruckManager()
	manager.AddTruck("1", 100)
	const numGoroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				manager.IncrementTruckCargo("1", 1)
			}
		}()
	}

	wg.Wait()
	truck, _ := manager.GetTruck("1")
	require.Equal(t, numGoroutines*iterations+100, truck.Cargo, "Final cargo should be equal to the number of goroutines times iterations plus initial cargo")
}
