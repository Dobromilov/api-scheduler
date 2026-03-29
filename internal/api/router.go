package api

import (
	"Scheduler-api/internal/api/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) //Присваивает каждому входящему HTTP-запросу уникальный ID (UUID).
	r.Use(middleware.RealIP)    //Определяет настоящий IP-адрес клиента
	r.Use(middleware.Recoverer) //Чтобы сервер не падал при панике в хендлере

	r.Get("/health", healthCheckHandler)
	r.Route("api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/scheduler/algorithms", listSchedulersHandler)
			r.Get("/mobility/models", listMobilityModelsHandler)
			r.Get("/traffic/models", listTrafficModelsHandler)
		})
		r.Route("simlation", func(r chi.Router) {
			r.Post("/run", handlers.RunSimulation)
		})
	})
	r.Get("/ws/simulation/{id}", handleSimulationWS)

	return r
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func listSchedulersHandler(w http.ResponseWriter, r *http.Request)     { /* TODO: return JSON */ }
func listMobilityModelsHandler(w http.ResponseWriter, r *http.Request) { /* TODO: return JSON */ }
func listTrafficModelsHandler(w http.ResponseWriter, r *http.Request)  { /* TODO: return JSON */ }
func runSimulationHandler(w http.ResponseWriter, r *http.Request) { /* TODO: parse config, start python, return 201 and ID */
}
func stopSimulationHandler(w http.ResponseWriter, r *http.Request) { /* TODO: kill process */ }
func pauseSimulationHandler(w http.ResponseWriter, r *http.Request) { /* TODO: send signal to python */
}
func resumeSimulationHandler(w http.ResponseWriter, r *http.Request) { /* TODO: send signal to python */
}
func getSimulationStatusHandler(w http.ResponseWriter, r *http.Request) { /* TODO: query store */ }
func handleSimulationWS(w http.ResponseWriter, r *http.Request)         { /* TODO: upgrade to WS */ }
