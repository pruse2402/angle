package route

import (
	"angle/src/handlers"

	"github.com/go-chi/chi"
)

func NewRouter(h *handlers.Provider) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Get("/vendors", h.GetVendors)              // GET Vendor Details
		r.Post("/vendor/new", h.InsertVendorDetails) // Save Vendor Details
		r.Put("/vendor/{id}", h.UpdateVendorDetails) // Update Vendor Details
		r.Get("/ping", h.Ping)

	})

	return r
}
