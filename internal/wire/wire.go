package wire

import (
	"aplikasi-pos-team-boolean/internal/adaptor"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func Wiring(repo repository.Repository) *chi.Mux {
	router := chi.NewRouter()
	rV1 := chi.NewRouter()
	wireStaff(rV1, repo)
	router.Mount("/api/v1", rV1)

	return router
}

func wireStaff(router *chi.Mux, repo repository.Repository) {
	usecaseStaff := usecase.NewStaffUseCase(repo)
	adaptorStaff := adaptor.NewStaffAdaptor(usecaseStaff)
	router.Get("/staff", adaptorStaff.GetList)
	router.Post("/staff", adaptorStaff.Create)
	router.Get("/staff/{id}", adaptorStaff.GetByID)
	router.Put("/staff/{id}", adaptorStaff.Update)
	router.Delete("/staff/{id}", adaptorStaff.Delete)
}
