package python

import (
	"Scheduler-api/internal/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Runner struct {
	mu        sync.RWMutex
	processes map[string]*exec.Cmd
}

func NewRunner() *Runner {
	return &Runner{
		processes: make(map[string]*exec.Cmd),
	}
}

func (r *Runner) RunSimulation(ctx context.Context, simID string, config *models.SimulationConfig) error {
	r.mu.Lock()

	scriptPath := findPythonScript()
	if scriptPath == "" {
		r.mu.Unlock()
		return fmt.Errorf("python simulation script not found")
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	cmd := exec.CommandContext(ctx, "python3", scriptPath, "--config", string(configJSON))
	cmd.Dir = filepath.Dir(scriptPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to start python process: %w", err)
	}

	r.processes[simID] = cmd
	r.mu.Unlock()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			_ = scanner.Text()
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.Canceled {
			return nil
		}
		return fmt.Errorf("python process exited with error: %w", err)
	}

	r.mu.Lock()
	delete(r.processes, simID)
	r.mu.Unlock()

	return nil
}

func (r *Runner) StopSimulation(simID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cmd, ok := r.processes[simID]
	if !ok {
		return fmt.Errorf("simulation %s process not found", simID)
	}

	return cmd.Process.Kill()
}

func (r *Runner) IsRunning(simID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.processes[simID]
	return ok
}

func findPythonScript() string {
	candidates := []string{
		"../../PyScheduler/SIMULATION_MANAGER.py",
		"../PyScheduler/SIMULATION_MANAGER.py",
		"PyScheduler/SIMULATION_MANAGER.py",
	}

	for _, path := range candidates {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	return ""
}
