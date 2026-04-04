package store

import (
	"Scheduler-api/internal/models"
	"context"
	"sync"
	"time"
)

type Simulation struct {
	ID         string
	Config     *models.SimulationConfig
	Status     string
	StartTime  time.Time
	EndTime    time.Time
	CurrentTTI int
	CancelFunc context.CancelFunc
	Ctx        context.Context
}

type Store struct {
	simulation map[string]*Simulation
	mu         sync.RWMutex
}

var (
	instance *Store
	once     sync.Once
)

func GetStore() *Store {
	once.Do(func() {
		instance = &Store{
			simulation: make(map[string]*Simulation),
		}
	})
	return instance
}

func (s *Store) Add(sim *Simulation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.simulation[sim.ID] = sim
}

func (s *Store) Get(id string) (*Simulation, bool) {
	s.mu.RLock()
	defer s.mu.Unlock()
	sim, ok := s.simulation[id]
	return sim, ok
}

func (s *Store) List() []*Simulation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Simulation, 0, len(s.simulation))
	for _, sim := range s.simulation {
		list = append(list, sim)
	}
	return list
}

func (s *Store) UpdateStatus(id string, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sim, ok := s.simulation[id]; ok {
		sim.Status = status
		if status == "stopped" || status == "completed" {
			sim.EndTime = time.Now()
		}
	}
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.simulation, id)
}
