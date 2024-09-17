package db

type GlobalConfig struct {
	Key   string
	Value string
}

func (ds *DbServiceImpl) GetConfigs() (map[string]string, error) {
	var configs []GlobalConfig

	err := ds.db.Find(&configs).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)

	for _, c := range configs {
		res[c.Key] = c.Value
	}
	return res, nil
}

func (ds *DbServiceImpl) GetConfigValueForKey(key string) (string, error) {
	var config GlobalConfig

	err := ds.db.Where(GlobalConfig{Key: key}).First(&config).Error
	if err != nil {
		return "", err
	}

	return config.Value, nil

}

func (ds *DbServiceImpl) SetConfigValueForKey(key string, value string) error {

	config := GlobalConfig{
		Key:   key,
		Value: value,
	}

	err := ds.db.Save(config).Error
	if err != nil {
		return err
	}

	return ds.NotifyChange(CHANGELOG_TABLE_CONFIG, key)
}
