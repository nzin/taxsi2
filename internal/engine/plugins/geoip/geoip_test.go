package geoip

import (
	"compress/gzip"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/nzin/taxsi2/internal/com"
	"github.com/nzin/taxsi2/internal/db"
	"github.com/oschwald/maxminddb-golang"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	t.Run("Download happy path", func(t *testing.T) {
		// let's start a fake http server
		// that will serve a gzipped file
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)

				gz := gzip.NewWriter(w)
				defer gz.Close()
				_, err := gz.Write([]byte("test"))
				if err != nil {
					http.Error(w, "Failed to write gzip data", http.StatusInternalServerError)
					return
				}
			}),
		)
		defer server.Close()

		// let's download the file
		data, err := downloadGzFileToMemory(server.URL)
		assert.Nil(t, err)
		assert.Equal(t, []byte("test"), data)
	})
}

const base64MmdContent = "AAABAAA0AAACAAAoAAADAACkAAAEAACkAAAFAAAlAAAGAACkAAAHAACkAACkAAAIAAAJAACkAAAKAACkAAALAACkAAAMAACkAAANAACkAAAOAACkAAAPAACkAACkAAAQAAARAACkAAASAACkAAATAACkAAAUAACkAAAVAACkAAAWAACkAAAXAACkAACkAAAYAAAZAACkAAAaAACkAAAbAAAgAAAcAAC0AAAdAADBAAAeAADNAAAfAADZAACkAADlAAAhAACkAAAiAACkAAAjAACkAAAkAACkAADxAACkAAAmAACkAACkAAAnAACkAACkAACkAAApAAAqAAAwAAArAACkAACkAAAsAAAtAACkAAAuAACkAAAvAACkAACkAACkAACkAAAxAACkAAAyAACkAAAzAACkAACkAAA1AABJAACkAAA2AAA3AACkAACkAAA4AAA5AABDAAA6AACkAACkAAA7AACkAAA8AACkAAA9AACkAAA+AACkAAA/AACkAABAAACkAABBAACkAABCAACkAACkAABEAACkAABFAACkAABGAACkAABHAACkAABIAACkAACkAACkAABKAACjAABLAACkAABMAACQAABNAAB6AABOAACkAABPAACkAABQAABzAABRAABlAABSAACkAABTAACkAABUAACkAABVAACkAABWAACkAABXAACkAABYAACkAABZAACkAABaAACkAABbAACkAABcAACkAABdAACkAABeAABkAABfAACkAABgAACkAABhAACkAABiAACkAABjAACkAACkAACkAACkAACkAABmAACkAACkAABnAACkAABoAABpAACkAABqAACkAABrAACkAABsAACkAACkAABtAACkAABuAABvAACkAABwAACkAABxAACkAACkAAByAACkAACkAAB0AACkAACkAAB1AAB2AACkAACkAAB3AAB4AACkAAB5AACkAACkAACkAACkAAB7AAB8AACkAAB9AACkAAB+AACkAAB/AACDAACkAACAAACBAACkAACCAACkAACkAACkAACkAACEAACFAACkAACGAACkAACkAACHAACkAACIAACJAACkAACkAACKAACkAACLAACMAACkAACNAACkAACkAACOAACPAACkAACkAACkAACRAACkAACkAACSAACkAACTAACUAACkAACVAACkAACWAACkAACXAACkAACYAACkAACZAACkAACaAACkAACbAACkAACcAACkAACkAACdAACkAACeAACkAACfAACgAACkAAChAACkAACiAACkAACkAACkAACkAACkAAAAAAAAAAAAAAAAAAAAAOFCaXBIMS4xLjEuMTbhQmlwRzEuMS4xLjjhQmlwRzEuMS4xLjThQmlwRzEuMS4xLjLhQmlwRzEuMS4xLjHhQmlwSDEuMS4xLjMyq83vTWF4TWluZC5jb23pW2JpbmFyeV9mb3JtYXRfbWFqb3JfdmVyc2lvbqECW2JpbmFyeV9mb3JtYXRfbWlub3JfdmVyc2lvbqBLYnVpbGRfZXBvY2gEAmLf/9ZNZGF0YWJhc2VfdHlwZURUZXN0S2Rlc2NyaXB0aW9u4kJlbk1UZXN0IERhdGFiYXNlQnpoVVRlc3QgRGF0YWJhc2UgQ2hpbmVzZUppcF92ZXJzaW9uoQRJbGFuZ3VhZ2VzAgRCZW5CemhKbm9kZV9jb3VudMGkS3JlY29yZF9zaXploRg="

/*
 * This is a mock implementation of the db.DBServiceGeoip interface
 */
type DbServiceGeoipMock struct {
	timestamp time.Time
	content   []byte
	countries map[string]bool
}

func (dsgm *DbServiceGeoipMock) GetGeoipSourceTimestamp() (time.Time, error) {
	return dsgm.timestamp, nil
}

func (dsgm *DbServiceGeoipMock) GetGeoipSource() (*db.GeoipSource, error) {
	return &db.GeoipSource{Timestamp: dsgm.timestamp, Source: dsgm.content}, nil
}

func (dsgm *DbServiceGeoipMock) SetGeoipSource(timestamp time.Time, source []byte) error {
	return nil
}

func (dsgm *DbServiceGeoipMock) SubscribeChanges(table int, listener db.DbChangeListener) {}

func (dsgm *DbServiceGeoipMock) GetGeoipCountries() (map[string]bool, error) {
	return dsgm.countries, nil
}

func (dsgm *DbServiceGeoipMock) SetGeoipCountries(countryCodes []string, allow bool) error {
	return nil
}

func TestInitMmdb(t *testing.T) {
	t.Run("InitMmdb happy path: db is uptodate", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		downloaded := false
		// let's start a fake http server
		// that will serve a gzipped file
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)

				gz := gzip.NewWriter(w)
				defer gz.Close()

				_, err = gz.Write(data)
				if err != nil {
					http.Error(w, "Failed to write gzip data", http.StatusInternalServerError)
					return
				}
				downloaded = true
			}),
		)
		defer server.Close()

		// let's get from the DB, to the http server to the file
		ds := &DbServiceGeoipMock{
			timestamp: time.Now(),
			content:   data,
			countries: make(map[string]bool),
		}
		mmdb, err := initMaxmindDb(ds, "/tmp/foobar.mmdb", server.URL+"/foobar-%d-%d.mmdb.gz")
		assert.Nil(t, err)
		assert.NotNil(t, mmdb)
		assert.False(t, downloaded)
	})

	t.Run("InitMmdb happy path: db is not uptodate, but http server is here", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		downloaded := false
		// let's start a fake http server
		// that will serve a gzipped file
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)

				gz := gzip.NewWriter(w)
				defer gz.Close()

				_, err = gz.Write(data)
				if err != nil {
					http.Error(w, "Failed to write gzip data", http.StatusInternalServerError)
					return
				}
				downloaded = true
			}),
		)
		defer server.Close()

		// let's get from the DB, to the http server to the file
		ds := &DbServiceGeoipMock{
			timestamp: time.Now().Add(-32 * 24 * time.Hour),
			content:   data,
			countries: make(map[string]bool),
		}
		mmdb, err := initMaxmindDb(ds, "/tmp/foobar.mmdb", server.URL+"/foobar-%d-%d.mmdb.gz")
		assert.Nil(t, err)
		assert.NotNil(t, mmdb)
		assert.True(t, downloaded)
	})

	t.Run("InitMmdb happy path: db is not uptodate, http server doesn't return correct content, but we have a tmp file", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		downloaded := false
		// let's start a fake http server
		// that will serve a gzipped file
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)

				gz := gzip.NewWriter(w)
				defer gz.Close()

				_, err = gz.Write([]byte("foobar"))
				if err != nil {
					http.Error(w, "Failed to write gzip data", http.StatusInternalServerError)
					return
				}
				downloaded = true
			}),
		)
		defer server.Close()

		// let's get from the DB, to the http server to the file
		ds := &DbServiceGeoipMock{
			timestamp: time.Now().Add(-32 * 24 * time.Hour),
			content:   []byte("foobar"),
			countries: make(map[string]bool),
		}

		// let's create a tmp file
		f, err := os.CreateTemp("", "foobar-*.mmdb.gz")
		assert.Nil(t, err)
		defer os.Remove(f.Name())
		_, err = f.Write(data)

		mmdb, err := initMaxmindDb(ds, f.Name(), server.URL+"/foobar-%d-%d.mmdb.gz")
		assert.Nil(t, err)
		assert.NotNil(t, mmdb)
		assert.True(t, downloaded)
	})

	t.Run("test NewGeoipWafPlugin", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		// let's get from the DB, to the http server to the file
		ds := &DbServiceGeoipMock{
			timestamp: time.Now(),
			content:   data,
			countries: make(map[string]bool),
		}
		ds.countries["UK"] = true

		// not downloaded, the timestamp is recent
		plugin, err := NewGeoipWafPlugin("/tmp/foobar.mmdb", "http://localhost:12345/foobar-%d-%d.mmdb.gz", ds)
		assert.Nil(t, err)
		assert.NotNil(t, plugin)
		assert.NotNil(t, plugin.(*GeoipWafPlugin).geoipdb)
		assert.NotNil(t, plugin.(*GeoipWafPlugin).allowDenyTable)
		assert.Equal(t, 1, len(plugin.(*GeoipWafPlugin).allowDenyTable))
		assert.True(t, plugin.(*GeoipWafPlugin).allowMode)
	})

}

func TestScan(t *testing.T) {
	t.Run("scan happy path", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		ds := &DbServiceGeoipMock{
			timestamp: time.Now(),
			content:   data,
			countries: make(map[string]bool),
		}
		mmdb, err := maxminddb.FromBytes(data)
		assert.Nil(t, err)

		// empty allow/deny list -> pass
		plugin := &GeoipWafPlugin{
			ds:             ds,
			allowDenyTable: map[string]bool{},
			allowMode:      true,
			geoipdb:        mmdb,
		}

		url, err := url.Parse("http://www.google.fr")
		assert.Nil(t, err)

		res := plugin.Scan(&com.TaxsiCom{
			RemoteAddr: "34.130.155.108",
			Url:        url,
			Method:     "GET",
		})

		assert.True(t, res)

	})

	t.Run("scan happy path: deny", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		ds := &DbServiceGeoipMock{
			timestamp: time.Now(),
			content:   data,
			countries: make(map[string]bool),
		}
		mmdb, err := maxminddb.FromBytes(data)
		assert.Nil(t, err)

		// empty allow/deny list -> pass
		plugin := &GeoipWafPlugin{
			ds:             ds,
			allowDenyTable: map[string]bool{"FR": false},
			allowMode:      false,
			geoipdb:        mmdb,
		}

		url, err := url.Parse("http://www.google.com")
		assert.Nil(t, err)

		res := plugin.Scan(&com.TaxsiCom{
			RemoteAddr: "34.130.155.108",
			Url:        url,
			Method:     "GET",
		})

		assert.True(t, res)

	})
}

func TestNotification(t *testing.T) {
	t.Run("happy path: receive notification db change", func(t *testing.T) {
		data, err := base64.StdEncoding.DecodeString(base64MmdContent)
		assert.Nil(t, err)

		ds := &DbServiceGeoipMock{
			timestamp: time.Now(),
			content:   data,
			countries: make(map[string]bool),
		}
		ds.countries["UK"] = true

		mmdb, err := maxminddb.FromBytes(data)
		assert.Nil(t, err)

		// empty allow/deny list -> pass
		plugin := &GeoipWafPlugin{
			ds:             ds,
			allowDenyTable: map[string]bool{"FR": false},
			allowMode:      false,
			geoipdb:        mmdb,
		}

		plugin.NotifyDbChange("allow")
		assert.Equal(t, 1, len(plugin.allowDenyTable))
		assert.Equal(t, true, plugin.allowDenyTable["UK"])
		assert.Equal(t, true, plugin.allowMode)
	})
}
