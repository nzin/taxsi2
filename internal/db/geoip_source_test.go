package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGeoipSource(t *testing.T) {
	t.Run("happy path: read empty table", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// try to get
		source, err := dbs.GetGeoipSource()
		assert.Nil(t, source)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// try to get the timestamp
		_, err = dbs.GetGeoipSourceTimestamp()
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("happy path: save a byte array", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// save a bin file
		data := []byte("foobar")
		now := time.Now()
		err = dbs.SetGeoipSource(now, data)
		assert.Nil(t, err)

		// let's try to get the data back
		source, err := dbs.GetGeoipSource()
		assert.Nil(t, err)
		assert.Equal(t, []byte("foobar"), source.Source)
		assert.Equal(t, now.UTC().String(), source.Timestamp.UTC().String())

		// differently
		timestamp, err := dbs.GetGeoipSourceTimestamp()
		assert.Nil(t, err)
		assert.Equal(t, now.UTC().String(), timestamp.UTC().String())
	})
}
