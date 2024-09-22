package handler

import (
	"bufio"
	"bytes"
	"io"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/config"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/nzin/taxsi2/internal/engine"
	"github.com/nzin/taxsi2/internal/engine/plugins/axsi"
	"github.com/nzin/taxsi2/internal/engine/plugins/geoip"
	"github.com/nzin/taxsi2/swagger_gen/models"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/health"
	"github.com/nzin/taxsi2/swagger_gen/restapi/operations/waf"
	"github.com/sirupsen/logrus"

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
	ds, err := db.NewDbService(
		config.Config.DBDriver,
		config.Config.DBConnectionStr,
		config.Config.DBConnectionRetryAttempts,
		config.Config.DBConnectionRetryDelay,
	)
	if err != nil {
		panic(err)
	}

	e, err := engine.NewWafEngineImpl(
		ds,
		config.Config.WafOutput,
		config.Config.WafOutputFormat,
	)
	if err != nil {
		panic(err)
	}

	// register scan plugins
	geoipPlugin, err := geoip.NewGeoipWafPlugin(config.Config.DefaultGeoipDbPath, config.Config.DefaultRemoteGeoipDbPath, ds)
	if err != nil {
		logrus.Errorf("unable to create geoip plugin: %v", err)
	} else {
		e.RegisterPlugin(geoipPlugin)
	}

	axsi := axsi.NewAxiWafPlugin(ds)
	e.RegisterPlugin(axsi)

	// for later
	// - botmanager?
	// - ratelimiter?

	return &crud{
		ds:        ds,
		wafEngine: e,
	}
}

type crud struct {
	ds        db.DbService
	wafEngine engine.WafEngine
}

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

	res := c.wafEngine.Scan(t)
	if !res {
		return &waf.PostSubmitForbidden{}
	}

	return &waf.PostSubmitOK{}
}
