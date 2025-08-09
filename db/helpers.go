package db

import "gorm.io/gorm"

func DeleteRecord(record interface{}) *gorm.DB {
	return GetDB().Delete(&record)
}

func CreateRecord(record interface{}) *gorm.DB {
	return GetDB().Create(&record)
}
