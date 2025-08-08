package main

import (
	"errors"
	"sync"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int)
	GetTruck(id string) (Truck, error)
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, cargo int) error
	IncrementTruckCargo(id string, increment int) error
}

type Truck struct {
	ID    string
	Cargo int
}

type truckManager struct {
	trucks map[string]*Truck
	sync.RWMutex
}

func NewTruckManager() FleetManager {
	return &truckManager{
		trucks: make(map[string]*Truck),
	}
}

func (m *truckManager) AddTruck(id string, cargo int) {
	m.Lock()
	defer m.Unlock()

	m.trucks[id] = &Truck{
		ID:    id,
		Cargo: cargo,
	}
}

func (m *truckManager) GetTruck(id string) (Truck, error) {
	m.RLock()
	defer m.RUnlock()

	truck, exists := m.trucks[id]
	if !exists {
		return Truck{}, ErrTruckNotFound
	}
	// Return a copy to prevent external modification
	return *truck, nil
}

func (m *truckManager) RemoveTruck(id string) error {
	m.Lock()
	defer m.Unlock()

	_, exists := m.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}
	delete(m.trucks, id)
	return nil
}

func (m *truckManager) UpdateTruckCargo(id string, cargo int) error {
	m.Lock()
	defer m.Unlock()

	truck, exists := m.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}
	// m.trucks[id].Cargo = cargo
	truck.Cargo = cargo
	return nil
}

func (m *truckManager) IncrementTruckCargo(id string, increment int) error {
	m.Lock()
	defer m.Unlock()

	truck, exists := m.trucks[id]
	if !exists {
		return ErrTruckNotFound
	}
	truck.Cargo += increment
	return nil
}
