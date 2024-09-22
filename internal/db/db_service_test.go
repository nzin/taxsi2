package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DbChangeListenerMock struct {
	notifications int
	key           string
}

func (l *DbChangeListenerMock) NotifyDbChange(key string) {
	l.notifications++
	l.key = key
}

func TestWatch(t *testing.T) {
	t.Run("happy path: receive notification", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "taxsi.temp*")
		assert.Nil(t, err)
		defer os.Remove(tmpFile.Name())

		dbs, err := NewDbService("sqlite3", tmpFile.Name(), 1, 1*time.Second)
		assert.Nil(t, err)
		assert.NotNil(t, dbs)

		// subscribe
		l := DbChangeListenerMock{}
		dbs.SubscribeChanges(CHANGELOG_TABLE_CONFIG, &l)

		// for fun let's subscribe a second time (it should noop)
		dbs.SubscribeChanges(CHANGELOG_TABLE_CONFIG, &l)

		// start the watcher
		stopChan := make(chan struct{})
		go dbs.Watch(stopChan)
		defer close(stopChan)

		err = dbs.NotifyChange(CHANGELOG_TABLE_CONFIG, "foo")
		assert.Nil(t, err)

		// wait for the event to be read
		time.Sleep(2 * time.Second)

		assert.Equal(t, "foo", l.key)
		assert.Equal(t, 1, l.notifications)
	})

	t.Run("not happy path: wrong notified table", func(t *testing.T) {
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

		// wrong table here
		err = dbs.NotifyChange(CHANGELOG_TABLE_GEOIP, "foo")
		assert.Nil(t, err)

		// wait for the event to be read
		time.Sleep(2 * time.Second)

		assert.Equal(t, "", l.key)
	})

}
