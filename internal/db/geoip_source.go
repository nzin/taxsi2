package db

import (
	"time"

	"gorm.io/gorm"
)

/*
GeoipSource is a one row table to store the source of the Geoip2-Country.mmdb
Database are usually not good at storing large binary data, but the database
is centric to the application, and I dont want to introduce a new subsystem
*/
type GeoipSource struct {
	Timestamp time.Time
	Source    []byte `gorm:"type:blob"` // Geoip2-Country.mmdb
}

func (ds *DbServiceImpl) GetGeoipSourceTimestamp() (time.Time, error) {
	var gs GeoipSource

	// fetch only the timestamp
	err := ds.db.Select("timestamp").First(&gs).Error

	if err != nil {
		return time.Time{}, err
	}
	return gs.Timestamp, nil
}

func (ds *DbServiceImpl) GetGeoipSource() (*GeoipSource, error) {
	var gs GeoipSource

	err := ds.db.First(&gs).Error
	if err != nil {
		return nil, err
	}
	return &gs, nil
}

func (ds *DbServiceImpl) SetGeoipSource(timestamp time.Time, source []byte) error {
	ds.db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&GeoipSource{})

	gs := GeoipSource{
		Timestamp: time.Now(),
		Source:    source,
	}
	return ds.db.Create(&gs).Error
}
