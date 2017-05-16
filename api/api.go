package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/vend/go-common/api/fail"
)

func Serve() error {
	app, err := NewServer()
	if err != nil {
		return err
	}

	return http.ListenAndServe(":"+"8080", app)
}

type HelloWorldHandler struct {
	http.Handler
}

func (h *HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}

// Serve serves the API. It only returns if there is an error.
func NewServer() (*HelloWorldHandler, error) {
	r := httprouter.New()
	r.GET("/api/2.0/taxes-groups/search", searchTaxes)
	r.GET("/healthcheck", healthCheck)
	return &HelloWorldHandler{r}, nil
}
func healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func searchTaxes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	queryValues := r.URL.Query()

	//retailer := squire.RequestRetailer(r)
	country := queryValues.Get("country")
	state := queryValues.Get("state")
	street := queryValues.Get("street")
	city := queryValues.Get("city")
	zipcode := queryValues.Get("zipcode")
	provider := queryValues.Get("provider")

	if zipcode == "" {
		RespondWithError(w, r, fail.NewAPIResponseError(errors.New("No zipcode informed"), http.StatusBadRequest))
		return
	}

	service := getService(r)
	retailerID := "dummy-retailer-id"
	obj, err := service.GetTaxesForAddress(provider, retailerID, country, state, city, zipcode, street)

	if err != nil {
		RespondWithError(w, r, err)
		return
	}

	RespondWithData(w, r, obj, http.StatusOK)
}
