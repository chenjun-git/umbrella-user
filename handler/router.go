package handler

import (
	"fmt"

	"github.com/go-chi/chi"

	"business/user/handler/v1.0"
)

func RegisterUserRouter() *chi.Mux {
	router := chi.NewRouter()
	v1_0.registerRouter(router)
	return router
}
