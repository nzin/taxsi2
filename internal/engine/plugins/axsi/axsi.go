package axsi

import (
	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/nzin/taxsi2/internal/engine"
)

type AxiWafPlugin struct {
	ds db.DbService
}

func NewAxiWafPlugin(ds db.DbService) engine.WafEnginePlugin {

	return &AxiWafPlugin{
		ds: ds,
	}
}

func (a *AxiWafPlugin) Name() string {
	return "axi"
}

func (a *AxiWafPlugin) Scan(payload *com.TaxsiCom) bool {
	return true
}
