package service

import (
	"Scheduler-api/internal/models"
	"Scheduler-api/internal/python"
	"Scheduler-api/internal/store"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SimulationService struct {
	store       *store.Store
	pythonRunner *python.Runner
}

func NewSimulationService() *SimulationService {
	return &SimulationService{
		store:       store.GetStore(),
		pythonRunner: python.NewRunner(),
	}
}

type RunSimulationResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type SimulationStatus struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	CurrentTTI int       `json:"current_tti"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

func (s *SimulationService) RunSimulation(config *models.SimulationConfig) (*RunSimulationResult, error) {
	id := uuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())

	sim := &store.Simulation{
		ID:         id,
		Config:     config,
		Status:     "running",
		StartTime:  time.Now(),
		CurrentTTI: 0,
		CancelFunc: cancel,
		Ctx:        ctx,
	}

	s.store.Add(sim)

	go func() {
		err := s.pythonRunner.RunSimulation(ctx, id, config)
		if err != nil {
			s.store.UpdateStatus(id, "error")
			return
		}
		s.store.UpdateStatus(id, "completed")
	}()

	return &RunSimulationResult{
		ID:     id,
		Status: "running",
	}, nil
}

func (s *SimulationService) StopSimulation(id string) error {
	sim, ok := s.store.Get(id)
	if !ok {
		return fmt.Errorf("simulation %s not found", id)
	}

	sim.CancelFunc()
	s.store.UpdateStatus(id, "stopped")
	return nil
}

func (s *SimulationService) PauseSimulation(id string) error {
	sim, ok := s.store.Get(id)
	if !ok {
		return fmt.Errorf("simulation %s not found", id)
	}

	if sim.Status != "running" {
		return fmt.Errorf("simulation %s is not running (current status: %s)", id, sim.Status)
	}

	s.store.UpdateStatus(id, "paused")
	return nil
}

func (s *SimulationService) ResumeSimulation(id string) error {
	sim, ok := s.store.Get(id)
	if !ok {
		return fmt.Errorf("simulation %s not found", id)
	}

	if sim.Status != "paused" {
		return fmt.Errorf("simulation %s is not paused (current status: %s)", id, sim.Status)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sim.Ctx = ctx
	sim.CancelFunc = cancel
	sim.Status = "running"

	go func() {
		err := s.pythonRunner.RunSimulation(ctx, id, sim.Config)
		if err != nil {
			s.store.UpdateStatus(id, "error")
			return
		}
		s.store.UpdateStatus(id, "completed")
	}()

	return nil
}

func (s *SimulationService) GetSimulationStatus(id string) (*SimulationStatus, error) {
	sim, ok := s.store.Get(id)
	if !ok {
		return nil, fmt.Errorf("simulation %s not found", id)
	}

	return &SimulationStatus{
		ID:         sim.ID,
		Status:     sim.Status,
		CurrentTTI: sim.CurrentTTI,
		StartTime:  sim.StartTime,
		EndTime:    sim.EndTime,
	}, nil
}

func (s *SimulationService) ListSimulations() []*SimulationStatus {
	sims := s.store.List()
	result := make([]*SimulationStatus, 0, len(sims))

	for _, sim := range sims {
		result = append(result, &SimulationStatus{
			ID:         sim.ID,
			Status:     sim.Status,
			CurrentTTI: sim.CurrentTTI,
			StartTime:  sim.StartTime,
			EndTime:    sim.EndTime,
		})
	}

	return result
}
