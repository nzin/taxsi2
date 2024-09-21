package geoip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/nzin/taxsi2/internal/engine"
	"github.com/oschwald/maxminddb-golang"
	"github.com/sirupsen/logrus"
)

type MaxMindDbReader interface {
	Lookup(ip net.IP, result any) error
}

type GeoipWafPlugin struct {
	ds             db.DBServiceGeoip
	allowDenyTable map[string]bool
	allowMode      bool
	geoipdb        MaxMindDbReader
}

type GeoIP struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func NewGeoipWafPlugin(mmdbPath string, remoteMmdPath string, ds db.DBServiceGeoip) (engine.WafEnginePlugin, error) {

	// let's init the maxmind db
	mmdb, err := initMaxmindDb(ds, mmdbPath, remoteMmdPath)
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

	// if the database changes, we need to reload the allow/deny table
	ds.SubscribeChanges(db.CHANGELOG_TABLE_GEOIP, &plugin)

	return &plugin, nil
}

func initMaxmindDb(ds db.DBServiceGeoip, mmdbPath string, remoteMmdPath string) (MaxMindDbReader, error) {
	// let's check if there is a mmdb file in the database
	// and if it is recent
	// else let's try to download it
	// else let's check on disk

	timestamp, err := ds.GetGeoipSourceTimestamp()
	if err == nil {
		// check if the file is recent
		if time.Since(timestamp) < 31*24*time.Hour {
			// file is recent
			src, err := ds.GetGeoipSource()
			if err == nil {
				mmdb, err := maxminddb.FromBytes(src.Source)
				if err == nil {
					return mmdb, nil
				}
			}
		}
	}

	// let's try to download the file
	// https://download.db-ip.com/free/dbip-country-lite-<year>-<month>.mmdb.gz
	now := time.Now()
	content, err := downloadGzFileToMemory(fmt.Sprintf(remoteMmdPath, now.Year(), int(now.Month())))
	if err == nil {
		mmdb, err := maxminddb.FromBytes(content)
		if err == nil {
			// let's store the file in the database
			err = ds.SetGeoipSource(now, content)
			if err != nil {
				logrus.Errorf("error storing fresly downloaded mmdb file in the database: %v", err)
				return nil, err
			}
			return mmdb, nil
		}
	}
	logrus.Infof("unable to download mmdb file (%s): %v", remoteMmdPath, err)

	mmdb, err := maxminddb.Open(mmdbPath)
	if err != nil {
		return nil, fmt.Errorf("issue opening file %s: %v", mmdbPath, err)
	}

	// the file will always be obsolete
	err = ds.SetGeoipSource(now.Add(-30*24*time.Hour), content)
	if err != nil {
		logrus.Errorf("error storing mmdb file in the database: %v", err)
	}

	return mmdb, nil
}

func downloadGzFileToMemory(url string) ([]byte, error) {
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check if the server response was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Read the response body into a byte slice
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return gunzipData(buf.Bytes())
}

func gunzipData(compressedData []byte) ([]byte, error) {
	// Create a bytes buffer from the compressed data
	buf := bytes.NewBuffer(compressedData)

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("unable to create a gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Read the decompressed data into a buffer
	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, gzipReader)
	if err != nil {
		return nil, fmt.Errorf("uname to get the uncompress stream: %v", err)
	}

	// Return the decompressed data as a byte slice
	return decompressedData.Bytes(), nil
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

// NotifyDbChange is called when the database changes
// we need to reload the allow/deny table
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
