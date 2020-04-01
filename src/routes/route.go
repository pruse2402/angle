package route

import (
	"angle/src/handlers"

	"github.com/go-chi/chi"
)

func NewRouter(h *handlers.Provider) *chi.Mux {

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		// Vendor routes
		r.Get("/vendors/list", h.GetVendors)                   // GET Vendors Details
		r.Post("/vendor/new", h.InsertVendorDetails)           // Save Vendor Details
		r.Put("/vendor/{id}/edit", h.UpdateVendorDetails)      // Update Vendor Details
		r.Delete("/vendor/{id}/delete", h.RemoveVendorDetails) // Remove Vendor Details

		// Customer routes
		r.Get("/customers/list", h.GetCustomers)                   // GET Customers Details
		r.Post("/customer/new", h.InsertCustomerDetails)           // Save Customer Details
		r.Put("/customer/{id}/edit", h.UpdateCustomerDetails)      // Update Customer Details
		r.Delete("/customer/{id}/delete", h.RemoveCustomerDetails) // Remove Customer Details

		r.Get("/ping", h.Ping)

	})

	return r
}
