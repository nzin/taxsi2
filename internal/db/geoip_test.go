package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeoip(t *testing.T) {
	t.Run("happy path: read empty table", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// get a countries
		countries, err := dbs.GetGeoipCountries()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(countries))
	})

	t.Run("happy path: set geoips", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// set countries
		err = dbs.SetGeoipCountries([]string{"CA", "UK"}, true)
		assert.Nil(t, err)

		// get a countries
		countries, err := dbs.GetGeoipCountries()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(countries))
		assert.Equal(t, true, countries["CA"])
	})
}
