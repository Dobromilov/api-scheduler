package api

import (
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
	// r.Route("api/v1", func(r chi.Router) {})
	return r
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
