package handlers

import (
	"Scheduler-api/internal/models"
	"Scheduler-api/pkg/logger"
	"encoding/json"
	"net/http"
)

func RunSimulation(w http.ResponseWriter, r *http.Request) {
	var config models.SimulationConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		logger.Error("failed to decode simulation config", logger.ErrorField(err))
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info("starting simulation", logger.Int("ue_count", config.UECnt), logger.String("scheduler", config.BSScheduler))
}
