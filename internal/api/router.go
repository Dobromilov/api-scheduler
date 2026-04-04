package api

import (
	"Scheduler-api/internal/api/handlers"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", healthCheckHandler)

	r.Get("/api/v1/scheduler/algorithms", listSchedulersHandler)
	r.Get("/api/v1/mobility/models", listMobilityModelsHandler)
	r.Get("/api/v1/traffic/models", listTrafficModelsHandler)

	r.Post("/api/v1/simulation/run", handlers.RunSimulation)
	r.Post("/api/v1/simulation/{id}/stop", stopSimulationHandler)
	r.Post("/api/v1/simulation/{id}/pause", pauseSimulationHandler)
	r.Post("/api/v1/simulation/{id}/resume", resumeSimulationHandler)
	r.Get("/api/v1/simulation/{id}", getSimulationStatusHandler)
	r.Get("/api/v1/simulations", listSimulationsHandler)

	r.Get("/ws/simulation/{id}", handlers.HandleSimulationWS)

	return r
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func listSchedulersHandler(w http.ResponseWriter, r *http.Request) {
	algorithms := map[string]interface{}{
		"algorithms": []string{
			"BestCQI",
			"RoundRobin",
			"ProportionalFair",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(algorithms)
}

func listMobilityModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := map[string]interface{}{
		"mobility_models": []string{
			"RandomWalk",
			"Static",
			"ConstantVelocity",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}

func listTrafficModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := map[string]interface{}{
		"traffic_models": []string{
			"Poisson",
			"CBR",
			"FTP",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}

func stopSimulationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"simulation ID is required"}`, http.StatusBadRequest)
		return
	}

	if err := handlers.SimService.StopSimulation(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"stopped","id":"` + id + `"}`))
}

func pauseSimulationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"simulation ID is required"}`, http.StatusBadRequest)
		return
	}

	if err := handlers.SimService.PauseSimulation(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"paused","id":"` + id + `"}`))
}

func resumeSimulationHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"simulation ID is required"}`, http.StatusBadRequest)
		return
	}

	if err := handlers.SimService.ResumeSimulation(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"resumed","id":"` + id + `"}`))
}

func getSimulationStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"simulation ID is required"}`, http.StatusBadRequest)
		return
	}

	status, err := handlers.SimService.GetSimulationStatus(id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func listSimulationsHandler(w http.ResponseWriter, r *http.Request) {
	simulations := handlers.SimService.ListSimulations()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(simulations)
}
