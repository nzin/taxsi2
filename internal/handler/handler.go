package handler

import (
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/health"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/waf"
)

// Setup initialize all the handler functions
func Setup(api *operations.Taxsi2API) {
	c := NewCRUD()

	// healthcheck
	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(c.GetHealthcheck)
	api.WafPostSubmitHandler = waf.PostSubmitHandlerFunc(c.PostSubmit)
}
