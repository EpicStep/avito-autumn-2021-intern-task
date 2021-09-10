package balance

import "github.com/go-chi/chi/v5"

// Routes add new routes to chi Router.
func (s *Service) Routes(r chi.Router) {
	r.Route("/balance", func(r chi.Router) {
		r.Get("/", s.GetBalance)
		r.Post("/", s.ControlBalance)

		r.Get("/history", s.TransactionsHistory)

		r.Post("/transfer", s.Transfer)
	})
}
