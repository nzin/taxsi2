package db

import "gorm.io/gorm"

const (
	CHANGELOG_TABLE_CONFIG = iota
)

type ChangeLog struct {
	gorm.Model
	Table int
	Key   string
}

func (ds *DbServiceImpl) NotifyChange(table int, key string) error {
	changelog := ChangeLog{
		Table: table,
		Key:   key,
	}
	return ds.db.Create(&changelog).Error
}
