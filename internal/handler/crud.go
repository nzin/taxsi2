package handler

import (
	"bufio"
	"bytes"
	"io"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/swagger_gen/models"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/health"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/waf"

	"github.com/go-openapi/runtime/middleware"
)

// CRUD is the CRUD interface
type CRUD interface {
	// healthcheck
	GetHealthcheck(health.GetHealthParams) middleware.Responder
	PostSubmit(waf.PostSubmitParams) middleware.Responder
}

// NewCRUD creates a new CRUD instance
func NewCRUD() CRUD {
	return &crud{}
}

type crud struct{}

func (c *crud) GetHealthcheck(params health.GetHealthParams) middleware.Responder {
	return health.NewGetHealthOK().WithPayload(&models.Health{Status: "OK"})
}

func (c *crud) PostSubmit(params waf.PostSubmitParams) middleware.Responder {
	data, err := io.ReadAll(params.HTTPRequest.Body)
	if err != nil {
		return waf.NewPostSubmitDefault(503).WithPayload(
			ErrorMessage("unable to read the body content"),
		)
	}

	// try to unmarshall
	buf := bytes.NewBuffer(data)
	t, err := com.Unmarshall(bufio.NewReader(buf))
	if err != nil {
		return waf.NewPostSubmitDefault(503).WithPayload(
			ErrorMessage("unable to unmarshall payload: %v", err),
		)
	}

	// test
	if t.Method == "PATCH" {
		return &waf.PostSubmitForbidden{}
	}

	return &waf.PostSubmitOK{}
}
