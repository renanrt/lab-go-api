package api

import (
	"context"
	"net/http"

	"github.com/renanrt/lab-go-api/service"
)

// serviceHandler is a middleware http.Handler that adds a service.TaxService to the context.
type serviceHandler struct {
	service service.TaxService
	inner   http.Handler
}

func newServiceHandler(service service.TaxService, inner http.Handler) http.Handler {
	return &serviceHandler{service, inner}
}

type serviceKeyType int

const serviceKey serviceKeyType = iota

func (h *serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), serviceKey, h.service)
	h.inner.ServeHTTP(w, r.WithContext(ctx))
}

// getService retrieves a service.TaxService from the context. The request that's passed in must
// have gone through dbHandler.
func getService(r *http.Request) service.TaxService {
	return service.TaxService{}
}
