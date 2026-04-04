package handlers

import (
	"Scheduler-api/internal/models"
	"Scheduler-api/internal/service"
	"Scheduler-api/internal/validator"
	"Scheduler-api/pkg/logger"
	"encoding/json"
	"net/http"
)

var SimService = service.NewSimulationService()

func RunSimulation(w http.ResponseWriter, r *http.Request) {
	var config models.SimulationConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		logger.Error("failed to decode simulation config", logger.ErrorField(err))
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("starting simulation", logger.Int("ue_count", config.UECnt), logger.String("scheduler", config.BSScheduler))

	if err := validator.ValidateConfig(&config); err != nil {
		logger.Error("validation failed for simulation config", logger.ErrorField(err))
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(config.UEIDs) != config.UECnt {
		http.Error(w, "Validation error: len(ue_ids) must equal ue_cnt", http.StatusBadRequest)
		return
	}

	result, err := SimService.RunSimulation(&config)
	if err != nil {
		logger.Error("failed to start simulation", logger.ErrorField(err))
		http.Error(w, "Failed to start simulation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("simulation started", logger.String("simulation_id", result.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"code":   http.StatusCreated,
		"data":   result,
	})
}
