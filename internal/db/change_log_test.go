package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChangeLog(t *testing.T) {
	t.Run("wrong db type", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = NewDbService("sqlite", tmpFile.Name(), 1, 1*time.Second)
		assert.NotNil(t, err)
	})

	t.Run("changelog empty", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		assert.Equal(t, uint(0), dbs.(*DbServiceImpl).lastChangeLog)
	})

	t.Run("create a first changelog entry", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		assert.Equal(t, uint(0), dbs.(*DbServiceImpl).lastChangeLog)

		err = dbs.NotifyChange(CHANGELOG_TABLE_CONFIG, "foo")
		assert.Nil(t, err)

		var c ChangeLog
		db := dbs.(*DbServiceImpl).db
		err = db.First(&c).Error
		assert.Nil(t, err)
		assert.Equal(t, uint(1), c.ID)
		assert.Equal(t, CHANGELOG_TABLE_CONFIG, c.Table)
		assert.Equal(t, "foo", c.Key)
	})
}
