package repositories

import "gorm.io/gorm"

func getDb(defaultDb *gorm.DB, dbs ...*gorm.DB) *gorm.DB {
	if len(dbs) > 0 && dbs[0] != nil {
		return dbs[0]
	}
	return defaultDb
}
