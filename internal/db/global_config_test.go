package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGlobalConfig(t *testing.T) {
	t.Run("happy path: read empty table", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// get all config
		config, err := dbs.GetConfigs()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(config))

		// get specific
		value, err := dbs.GetConfigValueForKey("foo")
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Equal(t, "", value)
	})

	t.Run("happy path: set config", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// set value
		err = dbs.SetConfigValueForKey("foo", "bar")
		assert.Nil(t, err)

		err = dbs.SetConfigValueForKey("foo2", "bar2")
		assert.Nil(t, err)

		// get back the value
		config, err := dbs.GetConfigs()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(config))

		// get specific
		value, err := dbs.GetConfigValueForKey("foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", value)
	})

	t.Run("happy path: set config and get back notification", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// subscribe
		l := DbChangeListenerMock{}
		dbs.SubscribeChanges(CHANGELOG_TABLE_CONFIG, &l)

		// start the watcher
		stopChan := make(chan struct{})
		go dbs.Watch(stopChan)
		defer close(stopChan)

		// set value
		err = dbs.SetConfigValueForKey("foo", "bar")
		assert.Nil(t, err)

		err = dbs.SetConfigValueForKey("foo2", "bar2")
		assert.Nil(t, err)

		// wait for the event to be read
		time.Sleep(2 * time.Second)

		assert.Equal(t, 2, l.notifications)
		assert.Equal(t, "foo2", l.key)
	})
}
