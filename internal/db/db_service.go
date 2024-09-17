package db

import (
	"fmt"
	"time"

	retry "github.com/avast/retry-go"
	logrus "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var AutoMigrateTables = []interface{}{
	ChangeLog{},
	GlobalConfig{},
}

type DbChangeListener interface {
	NotifyDbChange(key string)
}

type DbService interface {
	// Register to be notified when a changelog event
	// occurs on a given table
	SubscribeChanges(table int, listener DbChangeListener)

	// NotifyChange create a ChangeLog row to notify
	// to other nodes that there is a DB change
	// in a specific table, for a specific key
	NotifyChange(table int, key string) error

	// Watch() start a go routine, that will periodically
	// check if there is any changelog
	Watch(stopChannel chan struct{})

	GetConfigs() (map[string]string, error)

	GetConfigValueForKey(key string) (string, error)

	SetConfigValueForKey(key string, value string) error
}

type DbServiceImpl struct {
	db                   *gorm.DB
	lastChangeLog        uint
	changelogSubscribers map[int][]DbChangeListener
}

func NewDbService(dbtype string, dsn string, retryAttempt uint, retryDelay time.Duration) (dbservice DbService, err error) {
	logger := &Logger{
		LogLevel:                  gorm_logger.Warn,
		SlowThreshold:             1000 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
	var db *gorm.DB

	err = retry.Do(
		func() error {
			switch dbtype {
			case "postgres":
				db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
					Logger: logger,
				})
			case "sqlite3":
				db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
					Logger: logger,
				})
			case "mysql":
				db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
					Logger: logger,
				})
			}
			return err
		},
		retry.Attempts(retryAttempt),
		retry.Delay(retryDelay),
	)
	if err != nil {
		return nil, err
	}

	// let's load the last change log
	var lastChangeLog ChangeLog
	err = db.Order("id desc").First(&lastChangeLog).Error
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(AutoMigrateTables...)
	if err != nil {
		return nil, err
	}

	dbservice = &DbServiceImpl{
		db:                   db,
		lastChangeLog:        lastChangeLog.ID,
		changelogSubscribers: make(map[int][]DbChangeListener),
	}

	return dbservice, err
}

func (ds *DbServiceImpl) Watch(stopChannel chan struct{}) {
	for {
		select {
		case _, ok := <-stopChannel:
			if !ok {
				// channel closed
				return
			}
		case <-time.After(1 * time.Second):
			// check the changelog table
			// let's load the last change log
			var lastChangeLog ChangeLog
			err := ds.db.Order("id desc").First(&lastChangeLog).Error
			if err != nil {
				logrus.Errorf("error looking for changelog: %v", err)
				continue
			}

			for lastChangeLog.ID != ds.lastChangeLog {
				err := ds.notifySubscriberTo(lastChangeLog.ID)
				if err != nil {
					logrus.Errorf("Error reading changelogs: %v", err)
				}
			}
		}

	}
}

// this function will read the ChangeLog table starting with
// ID = ds.lastChangeLog + 1
// and will distribute notifications across notification subscribers
// (i.e. ds.changelogSubscribers)
func (ds *DbServiceImpl) notifySubscriberTo(lastId uint) error {
	for ds.lastChangeLog < lastId {
		newId := ds.lastChangeLog + 1
		var log ChangeLog
		err := ds.db.First(&log, newId).Error
		if err != nil {
			return fmt.Errorf("error looking for changelog id %d: %v", newId, err)
		}

		subscribers := ds.changelogSubscribers[log.Table]
		for _, s := range subscribers {
			s.NotifyDbChange(log.Key)
		}
		ds.lastChangeLog++
	}
	return nil
}

func (ds *DbServiceImpl) SubscribeChanges(table int, listener DbChangeListener) {
	s := ds.changelogSubscribers[table]
	if s == nil {
		s = []DbChangeListener{listener}
		ds.changelogSubscribers[table] = s
	} else {
		for _, l := range s {
			if l == listener {
				return
			}
		}
		ds.changelogSubscribers[table] = append(ds.changelogSubscribers[table], listener)
	}
}
