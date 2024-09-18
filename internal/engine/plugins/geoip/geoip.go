package axsi

import (
	"net"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/nzin/taxsi2/internal/engine"
	"github.com/oschwald/maxminddb-golang"
)

type MaxMindDbReader interface {
	Lookup(ip net.IP, result any) error
}

type GeoipWafPlugin struct {
	ds             db.DbService
	allowDenyTable map[string]bool
	allowMode      bool
	geoipdb        MaxMindDbReader
}

type GeoIP struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func NewGeoipWafPlugin(mmdbPath string, ds db.DbService) (engine.WafEnginePlugin, error) {
	mmdb, err := maxminddb.Open(mmdbPath)
	if err != nil {
		return nil, err
	}

	allowdeny, err := ds.GetGeoipCountries()
	if err != nil {
		return nil, err
	}
	allowMode := true

	// let's read one element to see if we are in an allow
	// or deny mode
	for _, v := range allowdeny {
		allowMode = v
		break
	}

	plugin := GeoipWafPlugin{
		ds:             ds,
		geoipdb:        mmdb,
		allowDenyTable: allowdeny,
		allowMode:      allowMode,
	}

	ds.SubscribeChanges(db.CHANGELOG_TABLE_GEOIP, &plugin)

	return &plugin, nil
}

func (g *GeoipWafPlugin) Name() string {
	return "geoip"
}

func (g *GeoipWafPlugin) Scan(payload *com.TaxsiCom) bool {
	// no geoip restriction?
	if len(g.allowDenyTable) == 0 {
		return true
	}

	ip := net.ParseIP(payload.RemoteAddr)
	var record GeoIP
	err := g.geoipdb.Lookup(ip, &record)
	if err != nil {
		return true
	}

	if g.allowMode {
		// allow only
		if _, found := g.allowDenyTable[record.Country.IsoCode]; !found {
			return false
		}
	} else {
		// deny
		if _, found := g.allowDenyTable[record.Country.IsoCode]; found {
			return false
		}
	}

	return true
}

func (g *GeoipWafPlugin) NotifyDbChange(key string) {
	allowdeny, err := g.ds.GetGeoipCountries()
	if err != nil {
		return
	}

	g.allowDenyTable = allowdeny
	allowMode := true

	// let's read one element to see if we are in an allow
	// or deny mode
	for _, v := range allowdeny {
		allowMode = v
		break
	}
	g.allowMode = allowMode
}
